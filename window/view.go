package CryptoTool

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		panic(err)
	}
	win.SetTitle("Golang GUI with GTK+")

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(800, 600)
	win.SetPosition(gtk.WIN_POS_CENTER)

	// Create a fixed container
	fixed, err := gtk.FixedNew()
	fixed.SetSizeRequest(800, 600)
	if err != nil {
		panic(err)
	}
	win.Add(fixed)

	// Create a button and add it to the fixed container
	btn, err := gtk.ButtonNewWithLabel("Click Me!")
	if err != nil {
		panic(err)
	}
	btn.SetSizeRequest(80, 35)
	fixed.Put(btn, 200, 100)

	btn.Connect("clicked", func() {
		btn.SetLabel("Test")
	})

	// Create a label and add it to the fixed container
	lbl, err := gtk.LabelNew("Hello World!")
	if err != nil {
		panic(err)
	}
	fixed.Put(lbl, 50, 80)

	menubar, _ := gtk.MenuBarNew()
	fileMi, _ := gtk.MenuItemNewWithLabel("File")
	edit, _ := gtk.MenuItemNewWithLabel("Edit")

	fixed.Add(menubar)
	menubar.Append(fileMi)
	menubar.Append(edit)

	fixed.Put(createScrollableTextArea(), 300, 100)

	// Set the default window size.
	win.SetDefaultSize(300, 200)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func createScrollableTextArea() *gtk.ScrolledWindow {
	scrollableTextArea, _ := gtk.ScrolledWindowNew(nil, nil)
	buf, _ := gtk.TextBufferNew(nil)
	textView, _ := gtk.TextViewNewWithBuffer(buf)

	// Connect the drag and drop area signals
	textView.Connect("drag-data-received", func(ddarea *gtk.Box, ctx *gdk.DragContext, x int, y int, data *gtk.SelectionData, info uint, time uint32) {

		textView.SetBuffer(buf)
		buf.SetText("File Dropped: " + string(data.GetData()))
	})

	scrollableTextArea.Add(textView)
	// Set the scrolling policy
	scrollableTextArea.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	textView.SetWrapMode(gtk.WRAP_CHAR)
	// Set the size of the scrollable text area
	scrollableTextArea.SetSizeRequest(400, 400)
	return scrollableTextArea
}
