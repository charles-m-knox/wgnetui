package queries

import (
	"fmt"

	"wgnetui/constants"
	"wgnetui/database"
	"wgnetui/models"

	"gorm.io/gorm"
)

func GetServer() (server models.WgConfig, err error) {
	if database.DB == nil {
		return server, fmt.Errorf(constants.ErrorMessageNoDB)
	}

	serverResult := database.DB.First(&server, "is_server = ?", true)
	if serverResult.Error != nil && serverResult.Error != gorm.ErrRecordNotFound {
		return server, fmt.Errorf(
			"failed to look up server: %v",
			serverResult.Error.Error(),
		)
	}

	return server, nil
}
