# msbotbuilder-go

#### Bot Framework SDK for GoLang

This repository is the Go version of the Microsoft Bot Framework SDK. It facilitates developers to build bot applications using the Go language .

#### Get started with example

The `samples\echobot\` contains a sample bot created using thie library which echoes any message received.

Before running this, two environment variables are needed viz. the Bot Framework application ID and the password. This can be received after [registration of a new bot](https://dev.botframework.com/).

```
export APP_ID=MICROSOFT_APP_ID
export APP_PASSWORD=MICROSOFT_APP_PASSWORD
```

Then, from the root of this repository,

```
cd sameple/echobot/
go run main.go
```

This starts a webserver on port 3978 by default.

This is the endpoint which the connector service for the registered bot should point to. For a descriptive understanding of the example refer the [doc](https://github.com/infracloudio/msbotbuilder-go/blob/golint/samples/echobot/doc.go)