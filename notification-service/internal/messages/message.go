package messages

import (
	"io/ioutil"
	"net/smtp"
	"os"
	"path/filepath"
)

// WELCOME_MESSAGE é a mensagem de boas-vindas que será enviada por email
const WELCOME_MESSAGE = "Bem vindo ao nosso serviço! Estamos felizes em tê-lo conosco."

// SendWelcomeEmailMessage envia um email de boas-vindas para o usuário
func SendWelcomeEmailMessage(to string) error {

	from := os.Getenv("ROOT_EMAIL")
	password := os.Getenv("ROOT_EMAIL_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	templatePath := filepath.Join("public", "welcome-email.html")

	body, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return err
	}

	message := []byte("Subject: Bem-vindo ao nosso serviço!\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		string(body))

	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message); err != nil {
		return err
	}	

	return nil
}