package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
)

var (
	CtrlRShortcuts map[string]*func()
	CtrlSShortcuts map[string]*func()
)

func init() {
	CtrlRShortcuts = make(map[string]*func())
	CtrlSShortcuts = make(map[string]*func())
}

func AddQuitShortcut(w *fyne.Window, a *fyne.App) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyQ,
			Modifier: fyne.KeyModifierControl,
		},
		func(shortcut fyne.Shortcut) {
			log.Println("quitting")
			(*a).Quit()
		},
	)
}

func AddSwitchTab1Shortcut(w *fyne.Window, t *container.AppTabs) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key1,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if t == nil {
				return
			}

			t.SelectIndex(0)
		},
	)
}

func AddSwitchTab2Shortcut(w *fyne.Window, t *container.AppTabs) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key2,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if t == nil {
				return
			}

			t.SelectIndex(1)
		},
	)
}

func AddSwitchTab3Shortcut(w *fyne.Window, t *container.AppTabs) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key3,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if t == nil {
				return
			}

			t.SelectIndex(2)
		},
	)
}

func AddSwitchTab4Shortcut(w *fyne.Window, t *container.AppTabs) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key4,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if t == nil {
				return
			}

			t.SelectIndex(3)
		},
	)
}

func AddSwitchTab5Shortcut(w *fyne.Window, t *container.AppTabs) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key5,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if t == nil {
				return
			}

			t.SelectIndex(4)
		},
	)
}

func AddSwitchTab6Shortcut(w *fyne.Window, t *container.AppTabs) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key6,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if t == nil {
				return
			}

			t.SelectIndex(5)
		},
	)
}

func AddNextTabShortcut(w *fyne.Window, t *container.AppTabs) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyTab,
			Modifier: fyne.KeyModifierControl,
		},
		func(shortcut fyne.Shortcut) {
			if t == nil {
				return
			}

			tabCount := len(t.Items)
			newIndex := t.SelectedIndex() + 1
			if newIndex > tabCount-1 {
				newIndex = 0
			}

			t.SelectIndex(newIndex)
		},
	)
}

func AddPrevTabShortcut(w *fyne.Window, t *container.AppTabs) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyTab,
			Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
		},
		func(shortcut fyne.Shortcut) {
			if t == nil {
				return
			}

			newIndex := t.SelectedIndex() - 1
			if newIndex < 0 {
				newIndex = len(t.Items) - 1
			}

			t.SelectIndex(newIndex)
		},
	)
}

// AddCtrlSShortcut allows different Ctrl+S behavior depending on which
// tab you're currently on.
func AddCtrlSShortcut(w *fyne.Window) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyS,
			Modifier: fyne.KeyModifierControl,
		},
		func(shortcut fyne.Shortcut) {
			for _, fn := range CtrlSShortcuts {
				if fn == nil {
					log.Println("ctrl+S shortcut func is nil")
					return
				}

				(*fn)()
			}
		},
	)
}

// AddCtrlRShortcut allows different Ctrl+S behavior depending on which
// tab you're currently on.
func AddCtrlRShortcut(w *fyne.Window) {
	(*w).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyR,
			Modifier: fyne.KeyModifierControl,
		},
		func(shortcut fyne.Shortcut) {
			for _, fn := range CtrlRShortcuts {
				if fn == nil {
					log.Println("ctrl+R shortcut func is nil")
					return
				}

				(*fn)()
			}
		},
	)
}
