package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	// window & app

	// W is the globally accessible Fyne window
	W *fyne.Window
	// A is the globally accessible Fyne app
	A *fyne.App

	// tabs

	// Tabs are the primary viewing container for the entire application
	Tabs         *container.AppTabs
	TabAbout     *container.TabItem
	TabGenerator *container.TabItem
	TabDevices   *container.TabItem
	// ActiveTab is the name of the current active tab. Compare this value to
	// the constants.TabFoo values.
	ActiveTab string

	// functions

	// OnOpenFunctions stores functions that will be executed in order when a file is
	// opened
	OnOpenFunctions []*func()
	// LoadAboutView is a function that contains logic to refresh the About tab
	// view
	LoadAboutView *func()
	// LoadGenFormView is a function that contains logic to refresh the generation
	// form view
	LoadGenFormView *func()
	// LoadDevicesView is a function that contains logic to refresh the devices
	// view
	LoadDevicesView *func()

	// views

	// AboutView is what gets shown in the About tab
	AboutView *fyne.Container
	// AboutView is what gets shown in the Generator tab
	GenFormView *fyne.Container
	// AboutView is what gets shown in the Devices tab
	DevicesView *container.Split

	// keyboard shortcuts

	// All functions to run when pressing Ctrl+R
	CtrlRShortcuts map[string]*func()
	// All functions to run when pressing Ctrl+S
	CtrlSShortcuts map[string]*func()
	// All functions to run when pressing Ctrl+O
	CtrlOShortcuts map[string]*func()

	// ProgressBarDialog is meant to be used in many places and shows
	// progress
	ProgressBarDialog *ProgressDialog
)

func init() {
	CtrlRShortcuts = make(map[string]*func())
	CtrlSShortcuts = make(map[string]*func())
	CtrlOShortcuts = make(map[string]*func())
	OnOpenFunctions = make([]*func(), 0)
}

// InitializeGlobals is different from init() in that it must be called
// after a Fyne app has been initialized via a := app.New() then
// a.NewWindow().
func InitializeGlobals() {
	ProgressBarDialog = &ProgressDialog{
		ProgressBar:   widget.NewProgressBar(),
		ProgressLabel: widget.NewLabel(""),
	}

	ProgressBarDialog.ProgressBar.Max = 100
	ProgressBarDialog.ProgressBar.Min = 0
}

// ExecuteOpenFunctions iterates through the list of global functions that
// should be called when opening a file.
func ExecuteOpenFunctions() {
	for i := range OnOpenFunctions {
		if OnOpenFunctions[i] == nil {
			log.Printf("open fn %v is nil", i)
			continue
		}

		log.Printf("running open fn %v", i)

		(*OnOpenFunctions[i])()
	}
}
