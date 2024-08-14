package mailopen

import "io"

type Email interface {
	GetFrom() string
	GetTo() []string
	GetSubject() string
	GetBodies() []Body
	GetAttachments() []Attachment
	GetCC() []string
	GetBcc() []string
}

type Body interface {
	ContentType() string
	Content() string
}

type Attachment interface {
	ContentType() string
	Name() string
	Reader() io.Reader
}
