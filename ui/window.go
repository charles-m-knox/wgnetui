package ui

import (
	"fmt"

	"wgnetui/constants"
	"wgnetui/database"
)

func SetWindowTitleOnFileOpen() {
	if W == nil {
		return
	}

	(*W).SetTitle(fmt.Sprintf("%v: %v",
		constants.DefaultWindowTitle,
		database.OpenedFileName,
	))
}
