package main

import (
	"github.com/gotk3/gotk3/gtk"
)

//A factory that produces the buttons and corresponding signals to be used in the view.

// adds buttons in a factory style to fixed context
func setupButtons(ctx *WindowCtx) {

	labelList := []string{"Compute Hash", "Compute Tag", "Encrypt With Password", "Decrypt With Password",
		"Generate Keypair", "Encrypt With Key", "Decrypt With Key", "Sign With Key", "Verify Signature"}

	buttonList := make([]gtk.Button, len(labelList))

	for i, label := range labelList {
		btn, _ := gtk.ButtonNewWithLabel(label)
		buttonList[i] = *btn
		ctx.fixed.Put(btn, 40, 80+i*45)
	}

	buttonList[0].SetTooltipMarkup("Computes a SHA3-512 hash of the text in the notepad.")
	buttonList[0].Connect("clicked", func() {
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		ctx.notePad.SetText(ComputeSHA3HASH(text, ctx.fileMode))
		ctx.updateStatus("SHA3 hash computed successfully")
	}) //etc....

	buttonList[1].SetTooltipMarkup("Computes a keyed hash of the notepad. The resulting hash can only be computed parties with knowledge of the password and the message.")
	buttonList[1].Connect("clicked", func() {
		password := showPasswordDialog(ctx.win, "authentication")
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		ctx.notePad.SetText(ComputeTaggedHash(password, []byte(text), "T"))
		ctx.updateStatus("Message tag computed successfully")
	}) //etc....

	buttonList[2].SetTooltipMarkup("Encrypts data under a passphrase. Can only be decrypted by parties with knowledge of the passphrase.")
	buttonList[2].Connect("clicked", func() {

		password := showPasswordDialog(ctx.win, "encryption")
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		ctx.notePad.SetText(BytesToHexString(encryptPW(password, []byte(text))))
		ctx.updateStatus("Encryption successful")
	}) //etc....

	buttonList[3].SetTooltipMarkup("Decrypts data under a passphrase. Can only be decrypted by parties with knowledge of the passphrase.")
	buttonList[3].Connect("clicked", func() {

		password := showPasswordDialog(ctx.win, "encryption")
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		ctx.notePad.SetText(string(decryptPW(password, HexToBytes(text))))
		ctx.updateStatus("Decryption successful")
	}) //etc....

	buttonList[4].SetTooltipMarkup("Generates a Schnorr E521 keypair from supplied password.")
	buttonList[4].Connect("clicked", func() {
		key := KeyObj{}
		generateKeyPair(ctx, &key)
		ctx.keytable.importKey(ctx, key)
	}) //etc....

	reset, _ := gtk.ButtonNewWithLabel("Reset")
	reset.SetName("resetButton") //for CSS styling
	reset.Connect("button-press-event", func() {
		ctx.Reset()
	})
	ctx.fixed.Put(reset, 40, 510)
}
