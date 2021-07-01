package main

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	fmt.Printf("Please enter the auto code from browser: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func loadMailList(path string) [][]string {
	if len(path) == 0 {
		log.Fatal("empty path to load mail list")
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Fail to open mail list file")
	}
	csvReader := csv.NewReader(f)
	mails, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Fail to parse mail list")
	}
	return mails
}

func loadMailTemplate(path string) map[string]interface{} {
	if len(path) == 0 {
		log.Fatal("empty path to load mail template")
	}

	var data map[string]interface{}
	b, err := ioutil.ReadFile("message.json")
	if err != nil {
		log.Fatal("Fail to read email content")
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Fatalf("Fail to parse email content")
	}
	return data
}

func createMail(from string, to string, subject string, content string, attachment string) *gmail.Message {
	subjectEnc := fmt.Sprintf("=?utf-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	mailBody := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n\r\n"+
			"%s",
		from,
		to,
		subjectEnc,
		content,
	)

	return &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(mailBody)),
	}
}

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	mailLists := loadMailList("maillist.csv")
	data := loadMailTemplate("message.json")

	for _, toUser := range mailLists {
		mailTo := toUser[1]
		message := createMail(data["from"].(string), mailTo, data["subject"].(string), data["content"].(string), "")

		_, err = srv.Users.Messages.Send("me", message).Do()
		if err != nil {
			log.Printf("[FAILED] error sending mail to: %s err: %v", mailTo, err)
			log.Println()
		} else {
			log.Printf("[SUCCESS] sending mail to: %s", mailTo)
		}
	}
}
