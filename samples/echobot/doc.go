/*Bot Framework echo bot sample.
This bot uses msbotbuilder-go: https://github.com/infracloudio/msbotbuilder-go. It shows
how to create a simple bot that accepts input from the user and echoes it back.


Run the example

Bring up a terminal and run. Set two variables for the session as APP_ID and APP_PASSWORD to the values of
your BotFramework app_id and password. Then, run:
		go run main.go
This will start a server which will listen on port 3978


Understanding the example

The program starts by creating a hanlder struct of type `activity.HandlerFuncs`.
This struct contains defination for the `OnMessageFunc` field which is a treated as a callback by the library
on the respective event.

	var customHandler = activity.HandlerFuncs{
		OnMessageFunc: func(turn *activity.TurnContext) (schema.Activity, error) {
			return turn.SendActivity(activity.MsgOptionText("Echo: " + turn.Activity.Text))
		},
	}


A webserver is started with a hanlder passed the received payload to `adapter.ParseRequest`
This methods authenticates the payload, parses the request and returns an Activity value.
	activity, err := adapter.ParseRequest(ctx, req)

The Activity is then passed to `adapter.ProcessActivity` with the hanlder created to process
the activity as per the hanlder functions and send the response to the connector service.
	err = adapter.ProcessActivity(ctx, activity, customHandler)

In case of no error, this web responds with a 200 status

To expose this local IP outside your local network, a tool like ngrok can be used.

	ngrok http 3978

The server is then available on a IP similar to http://92832de0.ngrok.io
*/
package main
