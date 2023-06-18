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
	database.Connect(constants.DefaultFileName)
	database.OpenedFileName = constants.DefaultFileName
	ui.InitializeGlobals()
	ui.InitializeUI()

	// ui.TabAbout = container.NewTabItemWithIcon(
	// 	constants.TabAbout,
	// 	theme.FileIcon(),
	// 	aboutView,
	// )

	// ui.TabGenerator = container.NewTabItemWithIcon(
	// 	constants.TabGenerator,
	// 	theme.FileIcon(),
	// 	getNewGenForm(),
	// )

	// ui.TabDevices = container.NewTabItemWithIcon(
	// 	constants.TabDevices,
	// 	theme.DocumentIcon(),
	// 	devicesView,
	// )

	// tabs := container.NewAppTabs(
	// 	ui.TabAbout,
	// 	ui.TabGenerator,
	// 	ui.TabDevices,
	// )

	// tabs.OnSelected = func(tab *container.TabItem) {
	// 	ui.ActiveTab = tab.Text
	// 	switch tab.Text {
	// 	case constants.TabDevices:
	// 		// refresh the devices view when switching tabs, if desired
	// 		// devicesView, err := ui.GetDevicesView(w)
	// 		// if err != nil {
	// 		// 	dialog.ShowError(
	// 		// 		fmt.Errorf(
	// 		// 			"Failed to refresh devices view: %v",
	// 		// 			err.Error(),
	// 		// 		),
	// 		// 		w,
	// 		// 	)
	// 		// }
	// 		tab.Content = devicesView
	// 		break
	// 	default:
	// 		return
	// 	}
	// }

	w.ShowAndRun()
}
