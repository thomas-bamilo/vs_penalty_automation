package goemail

import (
	"strings"

	emailconf "github.com/thomas-bamilo/goemail/emailconf"
	sendemail "github.com/thomas-bamilo/goemail/sendemail"
)

// GoEmail function reads the config.yaml file in the same folder as the application and sends the email according to config.yaml.
// To understand how config.yaml should be built, please look for yaml file format on Google and also look at goemail/emailconf/emailconf.go
func GoEmail() {

	var emailconf emailconf.EmailConf

	emailconf.ReadYamlEmailConf()

	sendemail.SendEmail(emailconf.EmailAttachPath, // attachPath
		strings.Split(emailconf.EmailRecipient, ","), // toEmail
		emailconf.EmailTitle,                         // title
		emailconf.EmailBody,                          // message
		emailconf.SenderName,
		emailconf.SenderEmail,
		emailconf.SenderPw,
	)

}
