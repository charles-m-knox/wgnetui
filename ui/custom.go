package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type CustomEntry struct {
	widget.Entry
}

func (e *CustomEntry) TypedShortcut(shortcut fyne.Shortcut) {
	switch shortcut := shortcut.(type) {
	// case *fyne.ShortcutCopy:
	// 	// If you want to ignore the copy shortcut
	// 	log.Println("shortcut aborted")
	// 	return
	default:
		// Call the parent's implementation for all other shortcuts
		e.Entry.TypedShortcut(shortcut)
	}
}

// GetCustomEntry returns an Entry widget that, for some reason, does not
// interrupt mouse wheel scroll events.
func GetCustomEntry() *CustomEntry {
	entry := &CustomEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

type CustomMultiLineEntry struct {
	widget.Entry
}

func (e *CustomMultiLineEntry) TypedShortcut(shortcut fyne.Shortcut) {
	switch shortcut := shortcut.(type) {
	case *fyne.ShortcutCopy:
		// If you want to ignore the copy shortcut
		log.Println("shortcut aborted")
		return
	default:
		// Call the parent's implementation for all other shortcuts
		e.Entry.TypedShortcut(shortcut)
	}
}

func GetCustomMultiLineEntry() *CustomMultiLineEntry {
	entry := &CustomMultiLineEntry{}
	entry.ExtendBaseWidget(entry)
	entry.MultiLine = true
	entry.Wrapping = fyne.TextTruncate
	return entry
}
