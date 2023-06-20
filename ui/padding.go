package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// GetPadding returns a container whose minimum size is configurable. Use it
// with something like an HBox or VBox for simple padding. The container has no
// color (0 alpha).
func GetPadding(minWidth, minHeight float32) fyne.CanvasObject {
	// rect := canvas.NewRectangle(&color.RGBA{R: 0, G: 0, B: 0, A: 0})
	// rect.Resize(fyne.NewSize(minWidth, minHeight))
	// rectContainer := container.NewWithoutLayout(rect)
	// rectContainer.Resize(fyne.NewSize(minWidth, minHeight))
	// rect := canvas.NewRectangle(&color.NRGBA{128, 128, 128, 255})
	rect := canvas.NewRectangle(&color.NRGBA{128, 128, 128, 0})
	rect.SetMinSize(fyne.NewSize(minWidth, minHeight))
	return rect
}
