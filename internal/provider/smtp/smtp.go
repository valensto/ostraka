package smtp

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/valensto/ostraka/internal/logger"
	"net/url"
	"strconv"
	"time"
)

const SMTP = "smtp"

type Publisher struct {
	parsedURL string
	params    *Params
}

func NewPublisher(params []byte) (*Publisher, error) {
	p, err := unmarshalParams(params)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(p.BaseURL)
	parsedURL := "localhost"
	if err == nil {
		parsedURL = u.Hostname()
	}

	return &Publisher{
		parsedURL: parsedURL,
		params:    p,
	}, nil
}

func (p Publisher) Publish(b []byte) {
	logger.Get().Info().Msgf("sending email from SMTP to %s", p.params.To)
	build := builder{
		from:    p.params.From,
		to:      p.params.To,
		noReply: p.params.From,
		subject: p.params.Subject,
		baseURL: p.parsedURL,
	}

	server := fmt.Sprintf("%s:%s", p.params.Host, p.params.Port)
	auth := authenticate(p.params.Username, p.params.Password, p.params.Host)
	err := build.send(server, p.params.EnableStartTLS, auth, b)
	if err != nil {
		logger.Get().Error().Err(err).Msgf("error sending email from SMTP to %s", p.params.To)
		return
	}
	logger.Get().Info().Msgf("email sent from SMTP to %s", p.params.To)
}

// The message ID is used to show which message is a reply to which other message.
// Define in RFC 2822
// https://datatracker.ietf.org/doc/html/rfc2822#section-3.6.4
func generateEmailID(localName string) string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	randStr := hex.EncodeToString(buf)
	messageID := fmt.Sprintf("<%s.%s@%s>", randStr, timestamp, localName)
	return messageID
}
