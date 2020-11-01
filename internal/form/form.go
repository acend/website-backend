package form

import (
	"backend/internal/config"
	"fmt"
	"gopkg.in/gomail.v2"
	"strings"
)

func Submit(fields map[string]string) error {
	cfg := config.Values.SMTP

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", cfg.To)
	m.SetHeader("Subject", "acend.ch request from: "+fields["email"])

	var str strings.Builder
	for key, value := range fields {
		str.WriteString(fmt.Sprintf("<h3>%s</h3>", key))
		str.WriteString(fmt.Sprintf("<p>%s</p>", value))
	}
	m.SetBody("text/html", str.String())

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	d.LocalName = "acend.ch"
	err := d.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("send mail: %w", err)
	}

	return nil
}
