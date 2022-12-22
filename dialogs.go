package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/*
Displays a password entry dialog that asks a user for a passphrase.
Any input is considered a valid passphrase including the empty string.
OK button is disabled if content of password entry fields do not match.
*/
func showPasswordDialog(parent *gtk.Window, message string) []byte {
	// Create a dialog
	dialog, _ := gtk.DialogNew()

	dialog.SetTransientFor(parent)
	dialog.SetTitle("Enter " + message + " password:")

	okButton, _ := dialog.AddButton("OK", gtk.RESPONSE_OK)
	okButton.SetSensitive(true)
	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)

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

	// Create a password entry (this might be where everything is blowing up)
	entry, _ := gtk.EntryNew()
	entry.SetVisibility(false)
	hBox1.Add(entry)

	confirm, _ := gtk.EntryNew()
	confirm.SetVisibility(false)

	password1, _ := entry.GetText()
	password2, _ := confirm.GetText()

	confirm.Connect("changed", func() {
		// Get the entered password
		password1, _ = entry.GetText()
		password2, _ = confirm.GetText()
		if password2 == password1 {
			okButton.SetSensitive(true)
		} else {
			okButton.SetSensitive(false)
		}
	})

	entry.Connect("changed", func() {
		// Get the entered password

		password1, _ = entry.GetText()
		password2, _ = confirm.GetText()
		if password2 == password1 {
			okButton.SetSensitive(true)
		} else {
			okButton.SetSensitive(false)
		}
	})

	hBox2.Add(confirm)

	// Show the dialog
	dialog.ShowAll()
	dialog.Run()

	// Hide the dialog
	dialog.Hide()

	// Print the password
	// return []byte(password1)
	return []byte(password2)
}

func constructKey(parent *gtk.Window, key *KeyObj) {

	// Create a dialog
	dialog, _ := gtk.DialogNew()
	dialog.SetTransientFor(parent)
	dialog.SetTitle("Create new key: ")

	okButton, _ := dialog.AddButton("OK", gtk.RESPONSE_OK)
	okButton.SetSensitive(true)
	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)

	// Create a horizontal box
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

	hBox1.Add(ownerLbl)
	hBox2.Add(lbl)
	hBox3.Add(conf)

	owner, _ := gtk.EntryNew()
	pwd, _ := gtk.EntryNew()
	pwd.SetVisibility(false)

	confirm, _ := gtk.EntryNew()
	confirm.SetVisibility(false)

	hBox1.Add(owner)
	hBox2.Add(pwd)
	hBox3.Add(confirm)

	key.Id = hex.EncodeToString(generateRandomBytes(6))
	key.Owner = "NONE"
	key.KeyType = "PRIVATE"

	confirm.Connect("changed", func() {
		// Get the entered password
		password1, _ := pwd.GetText()
		password2, _ := confirm.GetText()
		ot, _ := owner.GetText()
		if password2 == password1 {
			setKeyData(key, password2, ot)
			okButton.SetSensitive(true)
		} else {
			okButton.SetSensitive(false)
		}
	})

	pwd.Connect("changed", func() {
		// Get the entered password
		password1, _ := pwd.GetText()
		password2, _ := confirm.GetText()
		ot, _ := owner.GetText()
		if password2 == password1 {
			setKeyData(key, password2, ot)
			okButton.SetSensitive(true)
		} else {
			okButton.SetSensitive(false)
		}
	})
	dialog.ShowAll()
	dialog.Run()
	dialog.Hide()
}

func setKeyData(key *KeyObj, password2 string, owner string) {
	s := new(big.Int).SetBytes(KMACXOF256([]byte(password2), []byte{}, 512, "K"))
	V := *GenPoint()
	V = *V.SecMul(s)
	key.Owner = owner
	key.PrivKey = password2
	key.PubKeyX = V.x.String()
	key.PubKeyY = V.y.String()
	key.DateCreated = time.Now().Format(time.RFC1123)
	key.Signature = "test"
}

func rightCLickMenu(ctx *WindowCtx) {
	// Create a Menu to display when right-clicking a row.
	menu, _ := gtk.MenuNew()
	// Create a MenuItem to be used in our Menu.
	menuItem, _ := gtk.MenuItemNewWithLabel("Details")
	// Add the MenuItem to the Menu.
	menu.Append(menuItem)

	// Connect a signal handler to the MenuItem's "activate" signal.
	menuItem.Connect("activate", func() { showKeyDetails(ctx) })
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
	PrivKey, _ := gtk.LabelNew("PRIV KEY: " + "NONE")
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

	dialog.ShowAll()
}

// A dialog that exports key data to a file
func saveDialog(ctx *WindowCtx, name string) {
	// Create a dialog that allows the user to save a file
	dialog, err := gtk.FileChooserDialogNewWith2Buttons("Save File", ctx.win,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		"Cancel", gtk.RESPONSE_CANCEL,
		"Save", gtk.RESPONSE_ACCEPT)
	if err != nil {
		panic(err)
	}

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
		ctx.updateStatus("Key data saved to: " + filename)
	}
	dialog.Destroy()
}

// A dialog that opens a key file. Handles any error in parsing file to key
func openDialog(ctx *WindowCtx) {

	// Create a new FileChooserDialog to open a file
	fileDialog, err := gtk.FileChooserDialogNewWith2Buttons("Open File", ctx.win,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"_Cancel", gtk.RESPONSE_CANCEL,
		"_Open", gtk.RESPONSE_ACCEPT)
	if err != nil {
		fmt.Println(err)
		return
	}
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

func getScrollableTextArea(ctx *WindowCtx, textForBuffer string) gtk.ScrolledWindow {

	scrollableTextArea, _ := gtk.ScrolledWindowNew(nil, nil)
	buf, _ := gtk.TextBufferNew(nil)
	textView, _ := gtk.TextViewNewWithBuffer(buf)
	textView.SetBuffer(buf)
	scrollableTextArea.Add(textView)
	scrollableTextArea.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	textView.SetWrapMode(gtk.WRAP_CHAR)
	scrollableTextArea.SetSizeRequest(200, 100)
	buf.SetText(textForBuffer)
	return *scrollableTextArea
}
