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
	sender := mailopen.New(
		mailopen.Only("text/html"),
	)

	email := email{
		From:    "a@test.com",
		To:      []string{"a@test.com"},
		Subject: "Test",
		Bodies: []mailopen.Body{
			Body{
				contentType: "text/html",
				content:     html,
			},
			Body{
				contentType: "text/plain",
				content:     "Some plain text",
			},
		},
		Attachments: []mailopen.Attachment{
			Attachment{
				contentType: "application/pdf",
				name:        "blank.pdf",
				reader:      strings.NewReader(string(file)),
			},
		},
	}

	sender.Send(email)
}
