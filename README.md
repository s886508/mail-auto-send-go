# Decription
Sending emails through Sendgrid or GMail API

# Build
```go
go build -o <output_file> .\main.go
```

# General setup before sending email
1. Update `configs/message.json`.
  ```json
  {
    "from": "Mail send from",
    "email": "xxx@domain.com"
    "subject": "This my first mail",
    "content": "Test mail"
  }
  ```

2. Update `configs/maillist.csv`, the file listed email address to be sent.
  ```csv
  <Name>,<Email>,<Attachment>
  WILL HUANG,a12345@gmail.com,TEST.pdf
  WILL WEI,a2345@gmail.com,TEST2.txt
  ```

# Send EMail with Sendgrid
1. Setup a Sendgrid account and create the API key. [link](https://docs.sendgrid.com/for-developers/sending-email/quickstart-go#create-and-store-a-sendgrid-api-key)

2. Update `configs/config.json`
    - Update `sendgridApiKey` to the created one from first step.

3. Execute and EMail sent
  ```
  [SUCCESS] sending mail to: WILL HUANG a12345@gmail.com
  [SUCCESS] sending mail to: WILL WEI a2345@gmail.com
  ```

# Send EMail with GMail API
1. Google API Setup:
    - Enable google API, read the [page](https://developers.google.com/workspace/guides/create-project).
    - Create GMail authentication credential:<BR>
      Read the [page](https://developers.google.com/workspace/guides/create-credentials) and create a OAuth client id credential.<BR>
      Download credential file and save content to `configs/gmail/credentials.json`.

2. Update `configs/config.json`
    - Update "gamilCredentials" to `configs/gmail/credentials.json`.

3. First time execute and autorization:
   - Execute the program and you will see following statement, follow the link to authorize the application.
   ```
   Go to the following link in your browser then type the authorization code:
   https://<Authorization URL>
   ```
   - Copy the token from browser and enter the token as following.
   ```
   Please enter the auto code from browser: <COPY FROM Browser>
   ```
   
4. EMail sent
  ```
  [SUCCESS] sending mail to: WILL HUANG a12345@gmail.com
  [SUCCESS] sending mail to: WILL WEI a2345@gmail.com
  ```
