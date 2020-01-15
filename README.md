# Microsoft Bot Framework SDK for Go

[![Build Status](https://travis-ci.org/infracloudio/msbotbuilder-go.svg?branch=develop)](https://travis-ci.org/infracloudio/msbotbuilder-go) [![GoDoc](https://godoc.org/github.com/infracloudio/msbotbuilder-go?status.svg)](https://godoc.org/github.com/infracloudio/msbotbuilder-go)

This repository is the Go version of the Microsoft Bot Framework SDK. It facilitates developers to build bot applications using the Go language.

## Installing

```sh
$ go get -u github.com/infracloudio/msbotbuilder-go/...
```

## Get started with example

The [samples](samples/echobot) contains a sample bot created using this library which echoes any message received.

Before running this, two environment variables are needed viz. the Bot Framework application ID and the password. This can be received after [registration of a new bot](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/create-a-bot-for-teams#register-your-web-service-with-the-bot-framework).

```sh
$ export APP_ID=MICROSOFT_APP_ID
$ export APP_PASSWORD=MICROSOFT_APP_PASSWORD
```

Then, from the root of this repository,

```sh
$ cd samples/echobot
$ go run main.go
```

This starts a webserver on port `3978` by default.

This is the endpoint which the connector service for the registered bot should point to. For a descriptive understanding of the example refer the [sample](samples/).

## Contributing

We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:
- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features

## Credits

This project is highly inspired from the official Microsoft Bot Framework SDK - https://github.com/microsoft/botbuilder-python.

We have borrowed most of the design principles from the official Python SDKs.
