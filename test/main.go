package main

import (
	"strings"

	_ "embed"

	"github.com/paganotoni/mailopen/v2"
)

//go:embed blank.pdf
var file []byte

//go:embed template.html
var html string

func main() {
	sender := mailopen.WithOptions(
		mailopen.Only("text/html"),
	)

	email := mailopen.Email{
		From:    "a@test.com",
		To:      []string{"a@test.com"},
		Subject: "Test",
		Bodies: []mailopen.Body{
			{
				ContentType: "text/html",
				Content:     html,
			},
			{
				ContentType: "text/plain",
				Content:     "Some Message",
			},
		},
		Attachments: []mailopen.Attachment{
			{
				Name:        "file.txt",
				ContentType: "text/plain",
				Reader:      strings.NewReader("Some file content"),
			},
			{
				Name:        "file.pdf",
				ContentType: "application/pdf",
				Reader:      strings.NewReader(string(file)),
			},
		},
	}

	if err := sender.Send(email); err != nil {
		panic(err)
	}
}
