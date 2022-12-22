package main

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

func getSOAPMessage(message string, ctx *WindowCtx) string {

	// string length
	sl := utf8.RuneCountInString(message)

	// set line length
	ll := 40

	strL := "-----------BEGIN-SOAP-MESSAGE-----------\n"
	// loop through string
	ctx.updateStatus("Working...")
	ctx.fixed.Put(ctx.progressBar, 245, 530)
	ctx.progressBar.SetVisible(true)
	ctx.progressBar.SetFraction(0)
	for i := 0; i < len(message); i += ll {
		if !ctx.initialState {
			// if line length is more than remaining characters
			if i+ll > sl {
				strL = strL + message[i:] + "\n"
			} else {
				strL = strL + message[i:i+ll] + "\n"
			}
			ctx.progressBar.SetFraction(float64(i) / float64(len(message)))
			refreshWindow()
		} else {
			ctx.Reset()
			break
		}
	}
	strL += "------------END-SOAP-MESSAGE------------"
	ctx.fixed.Remove(ctx.progressBar)
	ctx.progressBar.SetVisible(false)
	return strL
}

// Parses a SOAP formatted string by removing header and footer and all newlines.
func parseSOAPMessage(message string) (string, error) {
	lines := strings.Split(message, "\n")

	if lines[0] == "-----------BEGIN-SOAP-MESSAGE-----------" &&
		lines[len(lines)-1] == "------------END-SOAP-MESSAGE------------" {

		str := strings.Join(lines[1:len(lines)-1], "\n")
		regex := regexp.MustCompile(`[\r\n]+`)
		strippedString := regex.ReplaceAllString(str, "")
		return strippedString, nil
	} else {
		return "test", errors.New("malformed cryptogram, unable to decrypt")
	}

}
