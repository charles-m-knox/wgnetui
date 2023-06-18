package ui

import (
	"log"

	"wgnetui/database"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func GetMainMenu() *fyne.MainMenu {
	openAction := func() {
		log.Println("opened")
		openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(
					err,
					*W,
				)
				return
			}

			database.OpenedFileName = reader.URI().Name()
			database.OpenedFilePath = reader.URI().Path()

			database.Connect(database.OpenedFilePath)

			ExecuteOpenFunctions()

			// fileContent, err := io.ReadAll(reader)
			// if err != nil {
			// 	dialog.ShowError(
			// 		fmt.Errorf(
			// 			"Failed to open file: %v",
			// 			err.Error(),
			// 		),
			// 		w,
			// 	)
			// 	return
			// }

			// fmt.Println(
			// 	"File contents:",
			// 	string(fileContent),
			// )
		}, *W)
		openDialog.Resize(fyne.NewSize(600, 600))
		openDialog.Show()
	}

	// exitAction := func() {
	// 	a.Quit()
	// }

	aboutAction := func() {
		dialog.ShowInformation(
			"About",
			"This is an example Fyne app.",
			*W,
		)
	}

	// Create "File" menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open", openAction),
		// fyne.NewMenuItemSeparator(),
		// fyne.NewMenuItem("Exit", exitAction),
	)

	// Create "Help" menu
	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", aboutAction),
	)

	// Add the menus to a main menu bar
	mainMenu := fyne.NewMainMenu(fileMenu, helpMenu)

	return mainMenu
}
