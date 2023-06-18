package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

func AddQuitShortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyQ,
			Modifier: fyne.KeyModifierControl,
		},
		func(shortcut fyne.Shortcut) {
			log.Println("quitting")
			(*A).Quit()
		},
	)
}

func AddSwitchTab1Shortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key1,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if Tabs == nil {
				return
			}

			Tabs.SelectIndex(0)
		},
	)
}

func AddSwitchTab2Shortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key2,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if Tabs == nil {
				return
			}

			Tabs.SelectIndex(1)
		},
	)
}

func AddSwitchTab3Shortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key3,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if Tabs == nil {
				return
			}

			Tabs.SelectIndex(2)
		},
	)
}

func AddSwitchTab4Shortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key4,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if Tabs == nil {
				return
			}

			Tabs.SelectIndex(3)
		},
	)
}

func AddSwitchTab5Shortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key5,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if Tabs == nil {
				return
			}

			Tabs.SelectIndex(4)
		},
	)
}

func AddSwitchTab6Shortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.Key6,
			Modifier: fyne.KeyModifierAlt,
		},
		func(shortcut fyne.Shortcut) {
			if Tabs == nil {
				return
			}

			Tabs.SelectIndex(5)
		},
	)
}

func AddNextTabShortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyTab,
			Modifier: fyne.KeyModifierControl,
		},
		func(shortcut fyne.Shortcut) {
			if Tabs == nil {
				return
			}

			tabCount := len(Tabs.Items)
			newIndex := Tabs.SelectedIndex() + 1
			if newIndex > tabCount-1 {
				newIndex = 0
			}

			Tabs.SelectIndex(newIndex)
		},
	)
}

func AddPrevTabShortcut() {
	(*W).Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyTab,
			Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift,
		},
		func(shortcut fyne.Shortcut) {
			if Tabs == nil {
				return
			}

			newIndex := Tabs.SelectedIndex() - 1
			if newIndex < 0 {
				newIndex = len(Tabs.Items) - 1
			}

			Tabs.SelectIndex(newIndex)
		},
	)
}

// AddCtrlSShortcut allows different Ctrl+S behavior depending on which
// tab you're currently on.
func AddCtrlSShortcut() {
	(*W).Canvas().AddShortcut(
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

// AddCtrlRShortcut allows different Ctrl+R behavior depending on which
// tab you're currently on.
func AddCtrlRShortcut() {
	(*W).Canvas().AddShortcut(
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
