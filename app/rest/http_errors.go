package rest

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

// All error codes for UI mapping and translation
const (
	ErrInternal   = 0 // any internal error
	ErrDecodeJSON = 1 // failed to unmarshal json
)

const errorHTML = `<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width"/>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
</head>
<body>
<div style="text-align: center; font-family: Arial, sans-serif; font-size: 18px;">
	<p style="position: relative; max-width: 20em; margin: 0 auto 1em auto; line-height: 1.4em;">{{.Error}}: {{.Details}}.</p>
</div>
</body>
</html>
`

// errData describes parameters of any error
type errData struct {
	error   string
	details string
}

// SendHTMLErrorPage returns an html page if any error, that relates to the user's request occurs
func SendHTMLErrorPage(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string, errCode int) {
	tmpl := template.Must(template.New("error").Parse(errorHTML))
	msg := bytes.Buffer{}
	if err := tmpl.Execute(&msg, errData{
		error:   err.Error(),
		details: details,
	}); err != nil {
		panic(err)
	}
	log.Printf("[WARN] error occured, while processing request: %s, internal error code: %d", err.Error(), errCode)
	render.Status(r, httpStatusCode)
	render.HTML(w, r, msg.String())
}

// SendJSONError returns a json representation of an occurred error
func SendJSONError(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string, errCode int) {
	log.Printf("[WARN] error occured, while processing request: %s, internal error code: %d", err.Error(), errCode)
	render.Status(r, httpStatusCode)
	render.JSON(w, r, errData{
		error:   err.Error(),
		details: details,
	})
}
