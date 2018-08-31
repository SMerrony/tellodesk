package main

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/gui/assets/icon"
	"github.com/g3n/engine/math32"
)

const fileSelectWidth, fileSelectHeight = 400, 300

// FileSelect is a general-purpose file selector
type FileSelect struct {
	*gui.Window
	parent *gui.Panel
	path   *gui.Label
	name   *gui.Edit
	list   *gui.List
	filter string
}

// NewFileSelect displays and returns a general-purpose file selector.
// The suffix should either be empty, or of the exact form "*.<suff>" and if present will cause
// non-directory files to be filtered according the suffix.
func NewFileSelect(parent *gui.Panel, initPath string, title string, suffix string) (fs *FileSelect, err error) {

	fs = new(FileSelect)
	fs.parent = parent
	fs.Window = gui.NewWindow(fileSelectWidth, fileSelectHeight)
	fs.SetPaddings(4, 4, 4, 4)
	fs.SetTitle(title)
	fs.SetColor(math32.NewColor("Gray"))
	fs.SetCloseButton(false)
	if suffix != "" && suffix[0] != '*' {
		return nil, errors.New("Suffix must start with '*' or be omitted")
	}
	fs.filter = suffix

	// Set vertical box layout for the whole panel
	vbl := gui.NewVBoxLayout()
	vbl.SetSpacing(4)
	fs.SetLayout(vbl)

	// Creates path label
	fs.path = gui.NewLabel("path")
	fs.Add(fs.path)

	// Creates list
	fs.list = gui.NewVList(0, 0)
	fs.list.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 5, AlignH: gui.AlignWidth})
	fs.list.Subscribe(gui.OnChange, func(evname string, ev interface{}) {
		fs.onSelect()
	})
	fs.Add(fs.list)

	// create name edit
	fs.name = gui.NewEdit(fileSelectWidth, suffix)
	fs.name.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 0, AlignH: gui.AlignWidth})
	// fs.name.SetBgColor(math32.NewColor("Gray"))
	// fs.name.SetColor(math32.NewColor("White"))
	fs.Add(fs.name)

	// Button container panel
	bc := gui.NewPanel(0, 0)
	bcl := gui.NewHBoxLayout()
	bcl.SetAlignH(gui.AlignWidth)
	bc.SetLayout(bcl)
	bc.SetLayoutParams(&gui.VBoxLayoutParams{Expand: 1, AlignH: gui.AlignWidth})
	fs.Add(bc)

	// Creates OK button
	bok := gui.NewButton("OK")
	bok.SetLayoutParams(&gui.HBoxLayoutParams{Expand: 0, AlignV: gui.AlignCenter})
	bok.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		fs.Dispatch("OnOK", nil)
	})
	bc.Add(bok)

	// Creates Cancel button
	bcan := gui.NewButton("Cancel")
	bcan.SetLayoutParams(&gui.HBoxLayoutParams{Expand: 0, AlignV: gui.AlignCenter})
	bcan.Subscribe(gui.OnClick, func(evname string, ev interface{}) {
		fs.Dispatch("OnCancel", nil)
	})
	bc.Add(bcan)

	fs.setPath(initPath)

	parent.Add(fs)
	fs.SetPosition(parent.Width()/2-fileSelectWidth/2, parent.Height()/2-fileSelectHeight/2)

	fs.Root().SetModal(fs)

	return fs, nil
}

// Close removes a file selector from the GUI
func (fs *FileSelect) Close() {
	fs.Root().SetModal(nil)
	fs.parent.Remove(fs)
}

func (fs *FileSelect) setPath(path string) error {

	// Open path file or dir
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Checks if it is a directory
	files, err := f.Readdir(0)
	if err != nil {
		return err
	}
	fs.path.SetText(path)

	// Sort files by name
	sort.Sort(listFileInfo(files))

	// Reads directory contents and loads into the list
	fs.list.Clear()
	// Adds previous directory
	prev := gui.NewImageLabel("..")
	prev.SetIcon(icon.FolderOpen)
	fs.list.Add(prev)
	// Adds directory files
	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			item := gui.NewImageLabel(files[i].Name())
			item.SetIcon(icon.FolderOpen)
			fs.list.Add(item)
		} else {
			if fs.filter != "" {
				if strings.HasSuffix(files[i].Name(), fs.filter[1:]) {
					item := gui.NewImageLabel(files[i].Name())
					item.SetIcon(icon.Note)
					fs.list.Add(item)
				}
			}
		}
	}
	return nil
}

// Selected returns the full path of the user-selected file
func (fs *FileSelect) Selected() string {
	if fs.name.Text() == "" {
		return ""
	}
	return filepath.Join(fs.path.Text(), fs.name.Text())
}

func (fs *FileSelect) onSelect() {

	// Get selected image label and its txt
	sel := fs.list.Selected()[0]
	label := sel.(*gui.ImageLabel)
	text := label.Text()

	// Checks if previous directory
	if text == ".." {
		dir, _ := filepath.Split(fs.path.Text())
		fs.setPath(filepath.Dir(dir))
		fs.name.SetText("")
		return
	}

	// Checks if it is a directory
	path := filepath.Join(fs.path.Text(), text)
	s, err := os.Stat(path)
	if err != nil {
		panic(err) // FIXME don't panic!
	}
	if s.IsDir() {
		fs.setPath(path)
		fs.name.SetText("")
	} else {
		fs.name.SetText(s.Name())
	}
}

// For sorting array of FileInfo by Name
type listFileInfo []os.FileInfo

func (fi listFileInfo) Len() int      { return len(fi) }
func (fi listFileInfo) Swap(i, j int) { fi[i], fi[j] = fi[j], fi[i] }
func (fi listFileInfo) Less(i, j int) bool {

	return fi[i].Name() < fi[j].Name()
}
