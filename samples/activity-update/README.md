# Bot Framework activity update sample.

This Microsoft Teams bot uses the [msbotbuilder-go](https://github.com/infracloudio/msbotbuilder-go) library. It shows how to create a simple bot that accepts input from the user and echoes response with an attachment.

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

The program starts by creating a handler struct of type `activity.HandlerFuncs`.

This struct contains a definition for the `OnMessageFunc` field which is treated as a callback by the library on the respective event.
In this sample, when the GetCard command is issued from the chat application, an adaptive card with an input box is sent from this bot. Once submit button on the card is clicked the card is deleted. Also, the bot updates the activity and sets the text message to activity as "Changed Activity" after the card is deleted.

```bash
var customHandler = activity.HandlerFuncs{
	OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
		if turn.Activity.Text == "getCard" {
			sJson := `{
				"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
				"type": "AdaptiveCard",
				"version": "1.0",
				"body": [
				  {
					"type": "TextBlock",
					"text": "Sample"
				  },
				  {  
					"type": "Input.Text",  
					"id": "firstName",  
					"placeholder": "What is your first name?"  
				  }
				],
				"actions": [
				  {
					"type": "Action.Submit",
					"title": "Submit"
				  }
				]
			  }`
			var obj map[string]interface{}
			err := json.Unmarshal(([]byte(sJson)), &obj)
			if err != nil {
				return schema.Activity{}, nil
			}
			attachments := []schema.Attachment{
				{
					ContentType: "application/vnd.microsoft.card.adaptive",
					Content:     obj,
				},
			}
			return turn.SendActivity(activity.MsgOptionAttachments(attachments))
		}
		if turn.Activity.Value != nil {
			fmt.Println("Activity=", turn.Activity.Value)
			activityID = turn.Activity.ReplyToID
		}

		return turn.SendActivity(activity.MsgOptionText("Echo: " + turn.Activity.Text))
	},
}
```
  

The `init` function picks up the `APP_ID` and `APP_PASSWORD` from the environment session and creates an `adapter` using this.


A webserver is started with a handler which passes the received payload to `adapter.ParseRequest`. This methods authenticates the payload, parses the request and returns an Activity value.

```
activity, err := adapter.ParseRequest(ctx, req)
```
  

The Activity is then passed to `adapter.ProcessActivity` with the handler created to process the activity as per the handler functions and send the response to the connector service.

```
err = adapter.ProcessActivity(ctx, activity, customHandler)
```

In case of no error, this web server responds with a 200 status.
