package main

import "github.com/gotk3/gotk3/gtk"

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
		ctx.notePad.SetText(ComputeSHA3HASH(text))
		ctx.updateStatus("SHA3 hash computed successfully")
	}) //etc....

	buttonList[1].SetTooltipMarkup("Computes a keyed hash of the notepad. The resulting hash can only be computed by a party who has knowledge of the password and the message.")
	buttonList[1].Connect("clicked", func() {
		password := showPasswordDialog(ctx.win, "authentication")
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		ctx.notePad.SetText(ComputeTaggedHash(password, []byte(text), "T"))
		ctx.updateStatus("Message tag computed successfully")
	}) //etc....

	reset, _ := gtk.ButtonNewWithLabel("Reset")
	reset.SetName("resetButton") //for CSS styling
	reset.Connect("button-press-event", func() {
		ctx.Reset()
	})
	ctx.fixed.Put(reset, 40, 510)
}
