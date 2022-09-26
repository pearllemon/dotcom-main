package pages

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Mail struct {
	Subject string
	To      string
	From    string
	CC      []string
	BCC     []string
	ReplyTo string
	Body    string
}

func sendMail(email string, subject string, cc, bcc []string, body string, sender string) error {

	m := Mail{
		Subject: subject,
		To:      email,
		From:    sender,
		CC:      cc,
		BCC:     bcc,
		Body:    body,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)

	res, err := http.Post("http://stl-mailer:9000", "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}
