package helpers

import (
	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/utils"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailSender struct {
	RegisterTemplateID string
	SengridAPIKey      string
}

func InitEmailSender(cfg *config.Config) *EmailSender {
	return &EmailSender{RegisterTemplateID: cfg.RegisterTemplateID, SengridAPIKey: cfg.SengridAPIKey}
}

func (es EmailSender) SendConfirmationEmail(to, confirmationURL string) error {
	m := mail.NewV3Mail()

	address := "test@truepnl.com"
	name := "Admin"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	m.SetTemplateID("d-99f48121d63f4e58b67098dd063bc82d")

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail("Hey", to),
	}
	p.AddTos(tos...)
	p.SetDynamicTemplateData("link", confirmationURL)

	m.AddPersonalizations(p)

	request := sendgrid.GetRequest(es.SengridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		return err
	} else {
		if response.StatusCode != 202 {
			return utils.ErrThirdParty
		}
	}
	return err
}
