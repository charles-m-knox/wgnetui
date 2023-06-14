package models

import (
	"gorm.io/gorm"
)

// KeyPair represents a public and private key pair. For now, private keys
// are stored unencrypted.
// type KeyPair struct {
// 	gorm.Model
// 	ID       uint `gorm:"primary_key"`
// 	Public   string
// 	Private  string
// 	WgConfig []WgConfig `gorm:"foreignKey:KeyPairID"`
// 	Device   []Device   `gorm:"foreignKey:KeyPairID"`
// }

// WgConfig represents a generated wireguard configuration for a single peer/server
type WgConfig struct {
	gorm.Model
	ID                  uint   `gorm:"primary_key"`
	Config              string // auto-generated
	Name                string // user-configurable
	Description         string // user-configurable
	Extra               string // user-configurable
	IP                  string // user-configurable
	AllowedIPs          string // user-configurable
	PersistentKeepAlive uint   // user-configurable
	MTU                 uint16 // user-configurable
	Endpoint            string // user-configurable
	EndpointPort        uint16 // user-configurable
	DNS                 string // user-configurable
	IsServer            bool   // not editable; determined by the GenerationForm
	PrivateKey          string
	PublicKey           string
	PreSharedKey        string
}

// GenerationForm represents a user-submitted form
type GenerationForm struct {
	gorm.Model
	ID                       uint `gorm:"primary_key"`
	CIDR                     string
	DNS                      string
	Server                   string // ip address of the server within CIDR
	ServerInterface          string // eth0, eno1, etc
	Endpoint                 string
	EndpointPort             uint16 // publicly exposed wireguard server port
	MTU                      uint16
	AllowedIPs               string
	PersistentKeepAlive      uint   // if 0, do not set
	Name                     string // for setting a placeholder name for peers
	Description              string // for setting a placeholder desc. for peers
	Extra                    string // extra interface lines for peers
	RegenerateKeys           bool
	ResetAll                 bool // if true, deletes everything
	ForceAllowedIPs          bool // replaces all previous values if true
	ForcePersistentKeepAlive bool // replaces all previous values if true
	ForceMTU                 bool // replaces all previous values if true
	ForceEndpoint            bool // replaces all previous values if true
	ForceEndpointPort        bool // replaces all previous values if true
	ForceDNS                 bool // replaces all previous values if true
	ForceName                bool // replaces all previous values if true
	ForceDescription         bool // replaces all previous values if true
	ForceExtra               bool // replaces all previous values if true
}

// Device represents any device or interface that will maintain a Wireguard
// connection, such as a phone, laptop, NIC (such as eth0), etc.
// type Device struct {
// 	gorm.Model
// 	ID   uint `gorm:"primary_key"`
// 	Name string
// }
