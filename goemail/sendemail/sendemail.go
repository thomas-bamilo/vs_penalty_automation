package sendemail

import (
	"log"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/scorredoira/email"
)

// SendEmail function is good
func SendEmail(attachPath string, toEmail []string, title string, message string, senderName string, senderEmail string, senderPw string) {

	today := time.Now().Format("Jan 2 15:04:05 MST 2006")

	// compose the message
	m := email.NewMessage(title+" "+today, message)
	m.From = mail.Address{Name: senderName, Address: senderEmail}
	m.To = toEmail

	// add attachments
	if attachPath != "NA" {
		if err := m.Attach(attachPath); err != nil {
			log.Fatal(err)
		}
	}

	// send it
	auth := smtp.PlainAuth("", senderEmail, senderPw, "smtp.aol.com")
	if err := email.Send("smtp.aol.com:587", auth, m); err != nil {
		log.Fatal(err)
	}
}
