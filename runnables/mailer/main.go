package main

import (
	"embed"
	"os"

	"github.com/scrambledeggs/booky-go-common/logs"
	"github.com/scrambledeggs/booky-go-common/mailer"
)

//go:embed templates/*
var templateFS embed.FS

var IMAGE_BASE_URL = os.Getenv("IMAGE_BASE_URL")

type User struct {
	Name  string
	Email string
}

func main() {
	user := User{
		Name:  "Mervyl Canlas",
		Email: "mervyl@phonebooky.com",
	}

	title := "Welcome to Booky!"

	html, err := mailer.RenderTemplate(mailer.RenderConfig{
		ImageBaseUrl:        IMAGE_BASE_URL,
		UseDefaultTemplates: true,
		Templates: []mailer.File{
			{Fs: templateFS, FileName: "templates/welcome.hbs"},
		},
		StyleSheets: []mailer.File{
			{Fs: templateFS, FileName: "templates/styles/welcome-email.css"},
		},
		Context: map[string]interface{}{
			"title": title,
			"name":  user.Name,
		},
	})

	if err != nil {
		logs.Error("RenderTemplate", err.Error())

		return
	}

	err = mailer.Send(mailer.SendConfig{
		Sender:    "no-reply@booky.ph",
		Recipient: user.Email,
		Body:      html,
		Subject:   title,
	})

	if err != nil {
		logs.Error("Send", err.Error())

		return
	}

	return
}
