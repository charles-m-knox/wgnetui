package main

import (
	"wgnetui/constants"
	"wgnetui/database"
	"wgnetui/ui"

	"fyne.io/fyne/v2/app"
)

func init() {
}

func main() {
	a := app.New()
	w := a.NewWindow(constants.DefaultWindowTitle)
	ui.A = &a
	ui.W = &w
	database.Initialize()
	ui.InitializeGlobals()
	ui.InitializeUI()
	w.ShowAndRun()
}
