package main

//Button manufacturing facility

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
			cg := encryptPW([]byte(password), &textBytes)
			temp := hex.EncodeToString(*cg)
			ctx.notePad.SetText(getSOAPMessage(temp, ctx))
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
			psdMsg, err1 := parseSOAPMessage(&notePadText)
			if err1 != nil {
				ctx.updateStatus("unable to decrypt")
			} else {
				cg, _ := parseCryptogram(psdMsg)
				dec, err := decryptPW([]byte(password), cg)
				if err != nil {
					ctx.updateStatus(err.Error())
				} else {
					ctx.notePad.SetText(string(*dec))
				}
			}
		} else {
			ctx.updateStatus("decryption cancelled")
		}
	})

	//Keygen
	buttonList[4].SetTooltipMarkup("Generates a Schnorr E521 keypair from supplied password.")
	buttonList[4].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		key := KeyObj{}
		result := generateKeyPair(ctx, &key)
		if result {
			ctx.keytable.importKey(ctx, key)
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

			result := encryptKey(key, &textBytes)
			ctx.notePad.SetText(getSOAPMessage(hex.EncodeToString(*result), ctx))

			ctx.toggleButtons(ctx.buttons, true)
			ctx.updateStatus("encryption successful")
		} else {
			ctx.updateStatus("encryption cancelled")
		}
	})

	// EC decryption. Searches keytable for corresponding private key.
	buttonList[6].SetTooltipMarkup("Decrypts data using passphrase that corresponds to a valid private key.")
	buttonList[6].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		password, result := passwordEntryDialog(ctx.win, "decryption")
		if result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			text2, _ := parseSOAPMessage(&text)
			psdMsg, err := parseCryptogram(text2)
			if err != nil {
				ctx.updateStatus(err.Error())
			} else {
				message, _ := decryptKey([]byte(password), psdMsg)
				ctx.notePad.SetText(*message)
				ctx.updateStatus("decryption successful")
			}
		} else {
			ctx.updateStatus("decryption cancelled")
		}
	})
	return &buttonList
}

// Disables buttons while operation is being performed, reenables buttons when finished.
// Reset toggles buttons on
func (ctx *WindowCtx) toggleButtons(buttonList *[]gtk.Button, setting bool) {
	for _, button := range *buttonList {
		button.SetSensitive(setting)
	}
}
