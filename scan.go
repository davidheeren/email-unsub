package main

import (
	"errors"
	"regexp"
	"strings"

	"google.golang.org/api/gmail/v1"
)

// finds the most likely unsub url ("" if none)
func scanUnsubUrls(msg *gmail.Message) (string, error) {
	if msg.Raw != "" {
		return "", errors.New("message has to be in format 'full', not 'raw'")
	}

	var unsubUrl string
	// we are only getting the url from this header for now
	// TODO: parse the body for urls
	for _, h := range msg.Payload.Headers {
		// fmt.Println(h.Name)
		if h.Name == "List-Unsubscribe" {
			// get url in <>
			re := regexp.MustCompile("<http.*?>")
			unsubUrl = re.FindString(h.Value)
			// strip <>
			unsubUrl = strings.TrimLeft(unsubUrl, "<")
			unsubUrl = strings.TrimRight(unsubUrl, ">")
		}
	}

	return unsubUrl, nil
}
