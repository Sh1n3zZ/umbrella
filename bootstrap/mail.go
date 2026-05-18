package bootstrap

import (
	"log"

	"github.com/wneessen/go-mail"
)

func NewMailClient(config *Config) *mail.Client {
	mc := config.Mail
	if mc.Host == "" {
		log.Fatal("mail.host is required")
	}

	port := mc.Port
	if port <= 0 {
		port = 587
	}

	opts := []mail.Option{
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithPort(port),
	}
	if mc.Username != "" {
		opts = append(opts, mail.WithUsername(mc.Username))
	}
	if mc.Password != "" {
		opts = append(opts, mail.WithPassword(mc.Password))
	}

	client, err := mail.NewClient(mc.Host, opts...)
	if err != nil {
		log.Fatal("Failed to create mail client: ", err)
	}

	log.Println("Successfully created mail client")
	return client
}

func CloseMailClient(client *mail.Client) {
	if client == nil {
		return
	}

	if err := client.Close(); err != nil {
		log.Println("Failed to close mail client: ", err)
		return
	}

	log.Println("Mail client closed.")
}
