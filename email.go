package mailopen

import "io"

type Email struct {
	From        string
	To          []string
	Subject     string
	Bodies      []Body
	Attachments []Attachment
	CC          []string
	Bcc         []string
}

// Body is the body of the email
type Body struct {
	ContentType string
	Content     string
}

type Attachment struct {
	ContentType string
	Name        string
	Reader      io.Reader
}
