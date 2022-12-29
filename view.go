package main

/* View structure for MVC. Responds to messages from controller.*/

import (
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
	win          *gtk.Window      // Main window containing fixed container
	fixed        *gtk.Fixed       // Fixed allows for precise arbitrary placement of widgets
	loadedFile   *os.File         // Represents the selected file either dropped in window or selected from chooser
	notePad      *gtk.TextBuffer  // The text area displaying input and output
	initialState bool             // Signals if window is waiting for user input, can also be used to cancel running ops
	status       *gtk.Label       // Outputs operation status and error messages
	keytable     *KeyTable        // A table storing all imported keys
	loadedKey    *KeyObj          // The key to be used for any asymmetric encryptions
	fileMode     bool             // Determines whether to process a loaded file or notepad text
	progressBar  *gtk.ProgressBar // A bar to display status of ongoing operations
	buttons      *[]gtk.Button    // A list of pointers to all buttons added to the window
}

// Entry point
func main() {
	// rune521Tests()
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

// Sets the window to the initial state and initializes all widgets
func initialize() *WindowCtx {

	ctx := WindowCtx{}
	ctx.initialState = true
	ctx.fileMode = false
	ctx.loadedFile = nil
	ctx.win = setupWindow()
	ctx.fixed = newFixed()
	ctx.win.Add(ctx.fixed)

	ctx.notePad = setupNotepad(&ctx)

	setupButtons(&ctx)
	setupLabels(&ctx)
	setupMenuBar(&ctx)
	setupKeyTable(&ctx)
	setupProgressBar(&ctx)
	setupStatus(&ctx)
	setupCSS()
	return &ctx
}

func setupButtons(ctx *WindowCtx) {
	createButtons(ctx)

}

// Sets up CSS for styling of window
func setupCSS() {
	cssProvider, _ := gtk.CssProviderNew()
	cssProvider.LoadFromPath("style.css")
	screen, _ := gdk.ScreenGetDefault()
	gtk.AddProviderForScreen(screen, cssProvider, 0)
}

// Sets up status indicator for window
func setupStatus(ctx *WindowCtx) {
	ctx.status, _ = gtk.LabelNew("Status: Ready")
	ctx.fixed.Put(ctx.status, 245, 540)
}

// Updates the status message following an operation
func (w *WindowCtx) updateStatus(message string) { w.status.SetText("Status: " + message) }

// Overrides layouts to provide direct placement of widgets
func newFixed() *gtk.Fixed {
	fixed, _ := gtk.FixedNew()
	fixed.SetSizeRequest(1050, 590)
	return fixed
}

// Create a new toplevel window, set its title, and connect it to the
// "destroy" signal to exit the GTK main loop when it is destroyed.
func setupWindow() *gtk.Window {

	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("SOAPTool v 0.1")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetResizable(false)
	win.SetPosition(gtk.WIN_POS_CENTER)
	return win
}

/*
The notepad is the primary location of interaction for the application.
A user can either enter text directly or drag and drop a file into the window
to perform cryptographic operations on the data. If the user edits the file details, the
session switches to text mode and any operations requested are performed over the notepad
text. Any UTF-8 characters can be processed by the notepad, including emojis and other
special symbols. The application is not tested on any language except english.
*/
func setupNotepad(ctxWin *WindowCtx) *gtk.TextBuffer {
	scrollableTextArea, _ := gtk.ScrolledWindowNew(nil, nil)
	buf, _ := gtk.TextBufferNew(nil)
	textView, _ := gtk.TextViewNewWithBuffer(buf)
	textView.SetName("textArea")
	buf.SetText("Enter text or drag and drop file...")
	// Drag and drop a file into the notepad to get the path and switch to file mode
	textView.Connect("drag-data-received", func(ddarea *gtk.Box, ctx *gdk.DragContext, x int, y int, data *gtk.SelectionData, info uint, time uint32) {
		ctxWin.initialState = false
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
		buf.SetText("Successfully loaded ")
		ctxWin.updateStatus("Switched to file processing mode")
		ctxWin.fileMode = true
	})

	//Clear initial message from notepad
	textView.Connect("button-press-event", func() {
		if ctxWin.initialState {
			buf.SetText("")
			ctxWin.initialState = false
			textView.SetProperty("editable", true)
		}
	})

	//connect buffer, set size and other properties
	textView.SetBuffer(buf)
	scrollableTextArea.Add(textView)
	scrollableTextArea.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	textView.SetWrapMode(gtk.WRAP_CHAR)
	scrollableTextArea.SetSizeRequest(440, 450)
	ctxWin.fixed.Put(scrollableTextArea, 245, 80)
	return buf
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
	keysImport.Connect("activate", func() { importKeyDialog(ctx) })
	keysExport.Connect("activate", func() { exportPrivateKey(ctx) })

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
	keysLabel.SetTooltipMarkup("Click on a key to use it for message encryption. Right click the key to view details or to save it to a file.")

	buttonsLabel.SetName("textOpsLabel") //for CSS styling
	notePadLabel.SetName("notepadLabel") //for CSS styling
	buttonsLabel.SetName("statusLabel")  //for CSS styling
	buttonsLabel.SetName("keysLabel")    //for CSS styling

	ctx.fixed.Put(buttonsLabel, 40, 50)
	ctx.fixed.Put(notePadLabel, 245, 50)
	ctx.fixed.Put(keysLabel, 710, 50)
}

// Sets up a progress bar to display status of operations
func setupProgressBar(ctx *WindowCtx) {
	pb, err := gtk.ProgressBarNew()
	if err != nil {
		panic(err)
	}
	pb.SetSizeRequest(100, 40)
	ctx.progressBar = pb
}
