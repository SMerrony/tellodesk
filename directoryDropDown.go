package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/g3n/engine/gui"
)

// DirectoryDropDown is a GUI dropdown that permits navigation up and down
// the directory tree.  It is intended for use where there is a need for choosing
// a directory for some operation(s).
type DirectoryDropDown struct {
	*gui.DropDown
}

// NewDirectoryDropDown creates a new DirectoryDropDown element with the specified width and
// initial directory selected.  Above the highlighted directory will be a list of superior
// directories (if any), and below will appear any inferior directories.
func NewDirectoryDropDown(width float32, initialDir string) (ddd *DirectoryDropDown, err error) {
	// check initial directory is accessible
	cwd, err := filepath.Abs(initialDir)
	if err != nil {
		return nil, err
	}

	ddd = new(DirectoryDropDown)
	ddd.DropDown = gui.NewDropDown(width, gui.NewImageLabel(""))

	ddd.populateDirs(cwd)

	ddd.DropDown.Subscribe(gui.OnMouseUp, ddd.onSelect)

	return ddd, nil
}

func (ddd *DirectoryDropDown) populateDirs(cwd string) (err error) {

	// clear out any items
	for ddd.Len() > 0 {
		ddd.RemoveAt(0)
	}

	// superior directories
	var superiors []string
	d := cwd
	for d[len(d)-1] != filepath.Separator {
		d = filepath.Dir(d)
		superiors = append(superiors, d)
	}
	// add in reverse order
	for i := len(superiors) - 1; i >= 0; i-- {
		ddd.Add(gui.NewImageLabel(superiors[i]))
	}

	// initial directory
	cwdLab := gui.NewImageLabel(cwd)
	ddd.Add(cwdLab)
	ddd.SetSelected(cwdLab) // make it the selected item

	// inferior directories
	subfiles, err := ioutil.ReadDir(cwd)
	if err != nil {
		return err
	}
	for _, d := range subfiles {
		if d.IsDir() {
			ddd.Add(gui.NewImageLabel(fmt.Sprintf("%s%c%s", cwd, filepath.Separator, d.Name())))
		}
	}

	return nil
}

func (ddd *DirectoryDropDown) onSelect(e string, ev interface{}) {
	// Get selected image label and its txt
	label := ddd.Selected()
	text := label.Text()
	ddd.populateDirs(text)
}
