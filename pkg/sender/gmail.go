package sender

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/s886508/mail-auto-send-go/pkg/mailutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func SendFromGMail(ctx context.Context) {
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

	mailList := mailutil.LoadMailList("maillist.csv")
	mailTemplate := mailutil.LoadMailTemplate("message.json")

	for _, receiver := range mailList {
		mailToName := receiver.Name
		mailTo := receiver.EMail
		attachmentPath := receiver.AttachmentFile
		message := createGMail(
			mailTemplate.EMail,
			mailTo,
			mailTemplate.Subject,
			mailTemplate.Content,
			attachmentPath,
		)
		if message == nil {
			continue
		}

		_, err = srv.Users.Messages.Send("me", message).Do()
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

func randStr(strSize int, randType string) string {
	var dictionary string
	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var strBytes = make([]byte, strSize)
	_, _ = rand.Read(strBytes)
	for k, v := range strBytes {
		strBytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(strBytes)
}

func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune
	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result = ""
	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}
	return result
}

func createGMail(from string, to string, subject string, content string, attachment string) *gmail.Message {
	_, err := os.Stat(attachment)
	if os.IsNotExist(err) {
		log.Printf("attachment does not exist: %s\n", attachment)
		return nil
	}

	fileBytes, err := ioutil.ReadFile(attachment)
	if err != nil {
		log.Printf("Fail to load attachment file: %v\n", err)
		return nil
	}

	fileMIMEType := http.DetectContentType(fileBytes)
	fileData := base64.StdEncoding.EncodeToString(fileBytes)
	boundary := randStr(32, "alphanum")

	subjectEnc := fmt.Sprintf("=?utf-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	mailBody := fmt.Sprintf(
		"Content-Type: multipart/mixed; boundary=%s\n"+
			"MIME-Version: 1.0\n"+
			"To: %s\n"+
			"Subject: %s\n\n"+
			"--%s\n"+
			"Content-Type: text/plain; charset=\"UTF-8\"\n"+
			"MIME-Version: 1.0\n"+
			"Content-Transfer-Encoding: 7bit\n\n"+
			"%s\n\n"+
			"--%s\n"+
			"Content-Type: %s; name=\"%s\"\n"+
			"MIME-Version: 1.0\n"+
			"Content-Transfer-Encoding: base64\n"+
			"Content-Disposition: attachment; filename=\"%s\"\n\n"+
			"%s--%s--",
		boundary,
		to,
		subjectEnc,
		boundary,
		content,
		boundary,
		fileMIMEType,
		attachment,
		attachment,
		chunkSplit(fileData, 76, "\n"),
		boundary,
	)

	return &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(mailBody)),
	}
}
