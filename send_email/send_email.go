package queryoms

import (
	"log"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/scorredoira/email"
)

// SendEmail function is good
func SendEmail(attachPath string, toEmail []string) {

	today := time.Now().Format("Jan 2 2006")

	// compose the message
	m := email.NewMessage("bamilo_automation_test "+today, "Please find here the attachement! :)))))")
	m.From = mail.Address{Name: "bamilo.automation", Address: "bamilo.automation@aol.com"}
	m.To = toEmail

	// add attachments
	if err := m.Attach(attachPath); err != nil {
		log.Fatal(err)
	}

	// send it
	auth := smtp.PlainAuth("", "bamilo.automation@aol.com", "golangisgreat", "smtp.aol.com")
	if err := email.Send("smtp.aol.com:587", auth, m); err != nil {
		log.Fatal(err)
	}
}
