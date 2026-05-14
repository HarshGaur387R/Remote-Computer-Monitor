package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ReadOnlyEntry struct {
	widget.Entry
}

func NewReadOnlyMultiLineEntry() *ReadOnlyEntry {
	e := &ReadOnlyEntry{}
	e.MultiLine = true
	e.ExtendBaseWidget(e)
	return e
}

// Block all text input
func (e *ReadOnlyEntry) TypedRune(_ rune) {}

// Block backspace, delete, enter, etc.
func (e *ReadOnlyEntry) TypedKey(_ *fyne.KeyEvent) {}

func (e *ReadOnlyEntry) TypedShortcut(shortcut fyne.Shortcut) {
	// Allow copy and select-all, block everything else
	switch shortcut.(type) {
	case *fyne.ShortcutCopy, *fyne.ShortcutSelectAll:
		e.Entry.TypedShortcut(shortcut)
	default:
		// blocks paste, cut, undo, redo, etc.
	}
}
