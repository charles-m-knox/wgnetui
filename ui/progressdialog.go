package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ProgressDialog struct {
	ProgressBar   *widget.ProgressBar
	Dialog        dialog.Dialog
	DialogContent *fyne.Container
	ProgressLabel *widget.Label
	Title         string
	ButtonText    string
}

// Fresh resets the values of the dialog to the given values.
func (p *ProgressDialog) Fresh(title, buttonText, label string) {
	if p == nil {
		return
	}

	p.Title = title
	p.ButtonText = buttonText
	p.SetLabel(label)

	p.SetDialog()
}

// SetDialog updates the dialog to reflect the latest settings for the dialog
// itself. This is required in order to propagate a new title or new button
// text.
func (p *ProgressDialog) SetDialog() {
	if p == nil {
		return
	}

	p.Dialog = dialog.NewCustom(
		p.Title,
		p.ButtonText,
		p.GetContent(),
		*W,
	)
}

func (p *ProgressDialog) GetContent() *fyne.Container {
	return container.NewVBox(
		p.ProgressLabel,
		p.ProgressBar,
	)
}

// SetLabel updates the label that is shown above the progress bar.
func (p *ProgressDialog) SetLabel(l string) {
	if p == nil || p.ProgressLabel == nil {
		return
	}

	p.ProgressLabel.SetText(l)
}

// SetTitle requires rebuilding the dialog (internally) and is perhaps
// a heftier operation than simply changing the progress bar label.
func (p *ProgressDialog) SetTitle(t string) {
	if p == nil {
		return
	}

	p.Title = t

	p.SetDialog()
}

// SetButtonText requires rebuilding the dialog (internally) and is perhaps
// a heftier operation than simply changing the progress bar label.
func (p *ProgressDialog) SetButtonText(t string) {
	if p == nil {
		return
	}

	p.ButtonText = t

	p.SetDialog()
}

func (p *ProgressDialog) SetProgress(val float64) {
	if p == nil || p.ProgressBar == nil {
		return
	}

	p.ProgressBar.SetValue(val)
}

func (p *ProgressDialog) Show() {
	if p == nil {
		return
	}

	p.Dialog.Show()
}

func (p *ProgressDialog) Hide() {
	if p == nil {
		return
	}

	p.Dialog.Hide()
}
