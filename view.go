package main

import (
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/*
A window context is an object that contains a Window and a Fixed.
Fixed holds all widgets so that they can be placed at precise pixels.
The style of the window is set by style.css.
*/
type WindowCtx struct {
	win          *gtk.Window
	fixed        *gtk.Fixed
	loadedFile   *os.File        // Represents any file currently pointed to
	notePad      *gtk.TextBuffer // The text area where input and output is processed
	initialState bool            // Signals if window has is waiting for user input
	status       string          // Outputs operation status and error messages
}

// Entry point
func main() {
	gtk.Init(nil)

	settings, _ := gtk.SettingsGetDefault()
	settings.SetProperty("gtk-application-prefer-dark-theme", true)
	window := initialize()

	window.win.ShowAll()
	gtk.Main()
}

// Overrides layouts to provide direct placement of widgets
func newFixed() *gtk.Fixed {
	fixed, err := gtk.FixedNew()
	fixed.SetSizeRequest(1000, 590)
	if err != nil {
		panic(err)
	}
	return fixed
}

// Create a new toplevel window, set its title, and connect it to the
// "destroy" signal to exit the GTK main loop when it is destroyed.
func setupWindow() *gtk.Window {

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		panic(err)
	}
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
func createScrollableTextArea(ctx *WindowCtx) *gtk.TextBuffer {
	scrollableTextArea, _ := gtk.ScrolledWindowNew(nil, nil)
	buf, _ := gtk.TextBufferNew(nil)
	textView, _ := gtk.TextViewNewWithBuffer(buf)
	textView.SetName("textArea")
	buf.SetText("Enter text or drag and drop file..")
	// Connect the drag and drop area signals
	textView.Connect("drag-data-received", func(ddarea *gtk.Box, ctx *gdk.DragContext, x int, y int, data *gtk.SelectionData, info uint, time uint32) {
		textView.SetBuffer(buf)
		buf.SetText("File Dropped: " + string(data.GetData()))
		textView.SetProperty("editable", false)
		textView.SetProperty("cursor-visible", false)
	})

	textView.Connect("button-press-event", func() {
		if ctx.initialState {
			buf.SetText("")
			ctx.initialState = false
		}
	})

	scrollableTextArea.Add(textView)
	// Set the scrolling policy
	scrollableTextArea.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	textView.SetWrapMode(gtk.WRAP_CHAR)
	// Set the size of the scrollable text area
	scrollableTextArea.SetSizeRequest(440, 450)
	ctx.fixed.Put(scrollableTextArea, 290, 80)
	return buf
}

func createTable() {}

// Sets the window to the initial state.
func initialize() *WindowCtx {

	ctx := WindowCtx{}
	ctx.win = setupWindow()
	ctx.fixed = newFixed()

	ctx.initialState = true
	ctx.win.Add(ctx.fixed)

	ctx.loadedFile = nil
	ctx.notePad = createScrollableTextArea(&ctx)
	ctx.status = "Status: Ready"

	setupButtons(&ctx)
	setupLabels(&ctx)
	setupMenuBar(ctx.fixed)

	cssProvider, _ := gtk.CssProviderNew()
	cssProvider.LoadFromPath("style.css")
	screen, _ := gdk.ScreenGetDefault()
	gtk.AddProviderForScreen(screen, cssProvider, 0)
	return &ctx
}

// Resets context to initial state
func (ctx *WindowCtx) Reset() {
	ctx.notePad = nil
	ctx.initialState = true
	ctx.loadedFile = nil
	ctx.notePad = createScrollableTextArea(ctx)
	ctx.status = "Status: Ready"
	ctx.win.ShowAll()
}

// Adds file menu bar to window.
func setupMenuBar(ctx *gtk.Fixed) {

	menubar, _ := gtk.MenuBarNew()
	fileMi, _ := gtk.MenuItemNewWithLabel("File")
	edit, _ := gtk.MenuItemNewWithLabel("Edit")

	menubar.Append(fileMi)
	menubar.Append(edit)
	ctx.Add(menubar)

}

// Adds labels to window
func setupLabels(ctx *WindowCtx) {
	// Create a label and add it to the fixed container
	buttonsLabel, _ := gtk.LabelNew("Text Operations:")
	notePadLabel, _ := gtk.LabelNew("Notepad:")
	statusLabel, _ := gtk.LabelNew("Status: Ready")

	buttonsLabel.SetName("textOpsLabel") //for CSS styling
	notePadLabel.SetName("notepadLabel") //for CSS styling
	buttonsLabel.SetName("statusLabel")  //for CSS styling

	ctx.fixed.Put(buttonsLabel, 40, 40)
	ctx.fixed.Put(notePadLabel, 290, 40)
	ctx.fixed.Put(statusLabel, 290, 545)

}

// adds buttons in a factory style to fixed context
func setupButtons(ctx *WindowCtx) {
	labelList := []string{"Compute Hash", "Compute Tag", "Encrypt With Password", "Decrypt With Password",
		"Generate Keypair", "Encrypt With Key", "Decrypt With Key", "Sign With Key", "Verify Signature"}

	for i, label := range labelList {
		btn, _ := gtk.ButtonNewWithLabel(label)
		ctx.fixed.Put(btn, 40, 80+i*45)
	}
	reset, _ := gtk.ButtonNewWithLabel("Reset")
	reset.SetName("resetButton") //for CSS styling
	reset.Connect("button-press-event", func() {
		ctx.Reset()
	})
	ctx.fixed.Put(reset, 40, 510)
}
