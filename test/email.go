package main

import (
	"io"

	"github.com/paganotoni/mailopen/v2"
)

type email struct {
	From        string
	To          []string
	Subject     string
	Bodies      []mailopen.Body
	Attachments []mailopen.Attachment
	CC          []string
	Bcc         []string
}

func (e email) GetFrom() string {
	return e.From
}

func (e email) GetTo() []string {
	return e.To
}

func (e email) GetSubject() string {
	return e.Subject
}

func (e email) GetBodies() []mailopen.Body {
	return e.Bodies
}

func (e email) GetAttachments() []mailopen.Attachment {
	return e.Attachments
}

func (e email) GetCC() []string {
	return e.CC
}

func (e email) GetBcc() []string {
	return e.Bcc
}

type Body struct {
	contentType string
	content     string
}

func (b Body) ContentType() string {
	return b.contentType
}

func (b Body) Content() string {
	return b.content
}

type Attachment struct {
	contentType string
	name        string
	reader      io.Reader
}

func (a Attachment) ContentType() string {
	return a.contentType
}

func (a Attachment) Name() string {
	return a.name
}

func (a Attachment) Reader() io.Reader {
	return a.reader
}
