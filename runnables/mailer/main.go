package main

import (
	"embed"

	"github.com/scrambledeggs/booky-go-common/logs"
	"github.com/scrambledeggs/booky-go-common/mailer"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed templates/scss/*.scss
var scssFS embed.FS

func main() {
	variant := ""
	templates := []mailer.File{
		{Fs: templateFS, FileName: "templates/welcome.hbs"},
	}

	user := mailer.User{
		Name:  "Mervyl Canlas",
		Email: "mervyl@phonebooky.com",
	}

	styleSheets := []mailer.File{
		{Fs: templateFS, FileName: "templates/styles/welcome-email.css"},
	}

	if variant == "scss" {
		styleSheets = []mailer.File{{Fs: scssFS, FileName: "templates/scss/welcome-email.scss"}}
	}

	const title = "Welcome To Booky!"

	config := mailer.RenderConfig{
		Templates:   templates,
		StyleSheets: styleSheets,
		Context: map[string]interface{}{
			"title":           title,
			"name":            user.Name,
			"unsubscribe_url": "https://booky.ph/unsubscribe",
		},
		ImageBaseUrl: "https://np-sih.booky.ph",
		CompileType:  variant,
	}

	if variant == "scss" {
		config.ScssFs = scssFS
		config.ScssDir = "templates/scss"
	}

	html, err := mailer.RenderTemplate(config)
	if err != nil {
		logs.Error("RenderTemplates", err.Error())

		return
	}

	err = mailer.Send(user.Email, "no-reply@phonebooky.com", html, title)
	if err != nil {
		logs.Error("Send", err.Error())

		return
	}

	return
}
