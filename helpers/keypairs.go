package helpers

// import (
// 	"fmt"

// 	"wgnetui/constants"
// 	"wgnetui/database"
// 	"wgnetui/models"

// 	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
// 	"gorm.io/gorm"
// )

// // GetKeyPair will either retrieve the corresponding keypair for the provided
// // device ID or it will generate a new key pair for it.
// func GetKeyPair(deviceID uint) (result models.KeyPair, err error) {
// 	// retrieve KeyPair from database first, if it exists.
// 	// if it doesn't exist, then create it and save it to the DB

// 	if database.DB == nil {
// 		return result, fmt.Errorf(constants.ErrorMessageNoDB)
// 	}

// 	d := database.DB.First(&result, "id = ?", deviceID)
// 	if d.Error != nil {
// 		if d.Error == gorm.ErrRecordNotFound {
// 			privKey, err := wgtypes.GeneratePrivateKey()
// 			if err != nil {
// 				return result, fmt.Errorf("error generating private key: %v", err.Error())
// 			}

// 			pubKey := privKey.PublicKey()
// 			result.ID = deviceID
// 			result.Public = privKey.String()
// 			result.Private = pubKey.String()

// 			c := database.DB.Create(&result)
// 			if c.Error != nil {
// 				return result, fmt.Errorf(
// 					"error creating new keypair: %v",
// 					err.Error(),
// 				)
// 			}
// 		} else {
// 			return result, fmt.Errorf(
// 				"error finding keypair by device id %v: %v",
// 				deviceID,
// 				d.Error.Error(),
// 			)
// 		}
// 	}

// 	return result, nil
// }
