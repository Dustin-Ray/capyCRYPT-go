package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"regexp"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

/*
Displays a password entry dialog that asks a user for a passphrase.
Any input is considered a valid passphrase including the empty string.
OK button is disabled if content of password entry fields do not match.
*/
func passwordEntryDialog(parent *gtk.Window, message string) (string, bool) {
	// Create a dialog
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle("Enter " + message + " password:")

	okButton, _ := dialog.AddButton("OK", gtk.RESPONSE_OK)
	okButton.SetSensitive(true)

	// Create a horizontal box
	vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	hBox1, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	hBox2, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	hBox4, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	cA, _ := dialog.GetContentArea()
	vBox.Add(hBox1)
	vBox.Add(hBox2)
	vBox.Add(hBox4)
	cA.Add(vBox)

	// Create a label
	lbl, _ := gtk.LabelNew("Password: ")
	conf, _ := gtk.LabelNew("Confirm:     ")

	hBox1.Add(lbl)
	hBox2.Add(conf)

	// Create a password entry
	entry, _ := gtk.EntryNew()
	confirm, _ := gtk.EntryNew()
	entry.SetVisibility(false)
	confirm.SetVisibility(false)
	hBox1.Add(entry)
	hBox2.Add(confirm)

	matched := true

	confirm.Connect("changed", func() { changed(entry, confirm, &matched, okButton) })
	entry.Connect("changed", func() { changed(entry, confirm, &matched, okButton) })

	// Show the dialog
	dialog.ShowAll()
	if dialog.Run() == gtk.RESPONSE_OK {
		password, _ := confirm.GetText()
		dialog.Destroy()
		return password, true
	}
	// close the dialog
	dialog.Destroy()
	return "", false
}

func constructKey(ctx *WindowCtx, key *KeyObj) bool {

	// Create a dialog
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle("Create new key: ")

	okButton, _ := dialog.AddButton("OK", gtk.RESPONSE_OK)
	okButton.SetSensitive(false)

	//boolean to check if pwds match
	matched := true
	//Signals whether operation was completed or cancelled
	opResult := false

	// Create a boxes to store entry fields and labels
	vBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	hBox1, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	hBox2, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	hBox3, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	hBox4, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	cA, _ := dialog.GetContentArea()
	vBox.Add(hBox1)
	vBox.Add(hBox2)
	vBox.Add(hBox3)
	vBox.Add(hBox4)
	cA.Add(vBox)

	ownerLbl, _ := gtk.LabelNew("Owner:       ")
	lbl, _ := gtk.LabelNew("Password: ")
	conf, _ := gtk.LabelNew("Confirm:     ")
	owner, _ := gtk.EntryNew()
	entry, _ := gtk.EntryNew()
	confirm, _ := gtk.EntryNew()

	hBox1.Add(ownerLbl)
	hBox2.Add(lbl)
	hBox3.Add(conf)
	hBox1.Add(owner)
	hBox2.Add(entry)
	hBox3.Add(confirm)

	ownerLbl.SetTooltipMarkup("Owner of key can be a name, email address, or some other identifier. Cannot be blank.")

	entry.SetVisibility(false)
	confirm.SetVisibility(false)

	//generate a random key id using sponge
	r := generateRandomBytes(200)
	r = append(r, 0x18) //Delim Suffix for key ID
	key.Id = hex.EncodeToString(SpongeSqueeze(SpongeAbsorb(&r, 256), 48, 136))
	key.Owner = "NONE"
	key.KeyType = "PRIVATE"

	confirm.Connect("changed", func() { changed(entry, confirm, &matched, okButton) })
	entry.Connect("changed", func() { changed(entry, confirm, &matched, okButton) })

	// Show the dialog
	dialog.ShowAll()
	if dialog.Run() == gtk.RESPONSE_OK {
		if matched {
			ot, _ := owner.GetText()
			password2, _ := confirm.GetText()
			setKeyData(key, password2, ot)
			dialog.Destroy()
			return true
		}
	}
	// close the dialog
	dialog.Destroy()
	return opResult
}

func rightCLickMenu(ctx *WindowCtx) {
	// Create a Menu to display when right-clicking a row.
	menu, _ := gtk.MenuNew()
	// Create a MenuItem to be used in our Menu.
	details, _ := gtk.MenuItemNewWithLabel("Details...")
	pubKeyExport, _ := gtk.MenuItemNewWithLabel("Export public key...")
	privKeyExport, _ := gtk.MenuItemNewWithLabel("Export private keypair...")

	// Add the MenuItem to the Menu.
	menu.Append(details)
	menu.Append(pubKeyExport)
	menu.Append(privKeyExport)

	// Connect a signal handler to the MenuItem's "activate" signal.
	details.Connect("activate", func() { showKeyDetails(ctx) })
	pubKeyExport.Connect("activate", func() { exportPublicKey(ctx) })
	privKeyExport.Connect("activate", func() { exportPrivateKey(ctx) })
	menu.ShowAll()
	// Connect the "button-press-event" signal to our handler.
	ctx.keytable.treeview.Connect("button-press-event", func(treeView *gtk.TreeView, event *gdk.Event) {
		// Cast the Event to a GdkEventButton.
		eventButton := gdk.EventButtonNewFromEvent(event)
		// Check if the right mouse button was pressed.
		if eventButton.Button() == 3 && ctx.loadedKey != nil {
			// Show the Menu at the position of the mouse click.
			menu.PopupAtPointer(event)
		}
	})
}

// showDialog displays a new smaller blank window dialog with text in it.
func showKeyDetails(ctx *WindowCtx) {
	dialog, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	fixed, _ := gtk.FixedNew()

	dialog.SetPosition(gtk.WIN_POS_CENTER)
	dialog.SetDefaultSize(490, 365)
	dialog.SetTitle("Key details: ")

	// Create a label and set its text.
	key := ctx.loadedKey
	ID, _ := gtk.LabelNew("ID: " + key.Id)
	Owner, _ := gtk.LabelNew("OWNER: " + key.Owner)
	KEYTYPE, _ := gtk.LabelNew("KEY TYPE: " + key.KeyType)
	PubKeyX, _ := gtk.LabelNew("PUB KEY X: ")
	PubKeyY, _ := gtk.LabelNew("PUB KEY Y: ")
	PrivKey, _ := gtk.LabelNew("PRIV KEY: ")
	if ctx.loadedKey.KeyType != "NONE" {
		PrivKey, _ = gtk.LabelNew("PRIV KEY: " + ctx.loadedKey.PrivKey)
	}
	DateCreated, _ := gtk.LabelNew("DATE CREATED: " + key.DateCreated)

	dialog.Add(fixed)
	fixed.Put(ID, 25, 25)
	fixed.Put(Owner, 25, 60)
	fixed.Put(KEYTYPE, 25, 95)

	pubKeyXWindow := getScrollableTextArea(ctx, ctx.loadedKey.PubKeyX)
	fixed.Put(PubKeyX, 25, 130)
	fixed.Put(&pubKeyXWindow, 25, 155)

	pubKeyYWindow := getScrollableTextArea(ctx, ctx.loadedKey.PubKeyY)
	fixed.Put(PubKeyY, 250, 130)
	fixed.Put(&pubKeyYWindow, 250, 155)

	fixed.Put(PrivKey, 25, 270)
	fixed.Put(DateCreated, 25, 305)

	ID.SetSelectable(true)
	Owner.SetSelectable(true)
	KEYTYPE.SetSelectable(true)
	PrivKey.SetSelectable(true)
	DateCreated.SetSelectable(true)

	dialog.ShowAll()
}

// A dialog that exports key data to a file
func exportKeyDialog(ctx *WindowCtx, name string) {
	// Create a dialog that allows the user to save a file
	dialog, err := gtk.FileChooserDialogNewWith2Buttons("Export Key", ctx.win,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		"Cancel", gtk.RESPONSE_CANCEL,
		"Save", gtk.RESPONSE_ACCEPT)
	if err != nil {
		panic(err)
	}

	//enforce SOAP file format
	filter, _ := gtk.FileFilterNew()
	dialog.SetCurrentName(ctx.loadedKey.Id + ".SOAP_KEY")
	filter.AddPattern("*.SOAP_KEY")
	dialog.SetFilter(filter)

	// Show the dialog and wait for the user to respond
	response := dialog.Run()
	if response == gtk.RESPONSE_ACCEPT {
		// Get the filename from the dialog
		filename := dialog.GetFilename()
		jsonKeyData, _ := KeyToJSON(ctx.loadedKey)

		// Create the file
		file, err := os.Create(filename)
		if err != nil {
			ctx.updateStatus("Failed to create file")
			dialog.Destroy()
		}
		defer file.Close()
		file.Write(jsonKeyData)
		if err != nil {
			ctx.updateStatus("Failed to write key")
			dialog.Destroy()
		}

		ctx.status.SetLineWrapMode(pango.WRAP_CHAR)
		reg := regexp.MustCompile(`\/[^\/]*\.SOAP_KEY$`)
		filePath := reg.ReplaceAllString(filename, "")
		ctx.updateStatus("Key data saved to: " + filePath)
	}
	dialog.Destroy()
}

// A dialog that opens a key file. Handles any error in parsing file to key
func importKeyDialog(ctx *WindowCtx) {

	// Create a new FileChooserDialog to open a file
	fileDialog, err := gtk.FileChooserDialogNewWith2Buttons("Import Key", ctx.win,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"_Cancel", gtk.RESPONSE_CANCEL,
		"_Open", gtk.RESPONSE_ACCEPT)
	if err != nil {
		fmt.Println(err)
		return
	}

	//enforce SOAP file format
	filter, _ := gtk.FileFilterNew()
	filter.AddPattern("*.SOAP_KEY")
	fileDialog.SetFilter(filter)

	fileDialog.SetSizeRequest(200, 100)

	// Show the dialog and wait for the user's response.
	response := fileDialog.Run()
	if response == gtk.RESPONSE_ACCEPT {
		// If a file was selected, print out the name
		filename := fileDialog.GetFilename()
		err := ctx.keytable.JsonToKey(ctx, filename)
		if err != nil {
			ctx.updateStatus("Import failed - invalid key selected")
		}
	}
	// Destroy the dialog when done
	fileDialog.Destroy()
}

// Displays a warning dialog when the reset button is pressed
func showRestWarningDialog(ctx *WindowCtx) {
	dialog := gtk.MessageDialogNew(ctx.win,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_WARNING,
		gtk.BUTTONS_OK_CANCEL,
		"Clears all data and keys from this session, "+
			"cancels any running operations, "+
			"and resets the notepad. Exported keys are unaffected.")

	dialog.SetTitle("Danger Zone!")
	response := dialog.Run()
	if response == gtk.RESPONSE_OK {
		ctx.Reset()
		dialog.Destroy()
	} else {
		dialog.Destroy()
	}
}

func refreshWindow() {

	for {
		if gtk.EventsPending() {
			gtk.MainIterationDo(false)
		} else {
			break
		}
	}

}
