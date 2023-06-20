package ui

import (
	"log"

	"fyne.io/fyne/v2"
)

// InitializeUI sets up keyboard shortcuts and many other things that are
// required in order to make the app work.
func InitializeUI() {
	if W == nil || A == nil {
		log.Fatalf("failed to initialize ui: window or app nil ptr")
	}

	// SetTheme()

	AddQuitShortcut()
	AddCtrlSShortcut()
	AddCtrlRShortcut()

	lav := func() {
		log.Println("running LoadAboutView")
		var err error
		AboutView, err = GetAboutView()
		if err != nil {
			log.Fatalf("failed to get about view: %v", err.Error())
		}
		if TabAbout != nil {
			TabAbout.Content = AboutView
			TabAbout.Content.Refresh()
		} else {
			log.Println("TabAbout is nil")
		}
	}
	LoadAboutView = &lav

	lgfv := func() {
		log.Println("running LoadGenFormView")
		GenFormView = GetWgGenForm()
		if TabGenerator != nil {
			TabGenerator.Content = GenFormView
			TabGenerator.Content.Refresh()
		} else {
			log.Println("TabGenerator is nil")
		}
	}
	LoadGenFormView = &lgfv

	var err error

	ldfv := func() {
		log.Println("running LoadDevicesView")
		DevicesView, err = GetDevicesView()
		if err != nil {
			log.Fatalf("failed to get devices view: %v", err.Error())
		}
		if TabDevices != nil {
			TabDevices.Content = DevicesView
			TabDevices.Content.Refresh()
		} else {
			log.Println("TabDevices is nil")
		}
	}
	LoadDevicesView = &ldfv

	// allows for the title of the window to be updated later upon refresh
	// hooks
	setWindowTitleOnFileOpen := func() {
		SetWindowTitleOnFileOpen()
	}

	// update the title of the window for the initial app startup only
	setWindowTitleOnFileOpen()

	OnOpenFunctions = append(
		OnOpenFunctions,
		LoadAboutView,
		LoadGenFormView,
		LoadDevicesView,
		&setWindowTitleOnFileOpen,
	)

	// we technically haven't executed all the functions relating to
	// opening a new file yet, so let's do that now
	ExecuteOpenFunctions()

	SetupTabs()
	SetupTabsShortcuts()

	menu := GetMainMenu()
	(*W).SetMainMenu(menu)
	(*W).SetContent(Tabs)
	(*W).Resize(fyne.NewSize(1024, 768))
	(*W).CenterOnScreen()
}
