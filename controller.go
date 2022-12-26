package main

/* Controller for MVC. Connects buttons to model functionality and
transmits messages to view from model. */

import (
	"encoding/hex"
	"math/big"

	"github.com/gotk3/gotk3/gtk"
)

// adds buttons in a factory style to fixed context
func setupButtons(ctx *WindowCtx) *[]gtk.Button {

	labelList := []string{"Compute Hash", "Compute Tag", "Encrypt With Password", "Decrypt With Password",
		"Generate Keypair", "Encrypt With Key", "Decrypt With Key", "Sign With Key", "Verify Signature"}

	buttonList := make([]gtk.Button, len(labelList))
	for i, label := range labelList {
		btn, _ := gtk.ButtonNewWithLabel(label)
		buttonList[i] = *btn
		ctx.fixed.Put(btn, 40, 80+i*45)
	}

	//Reset session
	reset, _ := gtk.ButtonNewWithLabel("Reset")
	reset.SetName("resetButton") //for CSS styling
	reset.Connect("clicked", func() {
		showRestWarningDialog(ctx)
	})
	ctx.fixed.Put(reset, 40, 510)

	/* BUTTON CONSTRUCTION:*/

	//SHA3 hash
	buttonList[0].SetTooltipMarkup("Computes a SHA3-512 hash of the text in the notepad.")
	buttonList[0].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		ctx.notePad.SetText(ComputeSHA3HASH(text, ctx.fileMode))
		ctx.updateStatus("SHA3 hash computed successfully")
	})

	//HMAC
	buttonList[1].SetTooltipMarkup("Computes a keyed hash of the notepad. The resulting hash could only have been computed by parties with knowledge of the password and the message.")
	buttonList[1].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		password, result := passwordEntryDialog(ctx.win, "authentication")
		if result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			ctx.notePad.SetText(ComputeTaggedHash([]byte(password), []byte(text), "T"))
			ctx.updateStatus("message tag computed successfully")
		} else {
			ctx.updateStatus("tag computation cancelled")
		}
	})

	//Symmetric encryption
	buttonList[2].SetTooltipMarkup("Encrypts data under a passphrase. Can only be encrypted by parties with knowledge of the passphrase.")
	buttonList[2].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		password, result := passwordEntryDialog(ctx.win, "encryption")
		if result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), false)
			textBytes := []byte(text)
			ctx.toggleButtons(ctx.buttons, false)
			cg := encryptWithPW([]byte(password), &textBytes)
			temp := hex.EncodeToString(*cg)
			res := getSOAP(&temp, ctx, soapMessageBegin, soapMessageEnd)
			ctx.notePad.SetText(*res)
			ctx.toggleButtons(ctx.buttons, true)
			ctx.updateStatus("encryption successful")
		} else {
			ctx.updateStatus("encryption cancelled")
		}
	}) //etc....

	//Symmetric Decryption. Emits ambiguous errors for security of password.
	buttonList[3].SetTooltipMarkup("Decrypts data under a passphrase. Can only be decrypted by parties with knowledge of the passphrase.")
	buttonList[3].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		password, result := passwordEntryDialog(ctx.win, "decryption")
		if result {
			notePadText, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			psdMsg, err1 := parseSOAP(&notePadText, soapMessageBegin, soapMessageEnd)
			if err1 != nil {
				ctx.updateStatus("unable to decrypt")
			} else {
				cg, _ := decodeSymCryptogram(psdMsg)
				dec, err := decryptWithPW([]byte(password), cg)
				if err != nil {
					ctx.updateStatus(err.Error())
				} else {
					ctx.notePad.SetText(string(*dec))
					ctx.updateStatus("decryption successful")
				}
			}
		} else {
			ctx.updateStatus("decryption cancelled")
		}
	})

	//Generate a keypair
	buttonList[4].SetTooltipMarkup("Generates a Schnorr E521 keypair from supplied password.")
	buttonList[4].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		key := KeyObj{}
		opResult := constructKey(ctx, &key)
		if opResult {
			ctx.keytable.importKey(ctx, key)
			ctx.updateStatus("key " + key.Id + " generated successfully")
		} else {
			ctx.updateStatus("key generation cancelled")
		}
	})

	// Encrypt with EC public key
	buttonList[5].SetTooltipMarkup("Encrypts data using a public key selected from the key table.")
	buttonList[5].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		if ctx.loadedKey != nil {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), false)
			textBytes := []byte(text)
			ctx.toggleButtons(ctx.buttons, false)
			//construct the key
			pubX, _ := big.NewInt(0).SetString(ctx.loadedKey.PubKeyX, 10)
			pubY, _ := big.NewInt(0).SetString(ctx.loadedKey.PubKeyY, 10)
			key := NewE521XY(*pubX, *pubY)

			result := hex.EncodeToString(*encryptWithKey(key, &textBytes))
			res := getSOAP(&result, ctx, soapMessageBegin, soapMessageEnd)
			ctx.notePad.SetText(*res)

			ctx.toggleButtons(ctx.buttons, true)
			ctx.updateStatus("encryption successful")
		} else {
			ctx.updateStatus("encryption cancelled")
		}
	})

	// EC decryption. TODO Searches keytable for corresponding private key.
	buttonList[6].SetTooltipMarkup("Decrypts data using passphrase that corresponds to a valid private key.")
	buttonList[6].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		password, result := passwordEntryDialog(ctx.win, "decryption")
		if password != "" && result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			text2, err := parseSOAP(&text, soapMessageBegin, soapMessageEnd)
			if err != nil {
				ctx.updateStatus(err.Error())
			} else {
				psdMsg, err := decodeECCryptogram(text2)
				if err != nil {
					ctx.updateStatus(err.Error())
				} else {
					message, err := decryptWithKey([]byte(password), psdMsg)
					if err != nil {
						ctx.updateStatus(err.Error())
					} else {
						ctx.notePad.SetText(*message)
						ctx.updateStatus("decryption successful")
					}
				}
			}
		} else {
			ctx.updateStatus("decryption cancelled")
		}
	})

	//Signs a message using a private key derived from a password.
	buttonList[7].SetTooltipMarkup("Signs a message with a selected key.")
	buttonList[7].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		password, result := passwordEntryDialog(ctx.win, "signature")
		if result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			textBytes := []byte(text)
			signature, err := signWithKey([]byte(password), &textBytes)
			if err != nil {
				ctx.updateStatus(err.Error())
			} else {
				sigHexString := hex.EncodeToString(*signature)
				soapFmttedSig := getSOAP(&sigHexString, ctx, signatureBegin, signatureEnd) //refactor
				ctx.notePad.SetText(*soapFmttedSig)
				ctx.updateStatus("signature generated")
			}
		} else {
			ctx.updateStatus("signature cancelled")
		}
	})

	//Verifies signature using public key.
	buttonList[8].SetTooltipMarkup("Verifies a signature against a public key.")
	buttonList[8].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		if ctx.loadedKey != nil {
			pubKeyObj := ctx.loadedKey                                //loaded key should maybe be keyoobj with E521 for public key instead of x/y
			keyX, _ := big.NewInt(0).SetString(pubKeyObj.PubKeyX, 10) //refactor
			keyY, _ := big.NewInt(0).SetString(pubKeyObj.PubKeyY, 10) //refactor
			key := NewE521XY(*keyX, *keyY)
			signatureBytes, err := parseSOAP(&text, signatureBegin, signatureEnd)
			if err != nil {
				ctx.updateStatus("error parsing signature")
			} else {
				signature, err2 := decodeSignature(signatureBytes)
				if err != nil || err2 != nil {
					ctx.updateStatus("unable to parse signature")
				} else {
					result := verify(key, signature, &signature.M)
					if result {
						ctx.updateStatus("good signature from key " + ctx.loadedKey.Id)
					} else {
						ctx.updateStatus("unable to verify signature")
					}
				}
			}
		} else {
			ctx.updateStatus("no key selected")
		}
	})
	return &buttonList
}
