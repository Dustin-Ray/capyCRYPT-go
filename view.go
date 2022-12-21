package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/*
A window context is an object that contains a Window and a Fixed.
Fixed holds all widgets so that they can be placed at precise pixels.
The style of the window is set by style.css.
*/
type WindowCtx struct {
	win          *gtk.Window     // Main window containing fixed container
	fixed        *gtk.Fixed      // Fixed allows for precise arbitrary placement of widgets
	loadedFile   *os.File        // Represents any file currently pointed to
	notePad      *gtk.TextBuffer // The text area where input and output is processed
	initialState bool            // Signals if window has is waiting for user input
	status       *gtk.Label      // Outputs operation status and error messages
	keytable     *KeyTable       // A table storing all imported keys
	loadedKey    *KeyObj         // The key to be used for any asymetric encryptions
	fileMode     bool            //Determines whether to process a loaded file or notepad text
}

// Entry point
func main() {
	gtk.Init(nil)
	window := initialize()
	settings, _ := gtk.SettingsGetDefault()
	err := settings.SetProperty("gtk-application-prefer-dark-theme", true) //try to default to dark theme
	if err != nil {
		window.notePad.SetText("Failed to load dark theme")
	}
	window.win.ShowAll()
	gtk.Main()
}

// Sets the window to the initial state.
func initialize() *WindowCtx {

	ctx := WindowCtx{}
	ctx.win = setupWindow()
	ctx.fixed = newFixed()

	ctx.initialState = true
	ctx.fileMode = false
	ctx.win.Add(ctx.fixed)

	ctx.loadedFile = nil
	ctx.notePad = createScrollableTextArea(&ctx)

	ctx.status, _ = gtk.LabelNew("")
	ctx.status.SetText("Status: Ready")
	ctx.fixed.Put(ctx.status, 245, 535)

	setupButtons(&ctx)
	setupLabels(&ctx)
	setupMenuBar(&ctx)
	setupKeyTable(&ctx)

	cssProvider, _ := gtk.CssProviderNew()
	cssProvider.LoadFromPath("style.css")
	screen, _ := gdk.ScreenGetDefault()
	gtk.AddProviderForScreen(screen, cssProvider, 0)
	return &ctx
}

// Updates the status message following an operation
func (w *WindowCtx) updateStatus(message string) { w.status.SetText("Status: " + message) }

// Overrides layouts to provide direct placement of widgets
func newFixed() *gtk.Fixed {
	fixed, _ := gtk.FixedNew()
	fixed.SetSizeRequest(1000, 590)
	return fixed
}

// Create a new toplevel window, set its title, and connect it to the
// "destroy" signal to exit the GTK main loop when it is destroyed.
func setupWindow() *gtk.Window {

	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("CryptoTool v 0.1")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

/*
Creates a scrollable text area to display input and output of operations. Contains
drag and drop functionality for file operations.
*/
func createScrollableTextArea(ctxWin *WindowCtx) *gtk.TextBuffer {
	scrollableTextArea, _ := gtk.ScrolledWindowNew(nil, nil)
	buf, _ := gtk.TextBufferNew(nil)
	textView, _ := gtk.TextViewNewWithBuffer(buf)
	textView.SetName("textArea")
	buf.SetText("Enter text or drag and drop file...")
	// Connect the drag and drop area signals
	textView.Connect("drag-data-received", func(ddarea *gtk.Box, ctx *gdk.DragContext, x int, y int, data *gtk.SelectionData, info uint, time uint32) {
		ctxWin.initialState = false
		ctxWin.fileMode = true
		buf.SetText("")
		var replacer = strings.NewReplacer("\r\n", "")
		filePath := replacer.Replace(string(data.GetData()))
		regex := regexp.MustCompile(`^.{7}`)
		filePath = regex.ReplaceAllString(filePath, "")
		// open the dropped file
		file, err := os.Open(filePath)
		if err != nil {
			buf.SetText(err.Error())
		} else {
			ctxWin.loadedFile = file
		}
		textView.SetCanFocus(false)
		textView.SetProperty("cursor-visible", false)
		buf.SetText("Successfully loaded: ")
		ctxWin.updateStatus("Switched to file processing mode")
	})

	textView.Connect("button-press-event", func() {
		if ctxWin.initialState {
			buf.SetText("")
			ctxWin.initialState = false
			textView.SetProperty("editable", true)
		}
	})
	textView.SetBuffer(buf)
	scrollableTextArea.Add(textView)
	scrollableTextArea.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	textView.SetWrapMode(gtk.WRAP_CHAR)
	scrollableTextArea.SetSizeRequest(440, 450)
	ctxWin.fixed.Put(scrollableTextArea, 245, 80)
	return buf
}

// Resets context to initial state
func (ctx *WindowCtx) Reset() {
	ctx.notePad = nil
	ctx.initialState = true
	ctx.loadedFile = nil
	ctx.notePad = createScrollableTextArea(ctx)
	setupKeyTable(ctx)
	ctx.status.SetText("Status: Ready")
	ctx.win.ShowAll()
}

// Adds file menu bar to window.
func setupMenuBar(ctx *WindowCtx) {

	menubar, _ := gtk.MenuBarNew()
	fileMenu, _ := gtk.MenuItemNewWithLabel("File")
	keysMenu, _ := gtk.MenuItemNewWithLabel("Keys")
	keysDropDown, _ := gtk.MenuNew()
	fileDropDown, _ := gtk.MenuNew()

	keysMenu.SetSubmenu(keysDropDown)
	fileMenu.SetSubmenu(fileDropDown)

	keysImport, _ := gtk.MenuItemNewWithLabel("Import")
	keysExport, _ := gtk.MenuItemNewWithLabel("Export")
	fileLoad, _ := gtk.MenuItemNewWithLabel("Load File")
	fileSave, _ := gtk.MenuItemNewWithLabel("Save File")
	help, _ := gtk.MenuItemNewWithLabel("How To Use")
	exit, _ := gtk.MenuItemNewWithLabel("Exit")

	//setup import and export funtionality
	keysImport.Connect("activate", func() { openDialog(ctx) })
	keysExport.Connect("activate", func() {
		if ctx.loadedKey != nil {
			KeyToJSON(ctx.loadedKey)
			saveDialog(ctx, "Save File")
		} else {
			ctx.updateStatus("Export failed - no key selected")
		}

	})

	keysDropDown.Append(keysImport)
	keysDropDown.Append(keysExport)

	fileDropDown.Append(fileLoad)
	fileDropDown.Append(fileSave)
	fileDropDown.Append(help)
	fileDropDown.Append(exit)

	menubar.Append(fileMenu)
	menubar.Append(keysMenu)
	ctx.fixed.Add(menubar)
}

// Adds labels to window
func setupLabels(ctx *WindowCtx) {
	buttonsLabel, _ := gtk.LabelNew("Text Operations:")
	notePadLabel, _ := gtk.LabelNew("Notepad:")
	keysLabel, _ := gtk.LabelNew("Select an encryption key: ")

	buttonsLabel.SetName("textOpsLabel") //for CSS styling
	notePadLabel.SetName("notepadLabel") //for CSS styling
	buttonsLabel.SetName("statusLabel")  //for CSS styling
	buttonsLabel.SetName("keysLabel")    //for CSS styling

	ctx.fixed.Put(buttonsLabel, 40, 50)
	ctx.fixed.Put(notePadLabel, 245, 50)
	ctx.fixed.Put(keysLabel, 710, 50)

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
