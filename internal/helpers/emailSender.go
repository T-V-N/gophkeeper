package helpers

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/T-V-N/gophkeeper/internal/config"
	"github.com/T-V-N/gophkeeper/internal/utils"
)

type EmailSender struct {
	SenderName         string
	RegisterSender     string
	RegisterTemplateID string
	SengridAPIKey      string
}

func InitEmailSender(cfg *config.Config) *EmailSender {
	return &EmailSender{SenderName: cfg.SenderName, RegisterSender: cfg.RegisterSender, RegisterTemplateID: cfg.RegisterTemplateID, SengridAPIKey: cfg.SengridAPIKey}
}

func (es EmailSender) SendConfirmationEmail(to, confirmationURL string) error {
	m := mail.NewV3Mail()

	address := es.RegisterSender
	name := es.SenderName
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	m.SetTemplateID(es.RegisterTemplateID)

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
