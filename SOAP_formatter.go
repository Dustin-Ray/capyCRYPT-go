package main

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

// refactor
const (
	soapMessageBegin = "-----------BEGIN-SOAP-MESSAGE-----------\n"
	soapMessageEnd   = "------------END-SOAP-MESSAGE------------"
	signatureBegin   = "----------BEGIN-SOAP-SIGNATURE----------\n"
	signatureEnd     = "-----------END-SOAP-SIGNATURE-----------"
)

// Formats a given message to SOAP format as specified in docs
func getSOAP(message *string, ctx *WindowCtx, l1, l2 string) *string {

	// string length
	sl := utf8.RuneCountInString(*message)
	// set line length
	lineLength := 40

	var sb strings.Builder
	sb.Write([]byte(l1))
	// Show a progress bar
	ctx.updateStatus("Working...")
	ctx.fixed.Put(ctx.progressBar, 245, 530)
	ctx.progressBar.SetVisible(true)
	ctx.progressBar.SetFraction(0)
	// loop through string
	for i := 0; i < len(*message); i += lineLength {
		if !ctx.initialState {
			// if line length is more than remaining characters
			if i+lineLength > sl {
				line := []byte((*message)[i:] + "\n")
				sb.Write(line)
			} else {
				line := []byte((*message)[i:i+lineLength] + "\n")
				sb.Write(line)
			}
			ctx.progressBar.SetFraction(float64(i) / float64(len(*message)))
			refreshWindow()
		} else {
			ctx.Reset()
			break
		}
	}
	sb.Write([]byte(l2))
	ctx.fixed.Remove(ctx.progressBar)
	ctx.progressBar.SetVisible(false)
	res := sb.String()
	return &res
}

// Parses a SOAP formatted string by removing header and footer and all newlines.
func parseSOAP(message *string, l1, l2 string) (*[]byte, error) {
	lines := strings.Split(*message, "\n")
	if lines[0]+"\n" == l1 &&
		lines[len(lines)-1] == l2 {
		str := strings.Join(lines[1:len(lines)-1], "\n")
		regex := regexp.MustCompile(`[\r\n]+`)
		strippedString := regex.ReplaceAllString(str, "")
		res, _ := hex.DecodeString(strippedString)
		return &res, nil
	} else {
		return nil, errors.New("malformed text, unable to parse")
	}
}
