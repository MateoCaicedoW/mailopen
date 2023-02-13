package mailopen_test

import (
	"fmt"
	"mime"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/paganotoni/mailopen/v2"

	"github.com/stretchr/testify/require"
)

type falseSender struct{}

func (ps falseSender) Send(m mail.Message) error {
	return nil
}

const (
	txtFormat = `From: %v <br>
		To: %v <br>
		Cc: %v <br>
		Bcc: %v <br>
		Subject: %v <br>
		----------------------------`
)

func Test_Send(t *testing.T) {
	r := require.New(t)

	mailopen.Testing = true

	m := mail.NewMessage()
	m.From = "testing@testing.com"
	m.To = []string{"testing@other.com"}
	m.CC = []string{"aa@other.com"}
	m.Bcc = []string{"aax@other.com"}
	m.Subject = "something"

	const testHTMLcontent = `<html><head></head><body><div>Some Message</div></body></html>`
	const testPlainContent = "Same message"

	t.Run("html and plain with attachments", func(t *testing.T) {
		sender := mailopen.WithOptions()
		sender.Open = false
		sender.TempDir = t.TempDir()

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: testHTMLcontent},
			{ContentType: "text/plain", Content: testPlainContent},
		}

		m.Attachments = []mail.Attachment{
			{Name: "txt_test", Reader: strings.NewReader(""), ContentType: "text/plain", Embedded: false},
			{Name: "csv_test", Reader: strings.NewReader(""), ContentType: "text/csv", Embedded: false},
			{Name: "img_test", Reader: strings.NewReader(""), ContentType: "image/jpeg", Embedded: false},
			{Name: "pdf_test", Reader: strings.NewReader(""), ContentType: "application/pdf", Embedded: false},
			{Name: "zip_test", Reader: strings.NewReader(""), ContentType: "application/zip", Embedded: false},
		}

		r.NoError(sender.Send(m))

		htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[0].ContentType, "/", "_")))
		txtFile := path.Join(sender.TempDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[1].ContentType, "/", "_")))

		r.FileExists(htmlFile)
		r.FileExists(txtFile)

		htmlHeader, err := os.ReadFile(htmlFile)
		r.NoError(err)

		fmt.Println(string(htmlHeader))

		r.Contains(string(htmlHeader), m.From)
		r.Contains(string(htmlHeader), m.To[0])
		r.Contains(string(htmlHeader), m.CC[0])
		r.Contains(string(htmlHeader), m.Bcc[0])
		r.Contains(string(htmlHeader), m.Subject)

		for _, a := range m.Attachments {
			r.Contains(string(htmlHeader), a.Name)

			ext, err := mime.ExtensionsByType(a.ContentType)
			r.NoError(err)

			filePath := path.Join(sender.TempDir, fmt.Sprintf("%s%s", a.Name, ext[0]))
			r.FileExists(filePath)
		}

		txtHeader, err := os.ReadFile(txtFile)
		r.NoError(err)
		format := strings.ReplaceAll(txtFormat, "\t", "")

		r.Equal(string(txtHeader), fmt.Sprintf(format, m.From, m.To[0], m.CC[0], m.Bcc[0], m.Subject))
	})

	t.Run("html only", func(t *testing.T) {
		sender := mailopen.WithOptions(mailopen.Only("text/html"))
		sender.Open = false
		sender.TempDir = t.TempDir()

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: "<html><head></head><body><div>Some Message</div></body></html>"},
			{ContentType: "text/plain", Content: "Same message"},
		}

		htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[0].ContentType, "/", "_")))
		txtFile := path.Join(sender.TempDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[1].ContentType, "/", "_")))

		r.NoError(sender.Send(m))

		r.FileExists(htmlFile)
		r.NoFileExists(txtFile)

		dat, err := os.ReadFile(htmlFile)
		r.NoError(err)

		r.NotContains(string(dat), "Attachment:")
	})

	t.Run("plain only", func(t *testing.T) {
		sender := mailopen.WithOptions(mailopen.Only("text/plain"))
		sender.Open = false
		sender.TempDir = t.TempDir()

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: "<html><head></head><body><div>Some Message</div></body></html>"},
			{ContentType: "text/plain", Content: "Same message"},
		}

		htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[0].ContentType, "/", "_")))
		txtFile := path.Join(sender.TempDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[1].ContentType, "/", "_")))

		r.NoError(sender.Send(m))

		r.NoFileExists(htmlFile)
		r.FileExists(txtFile)
	})

	t.Run("long subject and long file name`", func(t *testing.T) {
		sender := mailopen.WithOptions()
		sender.Open = false
		sender.TempDir = t.TempDir()

		m := mail.NewMessage()

		m.Subject = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam nec leo tellus. Aliquam ac facilisis est, condimentum pellentesque velit. In quis erat turpis. Morbi accumsan ante nec nunc dapibus, quis lacinia mi ornare. Vivamus venenatis accumsan dolor ac placerat. Sed pulvinar sem eu est accumsan, ut commodo mi viverra. Quisque turpis metus, ultrices id mauris vel, suscipit sollicitudin erat. Vivamus eget quam non sem volutpat eleifend eget in lacus. In vulputate, justo fringilla lacinia lobortis, neque turpis dignissim tellus, in placerat eros justo nec massa. Duis ex enim, convallis ut leo nec, condimentum consectetur mi. Vestibulum imperdiet pharetra ipsum. Etiam venenatis tincidunt odio, sed feugiat quam blandit sit amet. Donec eget nulla dui."
		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: "<html><head></head><body><div>Some Message</div></body></html>"},
			{ContentType: "text/plain", Content: "Same message"},
		}

		m.Attachments = []mail.Attachment{
			{Name: "123456789-123456789-123456789-123456789-123456789-1", Reader: strings.NewReader(""), ContentType: "text/plain", Embedded: false},
		}

		r.NoError(sender.Send(m))

		att := m.Attachments[0]

		exts, err := mime.ExtensionsByType(att.ContentType)
		r.NoError(err)

		filePath := path.Join(sender.TempDir, fmt.Sprintf("%s%s", att.Name[0:50], exts[0]))
		r.FileExists(filePath)
	})

	t.Run("only one body", func(t *testing.T) {
		sender := mailopen.WithOptions()
		sender.Open = false
		sender.TempDir = t.TempDir()

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: "<html><head></head><body><div>Some Message</div></body></html>"},
		}

		r.Error(sender.Send(m))
	})
}

func Test_Wrap(t *testing.T) {
	mailopen.Testing = true

	r := require.New(t)

	os.Setenv("GO_ENV", "development")
	s := mailopen.Wrap(falseSender{})
	r.IsType(mailopen.FileSender{}, s)

	os.Setenv("GO_ENV", "")
	s = mailopen.Wrap(falseSender{})
	r.IsType(mailopen.FileSender{}, s)

	os.Setenv("GO_ENV", "staging")
	s = mailopen.Wrap(falseSender{})
	r.IsType(falseSender{}, s)

	os.Setenv("GO_ENV", "production")
	s = mailopen.Wrap(falseSender{})
	r.IsType(falseSender{}, s)

}
