package main

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type TreeViewFilterWindow struct {
	grid                  *gtk.Grid
	softwareListstore     *gtk.ListStore
	currentFilterLanguage string
	languageFilter        *gtk.TreeModelFilter
	treeview              *gtk.TreeView
	buttons               []*gtk.Button
	scrollableTreelist    *gtk.ScrolledWindow
}

func (tvf *TreeViewFilterWindow) languageFilterFunc(model *gtk.TreeModel, iter *gtk.TreeIter, data interface{}) bool {
	if tvf.currentFilterLanguage == "" || tvf.currentFilterLanguage == "None" {
		return true
	} else {
		return model.GetValue(iter, 2).GetString() == tvf.currentFilterLanguage
	}
}

func (tvf *TreeViewFilterWindow) onSelectionButtonClicked(widget *gtk.Button) {
	tvf.currentFilterLanguage = widget.GetLabelText()
	fmt.Printf("%s language selected!\n", tvf.currentFilterLanguage)
	tvf.languageFilter.Refilter()
}

func NewTreeViewFilterWindow() (*TreeViewFilterWindow, error) {
	tvf := &TreeViewFilterWindow{}
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	win.SetTitle("Treeview Filter Demo")
	win.SetBorderWidth(10)

	grid, err := gtk.GridNew()
	if err != nil {
		return nil, err
	}
	grid.SetColumnHomogeneous(true)
	grid.SetRowHomogeneous(true)
	win.Add(grid)
	tvf.grid = grid

	softwareListstore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_INT, glib.TYPE_STRING)
	if err != nil {
		return nil, err
	}

	softwareList := []string{}
	for _, softwareRef := range softwareList {
		softwareListstore.Append(softwareRef)
	}
	tvf.softwareListstore = softwareListstore
	tvf.currentFilterLanguage = ""

	languageFilter, err := tvf.softwareListstore.FilterNew()
	if err != nil {
		return nil, err
	}
	languageFilter.SetVisibleFunc(tvf.languageFilterFunc)
	tvf.languageFilter = languageFilter

	treeview, err := gtk.TreeViewNewWithModel(languageFilter)
	if err != nil {
		return nil, err
	}
	for i, columnTitle := range []string{"Software", "Release Year", "Programming Language"} {
		renderer, err := gtk.CellRendererTextNew()
		if err != nil {
			return nil, err
		}
		column, err := gtk.TreeViewColumnNewWithAttribute(columnTitle, renderer, "text", i)
		if err != nil {
			return nil, err
		}
		treeview.AppendColumn(column)
	}

	tvf.buttons = []*gtk.Button{}
	for _, progLanguage := range []string{"Java", "C", "C++", "Python", "None"} {
		button, err := gtk.ButtonNewWithLabel(progLanguage)
		if err != nil {
			return nil, err
		}
		tvf.buttons = append(tvf.buttons, button)
		button.Connect("clicked", tvf.onSelectionButtonClicked)
	}

	scrollableTreelist, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	scrollableTreelist.SetVExpand(true)
	grid.Attach(scrollableTreelist, 0, 0, 8, 10)
	grid.AttachNextTo(tvf.buttons[0], scrollableTreelist, gtk.POS_BOTTOM, 1, 1)
	for i, button := range tvf.buttons[1:] {
		grid.AttachNextTo(button, tvf.buttons[i], gtk.POS_RIGHT, 1, 1)
	}
	scrollableTreelist.Add(treeview)
	tvf.scrollableTreelist = scrollableTreelist
	tvf.treeview = treeview
	win.ShowAll()
	return tvf, nil
}
