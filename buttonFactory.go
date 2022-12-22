package main

//Button manufacturing facility

import (
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
			ctx.updateStatus("Message tag computed successfully")
		} else {
			ctx.updateStatus("Tag computation cancelled")
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
			ctx.notePad.SetText(getSOAPMessage(BytesToHexString(encryptPW([]byte(password), &textBytes)), ctx))
			ctx.toggleButtons(ctx.buttons, true)
			ctx.updateStatus("Encryption successful")
		} else {
			ctx.updateStatus("Encryption cancelled")
		}
	}) //etc....

	//Symmettric Decryption. Emits ambiguous errors for security of password.
	buttonList[3].SetTooltipMarkup("Decrypts data under a passphrase. Can only be decrypted by parties with knowledge of the passphrase.")
	buttonList[3].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		password, result := passwordEntryDialog(ctx.win, "decryption")
		if result {
			text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
			//1st check: a valid message must be at least 128 bytes long
			if len(text) > 128 {
				fmttedMsg, err1 := parseSOAPMessage(text)
				//2nd check: a message must be formatted as SOAP
				if err1 != nil {
					ctx.updateStatus("error received: " + err1.Error())
				} else {
					message, err := decryptPW([]byte(password), HexToBytes(fmttedMsg))
					//3rd check: decryption password must be valid
					if err != nil {
						ctx.updateStatus("error received: " + err.Error())
					}
					ctx.notePad.SetText(string(message))
					ctx.updateStatus("decryption successful")
				}
			} else {
				ctx.updateStatus("malformed cryptogram, unable to decrypt")
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
			ctx.updateStatus("Key generation cancelled")
		}
	})

	// Keygen
	buttonList[5].SetTooltipMarkup("Encrypts data using a public key selected from the key table.")
	buttonList[5].Connect("clicked", func() {
		ctx.initialState = false
		ctx.fileMode = false
		key := ctx.loadedKey
		if key != nil {
			result := generateKeyPair(ctx, ctx.loadedKey)
			if result {
				ctx.keytable.importKey(ctx, *key)
			} else {
				ctx.updateStatus("Key generation cancelled")
			}
		} else {
			ctx.updateStatus("No key selected")
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
