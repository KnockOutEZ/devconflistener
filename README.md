DevConfBD Ticket Watcher

This repository contains a Go program that continuously monitors the DevConfBD website for the availability of tickets. When tickets become available, the program will send a notification to your phone via Whatsapp.


Getting Started

To get started, you will need to:

1. Install Go.
2. Clone this repository.
3. Create a .env file in the root directory of the repository and add the following environment variables:
```
TWILIO_ACCOUNT_SID=
TWILIO_AUTH_TOKEN=
SENDER_PHONE=
RECIEVER_PHONE=
```

Code snippet


4. Run the program with the following command:
```
go run main.go
```

The program will now continuously monitor the DevConfBD website for the availability of tickets. When tickets become available, the program will send a notification to your phone via Whatsapp.


Features

The DevConfBD Ticket Listener has the following features:
    Continuously monitors the DevConfBD website for the availability of tickets.
    Sends a notification to your phone via SMS when tickets become available.
    Supports WhatsApp notifications.
    Easy to configure.

Note: It is a really really small script to get things done real quick so it doesnt follow the idiomatic Go folder structure or code patterns. You can follow my other applications in Nexentra(https://github.com/orgs/Nexentra/repositories) for more idiomatic patterns.