package main

import (
	"flag"

	"wgnetui/constants"
	"wgnetui/database"
	"wgnetui/ui"

	"fyne.io/fyne/v2/app"
)

func init() {
	flag.StringVar(
		&database.OpenedFilePath,
		"db",
		"",
		"Path to the database file",
	)
	flag.Parse()
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
