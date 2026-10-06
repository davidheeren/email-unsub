package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const emailUser = "me"

func main() {
	srv, err := getService()
	if err != nil {
		log.Fatalf(err.Error())
	}

	if len(os.Args) > 1 && os.Args[1] == "download" {
		if len(os.Args) < 3 {
			log.Fatalf("download command needs a email count")
		}
		emailDownloadCount, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf(err.Error())
		}
		err = downloadTestData(srv, emailDownloadCount)
		if err != nil {
			log.Fatalf(err.Error())
		}
		return
	}

	// get first 10 messages
	resp, err := srv.Users.Messages.List(emailUser).MaxResults(100).Do()
	if err != nil {
		log.Fatalf("unable to retrieve messages: %v", err)
	}

	for _, m := range resp.Messages {
		msg, err := srv.Users.Messages.Get(emailUser, m.Id).Format("full").Do()
		if err != nil {
			fmt.Printf("could not retrieve message. continuing... : %v", err)
			continue
		}

		unsubUrl, err := scanUnsubUrls(srv, msg)
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		if (unsubUrl == "") {
			continue
		}

		// print message headers
		var subject string
		var from string
		var date string

		for _, h := range msg.Payload.Headers {
			switch h.Name {
			case "Subject":
				subject = fmt.Sprint("Subject:", h.Value)
			case "From":
				from = fmt.Sprint("From:", h.Value)
			case "Date":
				date = fmt.Sprint("Date:", h.Value)
			}
		}

		fmt.Println(from)
		fmt.Println(subject)
		fmt.Println(date)
		fmt.Printf("Unsub: %s\n", unsubUrl)
		fmt.Println()
	}
}
