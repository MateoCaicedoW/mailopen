package mailopen

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"io"
	"mime"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/gofrs/uuid"
	"github.com/pkg/browser"
)

var (

	//go:embed html-header.html
	htmlHeader string

	//go:embed plain-header.txt
	plainHeader string

	// config used to write the email body to files
	// depending on the type we need to do a few things differently.
	config = map[string]contentConfig{
		"text/html": {
			headerTemplate: htmlHeader,
			replaceRegexp:  `(.*<body[^>]*>)((.|[\n\r])*)(<\/body>.*)`,
		},

		"text/plain": {
			headerTemplate: plainHeader,
			replaceRegexp:  `(.*<pre[^>]*>)((.|[\n\r])*)(<\/pre>.*)`,

			preformatter: func(s string) string {
				return fmt.Sprintf("<html><head></head><body><pre>%v</pre></body></html>", s)
			},
		},
	}
)

type contentConfig struct {
	headerTemplate string
	replaceRegexp  string
	preformatter   func(string) string
}

type FileSender struct {
	Open bool
	// dir is the directory to save the files to (attachments and email templates).
	dir string

	// openContentTypes are those content types to open in browser
	openContentTypes []string
}

// AttFile is a file to be attached to the email
type AttFile struct {
	Path string
	Name string
}

func (ps FileSender) shouldOpen(contentType string) bool {
	if len(ps.openContentTypes) == 0 {
		return true
	}

	for _, v := range ps.openContentTypes {
		if v == contentType {
			return true
		}
	}

	return false
}

func (ps FileSender) Send(m Email) error {
	for _, v := range m.GetBodies() {
		if !ps.shouldOpen(v.ContentType()) {
			continue
		}

		cc := config[v.ContentType()]
		content := v.Content()
		if cc.preformatter != nil {
			content = cc.preformatter(v.Content())
		}

		header := fmt.Sprintf(
			cc.headerTemplate,

			html.EscapeString(m.GetFrom()),
			strings.Join(m.GetTo(), ","),
			strings.Join(m.GetCC(), ","),
			strings.Join(m.GetBcc(), ","),

			html.EscapeString(m.GetSubject()),
		)

		var re = regexp.MustCompile(cc.replaceRegexp)
		content = re.ReplaceAllString(content, fmt.Sprintf("$1\n%v\n$2$3", header))
		tmpName := strings.ReplaceAll(v.ContentType(), "/", "_") + "_body"

		path, err := ps.saveEmailBody(content, tmpName, m.GetAttachments())
		if err != nil {
			return err
		}

		if err := browser.OpenFile(path); err != nil {
			return err
		}
	}

	return nil
}

func (ps FileSender) saveEmailBody(content, tmpName string, attachments []Attachment) (string, error) {
	id := uuid.Must(uuid.NewV4())

	afs, err := ps.saveAttachmentFiles(attachments)
	if err != nil {
		return "", fmt.Errorf("mailopen: failed to save attachments: %w", err)
	}

	tmpl := template.Must(template.New("mail").Parse(content))
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, afs)
	if err != nil {
		return "", fmt.Errorf("mailopen: failed to execute template: %w", err)
	}

	filePath := fmt.Sprintf("%s_%s.html", tmpName, id)

	path := path.Join(ps.dir, filePath)
	err = os.WriteFile(path, tpl.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("mailopen: failed to write email body: %w", err)
	}

	return path, nil
}

func (ps FileSender) saveAttachmentFiles(attachments []Attachment) ([]AttFile, error) {
	var afs []AttFile

	for _, a := range attachments {
		name := a.Name()
		if len(name) > 50 {
			name = name[:50]
		}

		exts, err := mime.ExtensionsByType(a.ContentType())
		if err != nil {
			return []AttFile{}, fmt.Errorf("mailopen: failed to get extension for content type %s: %w", a.ContentType(), err)
		}

		filePath := path.Join(ps.dir, fmt.Sprintf("%s_%s%s", uuid.Must(uuid.NewV4()), name, exts[0]))

		b, err := io.ReadAll(a.Reader())
		if err != nil {
			return []AttFile{}, fmt.Errorf("mailopen: failed to read attachment %s: %w", name, err)
		}

		err = os.WriteFile(filePath, b, 0644)
		if err != nil {
			return []AttFile{}, fmt.Errorf("mailopen: failed to write attachment %s: %w", name, err)
		}

		afs = append(afs, AttFile{
			Path: filePath,
			Name: name,
		})
	}

	return afs, nil
}
