package mailopen

import (
	"os"
)

const MailOpenDirKey = "MAILOPEN_DIR"

// WithOptions creates a sender that writes emails into disk
// And applies the passed options.
func WithOptions(options ...Option) FileSender {
	s := FileSender{
		Open: true,
		dir:  getEnv(MailOpenDirKey, os.TempDir()),
	}

	for _, option := range options {
		option(&s)
	}

	return s
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
