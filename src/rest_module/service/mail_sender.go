package service

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"

	"github.com/go-mail/mail/v2"
)

type MailSender struct {
	smtpHost string
	smtpPort int
	smtpUser string
	smtpPass string
}

func InitMailSender() (*MailSender, error) {
	// Получение параметров из переменных окружения
	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("Не валидный порт SMTP %s", err.Error())
	}
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")

	sender := MailSender{}
	sender.smtpHost = host
	sender.smtpPort = port
	sender.smtpUser = user
	sender.smtpPass = password
	return &sender, nil
}

// Сообщение для отправки по почте
func (self *MailSender) createMessage(to string, subject string, body string) *mail.Message {
	m := mail.NewMessage()
	m.SetHeader("From", self.smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return m
}

// SMTP-диалог
func (self *MailSender) createDialer() *mail.Dialer {
	d := mail.NewDialer(self.smtpHost, self.smtpPort, self.smtpUser, self.smtpPass)
	d.TLSConfig = &tls.Config{
		ServerName:         self.smtpHost,
		InsecureSkipVerify: false, // Не отключать проверку сертификата
	}
	return d
}

// Отправка письма
func (self *MailSender) sendEmail(m *mail.Message) error {
	// Настройка подключения
	d := self.createDialer()

	// Отправка
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("SMTP error: %v", err)
		return fmt.Errorf("email sending failed")
	}

	return nil
}

// Уведомление о оплате через email
func (self *MailSender) SendEmailMessage(emailTo string, amount float64) error {
	// Создание контента
	content := fmt.Sprintf(`
        <h1>Спасибо за оплату!</h1>
        <p>Сумма: <strong>%.2f RUB</strong></p>
        <small>Это автоматическое уведомление</small>
    `, amount)
	// Подготовка сообщения
	m := self.createMessage(emailTo, "Платеж успешно проведен", content)

	// Отправка
	if err := self.sendEmail(m); err != nil {
		fmt.Printf("Error send mail %s", err.Error())
		return err
	}

	fmt.Printf("Email sent to %s", emailTo)
	fmt.Println()

	return nil
}
