# Decription
Sending emails through GMail API

# Build
```go
go build -o <output_file> .\main.go
```

# Setup
1. Google API Setup:
    - Enable google API, read the [page](https://developers.google.com/workspace/guides/create-project).
    - Create GMail authentication credential:<BR>
      Read the [page](https://developers.google.com/workspace/guides/create-credentials) and create a OAuth client id credential.<BR>
      Download credential file and save content to `credentials.json`.

2. Modify `message.json`.
  ```json
  {
      "from": "Mail Send from",
      "subject": "This my first mail",
      "content": "Test mail"
  }
  ```

3. Update `maillist.csv`, the file listed email address to be sent.
  ```csv
  <Name>,<Email>
  WILL HUANG,a12345@gmail.com
  WILL WEI,a2345@gmail.com
  ```

4. First time execute and autorization:
   - Execute the program and you will see following statement, follow the link to authorize the application.
   ```
   Go to the following link in your browser then type the authorization code:
   https://<Authorization URL>
   ```
   - Copy the token from browser and enter the token as following.
   ```
   Please enter the auto code from browser: <COPY FROM Browser>
   ```
   
 5. EMail sent
    ```
    [SUCCESS] sending mail to: WILL HUANG a12345@gmail.com
    [SUCCESS] sending mail to: WILL WEI a2345@gmail.com
    ```
