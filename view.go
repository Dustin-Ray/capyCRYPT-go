package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type KeyTable struct {
	treeview           *gtk.TreeView       // Displays list of keys currently imported
	scrollableTreelist *gtk.ScrolledWindow // Allows scrolling for long list of keys
	grid               *gtk.Grid           // Grid container for TreeView
}

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
	keytable     *KeyTable       // A table storing all imported keys
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
	ctx.fixed.Put(scrollableTextArea, 245, 80)
	return buf
}

func onEmissionReceived(tv *gtk.TreeView, path *gtk.TreePath, column *gtk.TreeViewColumn) {

	tvmodel, _ := tv.GetModel()
	sel, _ := tv.GetSelection()
	_, iter, _ := sel.GetSelected()
	nenf := tvmodel.ToTreeModel().IterNChildren(iter)
	var child gtk.TreeIter
	for ind := 0; ind < nenf; ind++ {
		tvmodel.ToTreeModel().IterParent(&child, iter)
	}
	fmt.Println("clicked!")
}

// Connect the signal handler to the object's signal *gtk.Object.Connect("signal-name", onEmissionReceived, nil)

func createAndFillModel() *gtk.TreeModel {

	inColumns := []int{0, 1, 2}
	inValues0 := []interface{}{000, "test0", "PRIVATE"}
	inValues1 := []interface{}{555, "test1", "PUBLIC"}
	inValues2 := []interface{}{111, "test3", "PRIVATE"}
	inValues3 := []interface{}{888, "test4", "PRIVATE"}
	inValues4 := []interface{}{-100, "test5", "PUBLIC"}

	store, _ := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	iter := store.Append()
	store.InsertWithValues(iter, 0, inColumns, inValues0)
	store.InsertWithValues(iter, 1, inColumns, inValues1)
	store.InsertWithValues(iter, 2, inColumns, inValues2)
	store.InsertWithValues(iter, 3, inColumns, inValues3)
	store.InsertWithValues(iter, 4, inColumns, inValues4)
	return store.ToTreeModel()
}

// Sets up the key table
func setupKeyTable(ctx *WindowCtx) {

	newGrid, _ := gtk.GridNew()
	newGrid.SetColumnHomogeneous(true)
	newGrid.SetRowHomogeneous(true)

	newTreeView, _ := gtk.TreeViewNew()
	for i, columnTitle := range []string{"Key ID:     ", "Key Name:    ", "Type:    "} {
		renderer, _ := gtk.CellRendererTextNew()
		column, _ := gtk.TreeViewColumnNewWithAttribute(columnTitle, renderer, "text", i)
		newTreeView.AppendColumn(column)
	}

	newScrollableTreeList, _ := gtk.ScrolledWindowNew(nil, nil)
	newScrollableTreeList.SetVExpand(true)
	newScrollableTreeList.SetSizeRequest(245, 450)

	newGrid.Attach(newScrollableTreeList, 0, 0, 8, 10)
	newScrollableTreeList.Add(newTreeView)

	ctx.keytable = &KeyTable{
		grid:               newGrid,
		treeview:           newTreeView,
		scrollableTreelist: newScrollableTreeList,
	}

	ctx.fixed.Put(ctx.keytable.grid, 700, 80)
	newTreeView.SetModel(createAndFillModel())
	newTreeView.SetGridLines(gtk.TREE_VIEW_GRID_LINES_BOTH)
	newTreeView.Connect("row-activated", func(tv *gtk.TreeView, path *gtk.TreePath) string {
		// Get the list store
		liststore, _ := tv.GetModel()
		sel, _ := tv.GetSelection()
		_, iter, _ := sel.GetSelected()

		// Get the value from the list store
		id, _ := liststore.ToTreeModel().GetValue(iter, 0)
		name, _ := liststore.ToTreeModel().GetValue(iter, 1)
		keyType, _ := liststore.ToTreeModel().GetValue(iter, 2)

		idVal, _ := id.GoValue()
		nameVal, _ := name.GetString()
		keyVal, _ := keyType.GetString()
		// Print the value to the console

		ctx.status = "Key data: " + idVal.(string) + nameVal + keyVal
		fmt.Println("Key data: ", idVal, nameVal, keyVal)
		return ""
	})
}

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
	setupKeyTable(&ctx)

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
	setupKeyTable(ctx)
	ctx.status = "Status: Ready"
	ctx.win.ShowAll()
}

// Adds file menu bar to window.
func setupMenuBar(ctx *gtk.Fixed) {

	menubar, _ := gtk.MenuBarNew()
	fileMi, _ := gtk.MenuItemNewWithLabel("File")
	edit, _ := gtk.MenuItemNewWithLabel("Keys")

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
	ctx.fixed.Put(notePadLabel, 245, 40)
	ctx.fixed.Put(statusLabel, 245, 545)

}

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

	buttonList[0].Connect("clicked", func() {
		text, _ := ctx.notePad.GetText(ctx.notePad.GetStartIter(), ctx.notePad.GetEndIter(), true)
		regex := regexp.MustCompile(`(?m)^\s*\r?\n`)
		text = regex.ReplaceAllString(text, "")
		ctx.notePad.SetText(ComputeSHA3HASH(text))
	}) //etc....

	reset, _ := gtk.ButtonNewWithLabel("Reset")
	reset.SetName("resetButton") //for CSS styling
	reset.Connect("button-press-event", func() {
		ctx.Reset()
	})
	ctx.fixed.Put(reset, 40, 510)
}
