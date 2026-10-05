package main

import (
	"fmt"
	"log"
)

func main() {
	srv, err := getService()
	if err != nil {
		log.Fatalf(err.Error())
	}

	user := "me"
	// get first 10 messages
	resp, err := srv.Users.Messages.List(user).MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	for _, m := range resp.Messages {
		msg, err := srv.Users.Messages.Get(user, m.Id).Format("full").Do()
		if err != nil {
			fmt.Printf("Could not retrieve message. continuing... : %v", err)
			continue
		}

		// print message headers
		var subject string
		var from string
		var date string
		var unsub string

		for _, h := range msg.Payload.Headers {
			switch h.Name {
			case "Subject":
				subject = fmt.Sprint("Subject:", h.Value)
			case "From":
				from = fmt.Sprint("From:", h.Value)
			case "Date":
				date = fmt.Sprint("Date:", h.Value)
			case "List-Unsubscribe":
				unsub = fmt.Sprint("Unsub:", h.Value)
			}
		}

		fmt.Println(from)
		fmt.Println(subject)
		fmt.Println(date)
		fmt.Println(unsub)
		fmt.Println()
	}
}
