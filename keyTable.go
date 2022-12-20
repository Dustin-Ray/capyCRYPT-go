package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type KeyTable struct {
	treeview           *gtk.TreeView       // Displays list of keys currently imported into context
	store              *gtk.ListStore      //Contains a list of keys
	scrollableTreelist *gtk.ScrolledWindow // Allows scrolling for long list of keys
	grid               *gtk.Grid           // Grid container for TreeView
	keyList            map[string]KeyPair  // A list of all keys currently stored in this session
}

type KeyPair struct {
	Id      string `json:"Id"`
	Owner   string `json:"Owner"`
	KeyType string `json:"KeyType"`
	PubKey  string `json:"PubKey"`
	PrivKey string `json:"PrivKey"`
}

func (kt *KeyTable) JsonToKey(ctx *WindowCtx, filename string) error {
	data, _ := os.ReadFile(filename)
	var result = KeyPair{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	kt.importKey(ctx, result)
	return nil
}

func KeyToJSON(key *KeyPair) ([]byte, error) {
	fmt.Println(key)
	u, err := json.Marshal(KeyPair{
		Id:      key.Id,
		Owner:   key.Owner,
		KeyType: key.KeyType,
		PubKey:  key.PubKey,
		PrivKey: key.PrivKey})
	return u, err
}

func (kt *KeyTable) importKey(ctx *WindowCtx, key KeyPair) {
	// key := KeyPair{id: "12300", owner: "Jack Smith", keyType: "PRIVATE"}
	query := kt.keyList[key.Id]
	emptyKey := KeyPair{}
	if query == emptyKey {
		kt.keyList[key.Id] = key
		updateStore(kt, &key)
		ctx.updateStatus("Key import successful")
	} else {
		ctx.updateStatus("Key " + query.Id + " already imported")
	}
}

// Populate the model with imported key
func createAndFillModel() *gtk.ListStore {
	store, _ := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	return store
}

func updateStore(keyTable *KeyTable, keyData *KeyPair) {

	//create a row of data to append
	inValues := []interface{}{keyData.Id, keyData.Owner, keyData.KeyType}
	//get an iterator to the first row
	iter, _ := keyTable.store.GetIterFromString("0:0:0")
	inColumns := []int{0, 1, 2}
	keyTable.store.InsertWithValues(iter, 0, inColumns, inValues)
	keyTable.treeview.SetModel(keyTable.store.ToTreeModel())
}

// Sets up the key table
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
	newScrollableTreeList.SetSizeRequest(255, 450)

	newGrid.Attach(newScrollableTreeList, 0, 0, 8, 10)
	newScrollableTreeList.Add(newTreeView)

	ctx.keytable = &KeyTable{
		grid:               newGrid,
		treeview:           newTreeView,
		scrollableTreelist: newScrollableTreeList,
		keyList:            make(map[string]KeyPair),
	}

	ctx.fixed.Put(ctx.keytable.grid, 710, 80)
	ctx.keytable.store = createAndFillModel()
	newTreeView.SetModel(ctx.keytable.store.ToTreeModel())
	newTreeView.SetGridLines(gtk.TREE_VIEW_GRID_LINES_BOTH)
	newTreeView.SetActivateOnSingleClick(true)
	newTreeView.Connect("row-activated", func(tv *gtk.TreeView, path *gtk.TreePath) {
		// Get the list store
		liststore, _ := tv.GetModel()
		sel, _ := tv.GetSelection()
		_, iter, _ := sel.GetSelected()

		// Get the value from the list store
		id, _ := liststore.ToTreeModel().GetValue(iter, 0)
		name, _ := liststore.ToTreeModel().GetValue(iter, 1)
		keyType, _ := liststore.ToTreeModel().GetValue(iter, 2)

		idVal, _ := id.GetString()
		nameVal, _ := name.GetString()
		keyVal, _ := keyType.GetString()
		// Print the value to the console
		var test = ctx.keytable.keyList[idVal]
		ctx.status.SetText("Key data: " + idVal + nameVal + keyVal)
		ctx.loadedKey = &test
		ctx.updateStatus("key " + ctx.loadedKey.Id + " selected")
	})
}
