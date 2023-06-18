package ui

import "fyne.io/fyne/v2/theme"

func SetTheme() {
	(*A).Settings().SetTheme(theme.DarkTheme())
}
