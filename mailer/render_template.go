package mailer

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/flowchartsman/handlebars/v3"
	"github.com/scrambledeggs/booky-go-common/logs"
	"github.com/vanng822/go-premailer/premailer"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed templates/scss/*.scss
var scssFS embed.FS

type File struct {
	Fs       embed.FS
	FileName string
}

type RenderConfig struct {
	Templates    []File
	Context      map[string]interface{}
	StyleSheets  []File
	ImageBaseUrl string
	CompileType  string
	ScssFs       embed.FS
	ScssDir      string
}

type User struct {
	FirstName string
	LastName  string
	Email     string
}

func RenderTemplate(config RenderConfig) (string, error) {
	headerContent, err := templateFS.ReadFile("templates/header.hbs")
	if err != nil {
		return "", err
	}

	footerContent, err := templateFS.ReadFile("templates/footer.hbs")
	if err != nil {
		return "", err
	}

	templateString := string(headerContent)
	// join templates in order
	for _, template := range config.Templates {
		templateContent, err := template.Fs.ReadFile(template.FileName)
		if err != nil {
			return "", err
		}
		templateString += string(templateContent)
	}
	templateString += string(footerContent)

	var styleString string

	if config.CompileType == "scss" {
		// TODO: use godartsass, extract to function
		sassPath := "/opt/bin/sass" // from lambda layer
		tempDir := "/tmp/" + config.ScssDir
		os.MkdirAll(tempDir, 0755)

		var files []fs.DirEntry

		files, err = config.ScssFs.ReadDir(config.ScssDir)
		if err != nil {
			return "", err
		}

		// copy scss files to tempDir
		for _, file := range files {
			if file.IsDir() {
				// TODO recursive copy?
				continue
			}

			var fileContent []byte

			fileContent, err = config.ScssFs.ReadFile(filepath.Join(config.ScssDir, file.Name()))
			if err != nil {
				return "", err
			}

			os.WriteFile(filepath.Join(tempDir, file.Name()), fileContent, 0644)
			// logs.Print("copied scss file", filepath.Join(tempDir, file.Name()))
		}

		for _, styleSheet := range config.StyleSheets {
			fileName := strings.ReplaceAll(styleSheet.FileName, config.ScssDir, tempDir)
			outputFile := fileName + ".css"
			cmd := exec.Command(sassPath, fileName, outputFile, "--style=compressed", "--no-source-map")

			var output []byte

			output, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("scss error", string(output))
				return "", err
			}

			cssContent, err := os.ReadFile(outputFile)
			if err != nil {
				fmt.Println("scss read error", string(output))
				return "", err
			}

			styleString += string(cssContent)
			// fmt.Println("styleString", styleString)
		}
	} else {
		styleString, err = extractStyleSheets(config.StyleSheets)
		if err != nil {
			return "", err
		}

		baseCssContent, err := templateFS.ReadFile("templates/styles/base.css")
		if err != nil {
			return "", err
		}
		styleString = string(baseCssContent) + styleString

		// logs.Print("css styleString", styleString)
	}

	mediaQueriesString := extractMediaQueries(styleString)

	template, err := handlebars.Parse(templateString)
	if err != nil {
		return "", err
	}

	stylePartial := ""

	if styleString != "" {
		styleString = "<style>" + styleString + "</style>"
		stylePartial = styleString
	}

	template.RegisterPartial("style", stylePartial)
	mediaQueriesPartial := ""

	if mediaQueriesString != "" {
		mediaQueriesString = "<style>" + mediaQueriesString + "</style>"
		mediaQueriesPartial = mediaQueriesString
	}

	template.RegisterPartial("mediaQueries", mediaQueriesPartial)

	template.RegisterHelper("image", func(path string) string {
		return strings.TrimSuffix(config.ImageBaseUrl, "/") + "/" + transformPath(path)
	})

	// get raw html with style tags
	baseHtml, err := template.Exec(config.Context)
	if err != nil {
		return "", err
	}

	premailer, err := premailer.NewPremailerFromString(baseHtml, &premailer.Options{
		RemoveClasses:   false, // false for media queries
		CssToAttributes: true,
	})
	if err != nil {
		return "", err
	}

	html, err := premailer.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func extractStyleSheets(styleSheets []File) (string, error) {
	styleString := ""

	for _, styleSheet := range styleSheets {
		styleSheetContent, err := styleSheet.Fs.ReadFile(styleSheet.FileName)
		if err != nil {
			return "", err
		}
		styleString += string(styleSheetContent)
	}

	return styleString, nil
}

func extractMediaQueries(cssString string) string {
	mediaQueryRegex := regexp.MustCompile(`@media[^{]+\{(?:[^{}]|\{[^{}]*\})*\}`)
	mediaQueries := mediaQueryRegex.FindAllString(cssString, -1)

	return strings.Join(mediaQueries, "\n")
}

func transformPath(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// SIH Path format: <BUCKET>,<PATH>,<WIDTH>,<HEIGHT>
	params := strings.Split(path, ",")
	paramsLength := len(params)

	if paramsLength < 2 {
		logs.Error("Invalid path format", path)

		return path
	}

	bucket := params[0]
	key := params[1]

	var width = ""
	var height = ""

	if len(params) > 2 {
		width = params[2]
	}

	if len(params) > 3 {
		height = params[3]
	}

	if paramsLength > 4 {
		logs.Error("Invalid path format", path)
	}

	payload := map[string]interface{}{
		"bucket": bucket,
		"key":    normalizeKey(key),
	}

	if width != "" || height != "" {
		edits := map[string]interface{}{}
		if width != "" {
			edits["width"] = width
		}
		if height != "" {
			edits["height"] = height
		}
		payload["edits"] = edits
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logs.Error("transformPath", err.Error())

		return path
	}

	// logs.Print("transformPath", string(jsonPayload))

	transformedPath := base64.URLEncoding.EncodeToString(jsonPayload)

	return transformedPath
}

func normalizeKey(key string) string {
	key = strings.TrimPrefix(key, "/")

	return decodeURI(key)
}

func decodeURI(encodedString string) string {
	decodedString, err := url.QueryUnescape(encodedString)
	if err != nil {
		logs.Error("decodeURI", err.Error())

		return ""
	}

	return decodedString
}
