# Bot Framework File Upload Sample

This Microsoft Teams bot example uses the [msbotbuilder-go](https://github.com/infracloudio/msbotbuilder-go) library. It shows how to handle `invoke` activity and upload local file from Bot to use.

As documented [here](https://developer.microsoft.com/en-us/microsoft-teams/blogs/working-with-files-in-your-microsoft-teams-bot/), following is the workflow for sending files from Bot to user:

1. **Request permission to upload the file:**
    First, ask the user for permission to upload a file by sending a file `consent card`.
2. **User accepts or declines the file**
    When the user presses either “Allow” or “Decline”, your bot will receive an invoke activity.
3. **Handle invoke activity**
    If the file was accepted, activity value contains an uploadInfo property. To set the file contents, issue PUT requests to the URL in uploadInfo.uploadUrl as described [here](https://docs.microsoft.com/en-us/onedrive/developer/rest-api/api/driveitem_createuploadsession?view=odsp-graph-online#upload-bytes-to-the-upload-session)
4. **Delete consent card**
    Once user responds to file upload request, the `consent card` is deleted.
4. **Send the user a link to the uploaded file**
    After you finish the upload, send a link to the file with file `info card`.

## Run the example

#### Step 1: Register MS BOT

Follow the official [documentation](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/create-a-bot-for-teams#register-your-web-service-with-the-bot-framework) to create and register your BOT.

Copy `APP_ID` and `APP_PASSWORD` generated for your BOT.

Do not set any messaging endpoint for now.

#### Step 2: Run echo local server

Set two variables for the session as `APP_ID` and `APP_PASSWORD` to the values of your BotFramework `APP_ID` and `APP_PASSWORD`. Then run the `main.go` file.

```bash
export APP_PASSWORD=MICROSOFT_APP_PASSWORD
export APP_ID=MICROSOFT_APP_ID

go run main.go
```

This will start a server which will listen on port `3978`

#### Step 3: Expose local server with ngrok

Now, in separate terminal, run [ngrok](https://ngrok.com/download) command to expose your local server to outside world.

```sh
$ ngrok http 3978
```

Copy `https` endpoint, go to [Bot Framework](https://dev.botframework.com/bots) dashboard and set messaging endpoint under Settings.

#### Step 4: Test the BOT

You can either test BOT on BotFramework portal or you can create app manifest and install the App on Teams as mentioned [here](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/create-a-bot-for-teams#create-your-app-manifest-and-package).


## Understanding the example

The program starts by creating a handler structs of type `activity.HandlerFuncs`.

This example contains following handlers:

`OnMessageFunc` - sends upload request to user using file `consent card`

`OnInvokeFunc` - parses `invoke` activity, uploads file and notifies user with file `info card`


The `init` function picks up the `APP_ID` and `APP_PASSWORD` from the environment session and creates an `adapter` using this.

The consent cleanup `cleanupConsents()` is started as a go routine to delete processed consent cards.

A webserver is started with a handler which passes the received payload to `adapter.ParseRequest`. This methods authenticates the payload, parses the request and returns an Activity value.

```
activity, err := adapter.ParseRequest(ctx, req)
```
  

The Activity is then passed to `adapter.ProcessActivity` with the handler created to process the activity as per the handler functions and send the response to the connector service.

```
err = adapter.ProcessActivity(ctx, activity, customHandler)
```

Once the incoming activity is processed, file `consent card` is sent to ask user for file upload permission

e.g

```
{
    "contentType": "application/vnd.microsoft.teams.card.file.consent",
    "name": "data.txt",
    "content": {
        "description": "Sample data",
        "sizeInBytes": 29,
        "acceptContext": {
            "resultId": <unique-id>
        },
        "declineContext": {
            "resultId": <unique-id>
        }
    }
```

Based on user's response, `invoke` activity is sent to Bot, which we handle in `OnInvokeFunc`

In `OnInvokeFunc`, we find `uploadInfo` from the `invoke` activity which has following fields:

```
    "type": "invoke",
    "name": "fileConsent/invoke",
    ...
    "value": {
        "type": "fileUpload",
        "action": "accept",
        "context": {
            "resultId": <unique-id>
        },
        "uploadInfo": {
            "contentUrl": "https://<onedrive_url>/data.txt",
            "name": "data.txt",
            "uploadUrl": "https://<onedrive_upload_url>",
            "uniqueId": <unique-id>,
            "fileType": "txt"
        }
```

Then, the file contents are uploaded on `uploadUrl` using `PUT` request

```
err = putRequest(uploadInfo.UploadURL, data)
```

Finally, user is notified with the file info using `file info card`

```
{
    "contentType": "application/vnd.microsoft.teams.card.file.info",
    "contentUrl": "<uploadInfo.contentUrl>",
    "name": "<uploadInfo.name>",
    "content": {
        "uniqueId": "<uploadInfo.uniqueId>",
        "fileType": "<uploadInfo.fileType>",
    }
}
```

In case of no error, this web server responds with a 200 status.
