package mailutil

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type MailListItem struct {
	Name           string
	EMail          string
	AttachmentFile string
}

type MailTemplate struct {
	From    string `json:"from"`
	EMail   string `json:"email"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// LoadMailList loads mail list from file.
func LoadMailList(filePath string) []*MailListItem {
	if len(filePath) == 0 {
		log.Fatal("empty path to load mail list")
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Fail to open mail list file")
	}
	csvReader := csv.NewReader(f)
	mails, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Fail to parse mail list")
	}
	mailList := make([]*MailListItem, 0)
	for _, row := range mails {
		item := &MailListItem{
			Name:           row[0],
			EMail:          row[1],
			AttachmentFile: row[2],
		}
		mailList = append(mailList, item)
	}
	return mailList

}

// LoadMailTemplate loads a mail body template
func LoadMailTemplate(filePath string) *MailTemplate {
	if len(filePath) == 0 {
		log.Fatal("empty path to load mail template")
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Fail to read email content")
	}
	template := &MailTemplate{}
	err = json.Unmarshal(b, template)
	if err != nil {
		log.Fatalf("Fail to parse email content")
	}
	return template
}
