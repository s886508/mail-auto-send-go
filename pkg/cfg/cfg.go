package cfg

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type MailSenderConfig struct {
	MaillistConf     string `json:"maillist"`
	MailtemplateConf string `json:"mailTemplate"`
	ApiKey           string `json:"apikey"`
}

// LoadMailSenderConfig loads config from given file path.
func LoadMailSenderConfig(filePath string) *MailSenderConfig {
	if len(filePath) == 0 {
		log.Fatal("empty path to load config")
	}

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Fail to read email content")
	}

	config := &MailSenderConfig{}
	err = json.Unmarshal(b, config)
	if err != nil {
		log.Fatalf("Fail to parse config file")
	}
	return config
}
