package sender

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/s886508/mail-auto-send-go/pkg/cfg"
	"github.com/s886508/mail-auto-send-go/pkg/mailutil"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendFromSendgrid(config *cfg.MailSenderConfig) {
	mailList := mailutil.LoadMailList(config.MaillistConf)
	mailTemplate := mailutil.LoadMailTemplate(config.MailtemplateConf)

	client := sendgrid.NewSendClient(config.ApiKey)
	for _, receiver := range mailList {
		mailToName := receiver.Name
		mailTo := receiver.EMail
		attachmentPath := receiver.AttachmentFile
		message := createSendgridMail(
			mailTemplate.From,
			mailTemplate.EMail,
			receiver.Name,
			mailTo,
			mailTemplate.Subject,
			mailTemplate.Content,
			attachmentPath,
		)
		if message == nil {
			continue
		}

		_, err := client.Send(message)
		if err != nil {
			log.Printf("[FAILED] error sending mail to: %s %s err: %v\n", mailToName, mailTo, err)
			continue
		}

		err = os.Remove(attachmentPath)
		if err != nil {
			log.Printf("Failed to remove attachment file: %s\n", attachmentPath)
		}
		log.Printf("[SUCCESS] sending mail to: %s %s", mailToName, mailTo)
	}
}

func createSendgridMail(
	fromName string,
	from string,
	toName string,
	to string,
	subject string,
	content string,
	attachmentPath string,
) *mail.SGMailV3 {
	_, err := os.Stat(attachmentPath)
	if os.IsNotExist(err) {
		log.Printf("attachment does not exist: %s\n", attachmentPath)
		return nil
	}

	b, err := ioutil.ReadFile(attachmentPath)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
	}
	attachment := mail.NewAttachment()
	attachment.SetContent(base64.StdEncoding.EncodeToString([]byte(b)))
	mimeType := http.DetectContentType(b)
	attachment.SetType(strings.Split(mimeType, ";")[0])
	attachment.SetFilename(filepath.Base(attachmentPath))
	attachment.SetDisposition("attachment")

	fromMail := mail.NewEmail(fromName, from)
	toMail := mail.NewEmail(toName, to)
	message := mail.NewSingleEmail(fromMail, subject, toMail, "", content)
	message.AddAttachment(attachment)

	return message
}
