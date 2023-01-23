package main

import (
	"encoding/json"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type KeyTable struct {
	treeview           *gtk.TreeView       // Displays list of keys currently imported into context
	store              *gtk.ListStore      // Contains a list of keys
	scrollableTreelist *gtk.ScrolledWindow // Allows scrolling for long list of keys
	grid               *gtk.Grid           // Grid container for TreeView
	keyList            map[string]KeyObj   // A list of all keys currently stored in this session
}

type KeyObj struct {
	Id      string `json:"Id"`      //Represents the unique ID of the key
	Owner   string `json:"Owner"`   //Represents the owner of the key, can be arbitrary
	KeyType string `json:"KeyType"` /*Acceptable values are PUBLIC or PRIVATE.
	PUBLIC keys are used only for encryptions, while PRIVATE keys can
	encrypt or decrypt.
	*/
	PubKeyX     string `json:"PubKeyX"`     //big.Int value representing E521 X coordinate
	PubKeyY     string `json:"PubKeyY"`     //big.Int value representing E521 X coordinate
	PrivKey     string `json:"PrivKey"`     //big.Int value representing secret scalar, nil if KeyType is PUBLIC
	DateCreated string `json:"DateCreated"` //Date key was generated
	Signature   string `json:"Signature"`   //Nil unless PUBLIC. Signs 128 bit SHA3 hash of this KeyObj
}

// Converts JSON to KeyObj. Returns error if conversion is unsuccessful.
func (kt *KeyTable) JsonToKey(ctx *WindowCtx, filename string) error {
	data, _ := os.ReadFile(filename)
	var result = KeyObj{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	kt.importKey(ctx, result)
	return nil
}

// Converts a key to JSON format
func KeyToJSON(key *KeyObj) ([]byte, error) {
	u, err := json.Marshal(KeyObj{
		Id:          key.Id,
		Owner:       key.Owner,
		KeyType:     key.KeyType,
		PubKeyX:     key.PubKeyX,
		PubKeyY:     key.PubKeyY,
		PrivKey:     key.PrivKey,
		DateCreated: key.DateCreated,
		Signature:   key.Signature})
	return u, err
}

// Attempts to parse a JSON file into a KeyObj. Declines to import duplicate keys.
func (kt *KeyTable) importKey(ctx *WindowCtx, key KeyObj) {
	query := kt.keyList[key.Id]
	emptyKey := KeyObj{}
	if query == emptyKey {
		kt.keyList[key.Id] = key
		updateStore(kt, &key)
		ctx.updateStatus("Key " + key.Id + " imported")
	} else {
		ctx.updateStatus("Key " + query.Id + " already imported")
	}
}

// updates the list store with a new key
func updateStore(keyTable *KeyTable, keyData *KeyObj) {

	//create a row of data to append
	inValues := []interface{}{keyData.Id, keyData.Owner, keyData.KeyType}
	//get an iterator to the first row
	iter, _ := keyTable.store.GetIterFromString("0:0:0")
	inColumns := []int{0, 1, 2}
	keyTable.store.InsertWithValues(iter, 0, inColumns, inValues)
	keyTable.treeview.SetModel(keyTable.store.ToTreeModel())
}

// Populate the model with imported key
func createAndFillModel() *gtk.ListStore {
	store, _ := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	return store
}

// Constructs the key table
func setupKeyTable(ctx *WindowCtx) {

	newGrid, _ := gtk.GridNew()
	newGrid.SetColumnHomogeneous(true)
	newGrid.SetRowHomogeneous(true)

	newTreeView, _ := gtk.TreeViewNew()
	for i, columnTitle := range []string{"Key ID:     ", "Key Owner:    ", "Type:    "} {
		renderer, _ := gtk.CellRendererTextNew()
		column, _ := gtk.TreeViewColumnNewWithAttribute(columnTitle, renderer, "text", i)
		newTreeView.AppendColumn(column)
	}

	newScrollableTreeList, _ := gtk.ScrolledWindowNew(nil, nil)
	newScrollableTreeList.SetVExpand(true)
	newScrollableTreeList.SetSizeRequest(295, 205)

	newGrid.Attach(newScrollableTreeList, 0, 0, 8, 10)
	newScrollableTreeList.Add(newTreeView)

	ctx.keytable = &KeyTable{
		grid:               newGrid,
		treeview:           newTreeView,
		scrollableTreelist: newScrollableTreeList,
		keyList:            make(map[string]KeyObj),
	}

	ctx.fixed.Put(ctx.keytable.grid, 710, 80)
	ctx.keytable.store = createAndFillModel()
	newTreeView.SetModel(ctx.keytable.store.ToTreeModel())
	newTreeView.SetGridLines(gtk.TREE_VIEW_GRID_LINES_BOTH)
	newTreeView.SetActivateOnSingleClick(true)
	newTreeView.SetHoverSelection(true)

	newTreeView.Connect("row-activated", func(tv *gtk.TreeView, path *gtk.TreePath) {
		// Get the list store
		liststore, _ := tv.GetModel()
		sel, _ := tv.GetSelection()
		_, iter, _ := sel.GetSelected()

		// Get the value from the list store
		id, _ := liststore.ToTreeModel().GetValue(iter, 0)
		idVal, _ := id.GetString()
		var lookupKey = ctx.keytable.keyList[idVal]
		ctx.loadedKey = &lookupKey
		ctx.updateStatus("key " + ctx.loadedKey.Id + " selected")

	})
	newTreeView.Connect("button-press-event", func(tv *gtk.TreeView, event *gdk.Event) {

		// Get the list store
		liststore, err := tv.GetModel()
		if err != nil {
			return
		}
		sel, err := tv.GetSelection()
		_, iter, ok := sel.GetSelected()
		if !ok || err != nil {
			return
		}

		// Get the value from the list store
		id, err := liststore.ToTreeModel().GetValue(iter, 0)
		if err != nil {
			return
		}
		idVal, err := id.GetString()
		if err != nil {
			return
		}
		var lookupKey = ctx.keytable.keyList[idVal]
		ctx.loadedKey = &lookupKey
		ctx.updateStatus("key " + ctx.loadedKey.Id + " selected")

	})
	rightCLickMenu(ctx)
}
