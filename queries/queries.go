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

// GetLatestGenForm retrieves the most recently created generation form configuration
// from the currently opened SQLite database.
func GetLatestGenForm() (gf models.GenerationForm, err error) {
	if database.DB == nil {
		return gf, fmt.Errorf(constants.ErrorMessageNoDB)
	}

	result := database.DB.Order("created_at desc").Limit(1).First(&gf)
	if result.Error != nil {
		return gf, fmt.Errorf(
			"Failed to retrieve previously stored configuration: %v",
			result.Error.Error(),
		)
	}

	return gf, nil
}
