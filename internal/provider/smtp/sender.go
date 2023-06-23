package smtp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

type builder struct {
	from    string
	to      string
	noReply string
	subject string
	baseURL string

	content bytes.Buffer
}

func (s *builder) set(key, value string) {
	s.content.WriteString(key)
	s.content.WriteString(": ")
	s.content.WriteString(value)
	s.content.WriteString("\r\n")
}

func (s *builder) body(body []byte) {
	s.content.WriteString("\r\n")
	s.content.Write(body)
}

func (s *builder) bytes() []byte {
	return s.content.Bytes()
}

func (s *builder) send(serverAddress string, enableStartTLS bool, auth smtp.Auth, body []byte) error {
	s.set("From", s.from)
	s.set("Reply-To", s.noReply)
	s.set("To", s.to)
	s.set("Subject", s.subject)
	s.set("MIME-version", "1.0")
	s.set("Content-Type", "text/html; charset=\"UTF-8\"")
	s.set("Date", time.Now().Format(time.RFC1123Z))
	s.set("Message-ID", generateEmailID(s.baseURL))
	s.body(body)

	c, err := smtp.Dial(serverAddress)
	if err != nil {
		return fmt.Errorf("smtp server dial error: %w", err)
	}
	defer c.Close()

	if err = c.Hello(s.baseURL); err != nil {
		return fmt.Errorf("smtp server hello error: %w", err)
	}

	host, _, _ := net.SplitHostPort(serverAddress)
	if enableStartTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			config := &tls.Config{ServerName: host}
			if err = c.StartTLS(config); err != nil {
				return fmt.Errorf("smtp server starttls error: %w", err)
			}
		}
	}

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return fmt.Errorf("smtp server does not support AUTH")
		}
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("smtp server auth error: %w", err)
		}
	}
	if err = c.Mail(s.from); err != nil {
		return fmt.Errorf("smtp server mail error: %w", err)
	}

	if err = c.Rcpt(s.to); err != nil {
		return fmt.Errorf("smtp server rcpt error: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp server data error: %w", err)
	}

	_, err = w.Write(s.bytes())
	if err != nil {
		return fmt.Errorf("smtp server write error: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("smtp server close error: %w", err)
	}

	return c.Quit()
}
