package mailer

import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"regexp"
	"strings"

	"github.com/flowchartsman/handlebars/v3"
	"github.com/scrambledeggs/booky-go-common/logs"
	"github.com/vanng822/go-premailer/premailer"
)

type File struct {
	Fs       embed.FS
	FileName string
}

type RenderConfig struct {
	Context             map[string]interface{}
	Templates           []File
	StyleSheets         []File
	ImageBaseUrl        string
	CompileType         string
	UseDefaultTemplates bool
}

//go:embed templates/*
var templateFS embed.FS

func RenderTemplate(config RenderConfig) (string, error) {
	templateString, err := extractTemplates(config)

	if err != nil {
		return "", err
	}

	styleString, err := extractStyleSheets(config)

	if err != nil {
		return "", err
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

func extractTemplates(config RenderConfig) (string, error) {
	templateString := ""

	if config.UseDefaultTemplates {
		config.Templates = append([]File{{Fs: templateFS, FileName: "templates/header.hbs"}}, config.Templates...)
		config.Templates = append(config.Templates, File{Fs: templateFS, FileName: "templates/footer.hbs"})
	}

	for _, template := range config.Templates {
		templateContent, err := template.Fs.ReadFile(template.FileName)

		if err != nil {
			return "", err
		}

		templateString += string(templateContent)
	}

	return templateString, nil
}

func extractStyleSheets(config RenderConfig) (string, error) {
	styleString := ""

	if config.UseDefaultTemplates {
		config.StyleSheets = append([]File{{Fs: templateFS, FileName: "templates/styles/base.css"}}, config.StyleSheets...)
	}

	for _, styleSheet := range config.StyleSheets {
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
