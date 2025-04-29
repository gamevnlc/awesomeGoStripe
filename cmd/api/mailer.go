package main

import (
	"bytes"
	"embed"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"time"
)

//go:embed templates
var emailTemplateFS embed.FS

func (app *application) SendMail(from, to, subject, tmpl string, data interface{}) error {
	templateToRender := fmt.Sprintf("templates/%s.html.tmpl", tmpl)

	t, err := template.New("email-html").ParseFS(emailTemplateFS, templateToRender)

	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", data)

	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	formattedMessage := tpl.String()

	templateToRender = fmt.Sprintf("templates/%s.plain.tmpl", tmpl)

	t, err = template.New("email-plain").ParseFS(emailTemplateFS, templateToRender)

	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	err = t.ExecuteTemplate(&tpl, "body", formattedMessage)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	plainMessage := tpl.String()

	app.infoLog.Println(plainMessage)

	server := mail.NewSMTPClient()

	server.Host = app.config.smtp.host
	server.Port = app.config.smtp.port
	server.Username = app.config.smtp.username
	server.Password = app.config.smtp.password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second

	smptClient, err := server.Connect()
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(from)
	email.AddTo(to)
	email.SetSubject(subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextHTML, plainMessage)

	err = email.Send(smptClient)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	app.infoLog.Println("send mail")

	return nil
}
