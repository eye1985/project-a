package email

import (
	"context"
	"fmt"
	"github.com/mailersend/mailersend-go"
	"log"
	"time"
)

type SendEmailFromMailSendArgs struct {
	Email    string
	Name     string
	IsSignUp bool
	Code     string
	Origin   string
}

func SendEmailFromMailSend(token string, args *SendEmailFromMailSendArgs) error {
	email := args.Email
	name := args.Name
	isSignUp := args.IsSignUp
	code := args.Code
	origin := args.Origin

	templateId := "zr6ke4n80q34on12"
	ms := mailersend.NewMailersend(token)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	subject := "Sign in link"
	text := "Below is your sign in link."
	signInOrUp := "signin"
	if isSignUp {
		subject = "Sign up link"
		text = "Thank you for signing up. Below is your sign up link."
		signInOrUp = "signup"
	}
	url := fmt.Sprintf("%s/%s/%s", origin, signInOrUp, code)

	from := mailersend.From{
		Name:  "noreply",
		Email: "noreply@test-ywj2lpno30pg7oqz.mlsender.net",
	}

	recipients := []mailersend.Recipient{
		{
			Name:  name,
			Email: email,
		},
	}

	personalization := []mailersend.Personalization{
		{
			Email: email,
			Data: map[string]interface{}{
				"name":       name,
				"action_url": url,
				"text":       text,
			},
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetTemplateID(templateId)
	message.SetPersonalization(personalization)

	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		return err
	}
	sentEmailId := res.Header.Get("X-Message-Id")
	log.Println(sentEmailId)

	return nil
}
