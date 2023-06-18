package ui

import (
	"wgnetui/constants"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func SetupTabs() {
	TabAbout = container.NewTabItemWithIcon(
		constants.TabAbout,
		theme.FileIcon(),
		AboutView,
	)

	TabGenerator = container.NewTabItemWithIcon(
		constants.TabGenerator,
		theme.FileIcon(),
		GenFormView,
	)

	TabDevices = container.NewTabItemWithIcon(
		constants.TabDevices,
		theme.DocumentIcon(),
		DevicesView,
	)

	Tabs = container.NewAppTabs(
		TabAbout,
		TabGenerator,
		TabDevices,
	)
	Tabs.SetTabLocation(container.TabLocationTop)

	Tabs.OnSelected = func(tab *container.TabItem) {
		ActiveTab = tab.Text
		switch tab.Text {
		case constants.TabDevices:
			// refresh the devices view when switching tabs, if desired
			// devicesView, err := ui.GetDevicesView(w)
			// if err != nil {
			// 	dialog.ShowError(
			// 		fmt.Errorf(
			// 			"Failed to refresh devices view: %v",
			// 			err.Error(),
			// 		),
			// 		w,
			// 	)
			// }

			// tab.Content = devicesView // TODO: is this even necessary?
			break
		default:
			return
		}
	}
}

func SetupTabsShortcuts() {
	AddSwitchTab1Shortcut()
	AddSwitchTab2Shortcut()
	AddSwitchTab3Shortcut()
	AddNextTabShortcut()
	AddPrevTabShortcut()
}
