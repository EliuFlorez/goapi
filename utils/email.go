package utils

import (
	"fmt"
	"os"

	"github.com/mailjet/mailjet-apiv3-go"
)

func EmailSend(templateId int, subject string, to string, name string, message string) {
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "eliugdx@gmail.com",
				Name:  "Eliu",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: "eliufz@gmail.com",
					Name:  "Eliu",
				},
			},
			Subject: "Greetings from Mailjet.",
			Variables: map[string]interface{}{
				"name":    name,
				"email":   to,
				"message": message,
			},
			TemplateID: templateId,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)

	if err != nil {
		fmt.Printf("Data: %+v\n", res)
		panic(err)
	}

	//fmt.Printf("Data: %+v\n", res)
}
