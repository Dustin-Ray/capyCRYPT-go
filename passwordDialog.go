package main

import (
	"math/big"
	"time"

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

	// Create a password entry
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

	// Create a label
	ownerLbl, _ := gtk.LabelNew("Owner:       ")
	lbl, _ := gtk.LabelNew("Password: ")
	conf, _ := gtk.LabelNew("Confirm:     ")

	hBox1.Add(ownerLbl)
	hBox2.Add(lbl)
	hBox3.Add(conf)

	owner, _ := gtk.EntryNew()
	// Create a password pwd
	pwd, _ := gtk.EntryNew()
	pwd.SetVisibility(false)

	confirm, _ := gtk.EntryNew()
	confirm.SetVisibility(false)

	ownerField, _ := owner.GetText()
	password1, _ := pwd.GetText()
	password2, _ := confirm.GetText()

	hBox1.Add(owner)
	hBox2.Add(pwd)
	hBox3.Add(confirm)
	confirm.Connect("changed", func() {
		// Get the entered password
		password1, _ = pwd.GetText()
		password2, _ = confirm.GetText()
		if password2 == password1 {
			s := new(big.Int).SetBytes(KMACXOF256([]byte(password2), []byte{}, 512, "K"))
			V := *GenPoint()
			V = *V.SecMul(s)
			key.Id = BytesToHexString(generateRandomBytes(6))
			key.Owner, _ = owner.GetText()
			key.PrivKey = password2
			key.KeyType = "PRIVATE"
			key.PubKeyX = V.x.String()
			key.PubKeyY = V.y.String()
			key.DateCreated = time.Now().Format(time.RFC1123)
			key.Signature = "test"
			okButton.SetSensitive(true)
		} else {
			okButton.SetSensitive(false)
		}
	})

	pwd.Connect("changed", func() {
		// Get the entered password
		password1, _ = pwd.GetText()
		password2, _ = confirm.GetText()
		if password2 == password1 {
			key.Id = "0000"
			key.Owner = ownerField
			key.KeyType = "PRIVATE"
			key.PrivKey = password2

			key.DateCreated = time.Now().Format(time.RFC1123)

			okButton.SetSensitive(true)
		} else {
			okButton.SetSensitive(false)
		}
	})

	// Show the dialog
	dialog.ShowAll()
	dialog.Run()

	// Hide the dialog
	dialog.Hide()

}
