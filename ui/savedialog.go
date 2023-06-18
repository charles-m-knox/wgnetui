package ui

import (
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func SaveToFileDialog(content, fileName string, fileFilter []string) {
	openDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err == nil && writer == nil {
			return
		}
		if err != nil {
			dialog.ShowError(
				err,
				*W,
			)
			return
		}

		_, err = io.WriteString(writer, content)
		if err != nil {
			dialog.ShowError(
				fmt.Errorf(
					"Failed to open file: %v",
					err.Error(),
				),
				*W,
			)
			return
		}
	}, *W)
	openDialog.Resize(fyne.NewSize(600, 600))
	openDialog.SetFilter(storage.NewExtensionFileFilter(fileFilter))
	openDialog.SetFileName(fileName)
	openDialog.Show()
}
