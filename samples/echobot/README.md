# Bot Framework echo bot sample.

This Microsoft Teams bot uses the [msbotbuilder-go](https://github.com/infracloudio/msbotbuilder-go) library. It shows how to create a simple bot that accepts input from the user and echoes it back.

## Run the example

#### Step 1: Register MS BOT

Follow the official [documentation](https://docs.microsoft.com/en-us/azure/bot-service/abs-quickstart) to create and register your BOT.

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

Copy `https` endpoint, go to [Azure Portal](https://portal.azure.com) dashboard, select your bot under `Resources`, click `Configuration`, and paste the URL into `Messaging endpoint`. Append `/api/messages` to the URL, i.e. it should look like `https://b7fd-92-6-238-167.eu.ngrok.io/api/messages`.

#### Step 4: Test the BOT

Test the bot by going to [Azure Portal](https://portal.azure.com), select your bot under `Resources`, and click `Test in Web Chat` (under the `Settings` header). Note that you will need to update the `Messaging endpoint` URL in the bot configuration whenever you restart ngrok.


## Understanding the example

The program starts by creating a handler struct of type `activity.HandlerFuncs`.

This struct contains definition for the `OnMessageFunc` field which is a treated as a callback by the library on the respective event.

```bash
var customHandler = activity.HandlerFuncs{
	OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
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
