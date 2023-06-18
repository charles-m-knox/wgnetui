package ui

import (
	"fmt"

	"wgnetui/embedded"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func GetAboutView() (*fyne.Container, error) {
	aboutFileName := "about.md"
	aboutMarkdown, err := embedded.ReadFile(aboutFileName)
	if err != nil {
		return nil, fmt.Errorf(
			"error reading embedded file about.md: %v",
			err.Error(),
		)
	}
	aboutText := widget.NewRichTextFromMarkdown(aboutMarkdown)
	aboutText.Wrapping = fyne.TextWrapWord
	aboutView := container.NewMax(container.NewScroll(aboutText))

	return aboutView, nil
}
