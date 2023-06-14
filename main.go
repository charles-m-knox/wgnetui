package main

import (
	"fmt"
	"log"

	"wgnetui/constants"
	"wgnetui/database"
	"wgnetui/embedded"
	"wgnetui/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func init() {
}

func main() {
	a := app.New()
	w := a.NewWindow("Wireguard Network UI")

	database.Connect("test.db")

	filename := "about.md"
	aboutMarkdown, err := embedded.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading embedded file about.md: %v", err.Error())
		return
	}

	aboutText := widget.NewRichTextFromMarkdown(aboutMarkdown)
	aboutText.Wrapping = fyne.TextWrapWord
	// input3 := widget.NewRichTextFromMarkdown(constants.PlaceholderMarkdown)
	// input4 := widget.NewRichTextFromMarkdown(constants.PlaceholderMarkdown)
	// input5 := widget.NewRichTextFromMarkdown(constants.PlaceholderMarkdown)
	// input6 := widget.NewRichTextFromMarkdown(constants.PlaceholderMarkdown)
	// input7 := widget.NewRichTextFromMarkdown(constants.PlaceholderMarkdown)
	// input8 := widget.NewRichTextFromMarkdown(constants.PlaceholderMarkdown)

	generatorProgressBar := widget.NewProgressBar()
	generatorProgressBar.Min = 0
	generatorProgressBar.Max = 100
	generatorProgressBar.SetValue(0)

	generatorLabel := widget.NewLabel("Progressing")

	generatorProgressDialogContent := container.NewVBox(
		generatorLabel,
		generatorProgressBar,
	)

	generatorProgressDialog := dialog.NewCustom("Generating....", "Close", generatorProgressDialogContent, w)
	showGeneratorProgressDialog := func() {
		generatorProgressDialog.Show()
	}
	hideGeneratorProgressDialog := func() {
		generatorProgressDialog.Hide()
	}
	setGeneratorProgressLabel := func(label string) {
		if generatorLabel == nil {
			log.Printf(
				"setProgressLabel generatorLabel was nil for message: %v",
				label,
			)
			return
		}

		generatorLabel.SetText(label)
	}
	setGeneratorProgressValue := func(v float64) {
		if generatorProgressBar == nil {
			log.Printf(
				"setGeneratorProgressValue generatorProgressBar was nil for value: %v",
				v,
			)
			return
		}

		generatorProgressBar.SetValue(v)
	}

	aboutView := container.NewScroll(aboutText)
	genform := ui.GetWgGenForm(
		w,
		&showGeneratorProgressDialog,
		&hideGeneratorProgressDialog,
		&setGeneratorProgressLabel,
		&setGeneratorProgressValue,
	)
	genformView := container.NewScroll(genform)
	devicesView, err := ui.GetDevicesView(w)
	if err != nil {
		log.Fatalf("failed to get devices view: %v", err.Error())
	}
	// confViewerView := container.NewHSplit(input3, input4)
	// keysView := container.NewHSplit(input5, input6)
	// statsView := container.NewHSplit(input7, input8)

	ui.AddQuitShortcut(&w, &a)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("About", theme.FileIcon(), aboutView),
		container.NewTabItemWithIcon("Generator", theme.FileIcon(), genformView),
		container.NewTabItemWithIcon(constants.TabDevices, theme.DocumentIcon(), devicesView),
		// container.NewTabItemWithIcon("Config Viewer", theme.InfoIcon(), confViewerView),
		// container.NewTabItemWithIcon("Keys", theme.InfoIcon(), keysView),
		// container.NewTabItemWithIcon("Stats", theme.InfoIcon(), statsView),
	)

	tabs.OnSelected = func(tab *container.TabItem) {
		switch tab.Text {
		case constants.TabDevices:
			// refresh the devices view when switching tabs, if desired
			devicesView, err := ui.GetDevicesView(w)
			if err != nil {
				dialog.ShowError(
					fmt.Errorf(
						"Failed to refresh devices view: %v",
						err.Error(),
					),
					w,
				)
			}
			tab.Content = devicesView
			break
		default:
			return
		}
	}

	ui.AddSwitchTab1Shortcut(&w, tabs)
	ui.AddSwitchTab2Shortcut(&w, tabs)
	ui.AddSwitchTab3Shortcut(&w, tabs)
	// ui.AddSwitchTab4Shortcut(&w, tabs)
	// ui.AddSwitchTab5Shortcut(&w, tabs)
	ui.AddNextTabShortcut(&w, tabs)
	ui.AddPrevTabShortcut(&w, tabs)

	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1024, 768))
	w.CenterOnScreen()
	w.ShowAndRun()
}
