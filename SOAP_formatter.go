package main

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Formats a given message to SOAP format as specified in docs
func getSOAPMessage(message string, ctx *WindowCtx) string {

	// string length
	sl := utf8.RuneCountInString(message)
	// set line length
	ll := 40

	var sb strings.Builder
	sb.Write([]byte("-----------BEGIN-SOAP-MESSAGE-----------\n"))
	// Show a progress bar
	ctx.updateStatus("Working...")
	ctx.fixed.Put(ctx.progressBar, 245, 530)
	ctx.progressBar.SetVisible(true)
	ctx.progressBar.SetFraction(0)
	// loop through string
	for i := 0; i < len(message); i += ll {
		if !ctx.initialState {
			// if line length is more than remaining characters
			if i+ll > sl {
				sb.Write([]byte(message[i:] + "\n"))
			} else {
				sb.Write([]byte(message[i:i+ll] + "\n"))
			}
			ctx.progressBar.SetFraction(float64(i) / float64(len(message)))
			refreshWindow()
		} else {
			ctx.Reset()
			break
		}
	}
	sb.Write([]byte("------------END-SOAP-MESSAGE------------"))
	ctx.fixed.Remove(ctx.progressBar)
	ctx.progressBar.SetVisible(false)
	return sb.String()
}

// Parses a SOAP formatted string by removing header and footer and all newlines.
func parseSOAPMessage(message *string) (*[]byte, error) {
	lines := strings.Split(*message, "\n")
	if lines[0] == "-----------BEGIN-SOAP-MESSAGE-----------" &&
		lines[len(lines)-1] == "------------END-SOAP-MESSAGE------------" {
		str := strings.Join(lines[1:len(lines)-1], "\n")
		regex := regexp.MustCompile(`[\r\n]+`)
		strippedString := regex.ReplaceAllString(str, "")
		res, _ := hex.DecodeString(strippedString)
		return &res, nil
	} else {
		return nil, errors.New("malformed cryptogram, unable to decrypt")
	}
}
