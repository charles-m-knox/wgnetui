package generator

import (
	"fmt"
	"net"
	"strings"

	"wgnetui/constants"
	"wgnetui/database"
	"wgnetui/helpers"
	"wgnetui/models"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gorm.io/gorm"
)

// applySoftRules allows non-empty values to be preserved for existing peers,
// and desired values will be set for all other empty peer values
func applySoftRules(conf models.GenerationForm, w *models.WgConfig) error {
	if w == nil {
		return fmt.Errorf("received nil w ptr when applying soft rules")
	}

	if w.Name == "" {
		w.Name = conf.Name
		w.Name = strings.ReplaceAll(w.Name, "${id}", fmt.Sprintf("%v", w.ID))
	}

	if w.Description == "" {
		w.Description = conf.Description
		w.Description = strings.ReplaceAll(w.Description, "${id}", fmt.Sprintf("%v", w.ID))
		w.Description = strings.ReplaceAll(w.Description, "${name}", fmt.Sprintf("%v", w.Name))
	}

	if w.Extra == "" {
		w.Extra = conf.Extra
		w.Extra = strings.ReplaceAll(w.Extra, "${id}", fmt.Sprintf("%v", w.ID))
		w.Extra = strings.ReplaceAll(w.Extra, "${name}", fmt.Sprintf("%v", w.Name))
	}

	if w.AllowedIPs == "" {
		w.AllowedIPs = constants.DefaultAllowedIPs
	}

	if w.DNS == "" {
		w.DNS = conf.DNS
	}

	if w.PersistentKeepAlive == 0 {
		w.PersistentKeepAlive = conf.PersistentKeepAlive
	}

	if w.MTU == 0 {
		w.MTU = constants.DefaultMTU
	}

	if w.Endpoint == "" {
		w.Endpoint = conf.Endpoint
	}

	if w.EndpointPort == 0 {
		w.EndpointPort = conf.EndpointPort
	}

	return nil
}

// applyKeyRules generates and assigns Wireguard public, private, and pre-shared
// keys
func applyKeyRules(conf models.GenerationForm, w *models.WgConfig) error {
	if w == nil {
		return fmt.Errorf("received nil w ptr when applying key rules")
	}

	if w.PrivateKey == "" || w.PublicKey == "" || conf.RegenerateKeys {
		privKey, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return fmt.Errorf(
				"error generating private key: %v",
				err.Error(),
			)
		}
		w.PrivateKey = privKey.String()
		w.PublicKey = privKey.PublicKey().String()
	}

	if w.PreSharedKey == "" || conf.RegenerateKeys {
		w.PreSharedKey = helpers.GeneratePreSharedKey()
	}

	return nil
}

// applyForcedRules ensures that rules to override all other values are obeyed,
// for example, forceably setting the AllowedIPs to the value specified in
// the generation form.
func applyForcedRules(conf models.GenerationForm, w *models.WgConfig) error {
	if w == nil {
		return fmt.Errorf("received nil w ptr when applying forced rules")
	}

	if conf.ForceAllowedIPs {
		w.AllowedIPs = conf.AllowedIPs
	}
	if conf.ForcePersistentKeepAlive {
		w.PersistentKeepAlive = conf.PersistentKeepAlive
	}
	if conf.ForceMTU {
		w.MTU = conf.MTU
	}
	if conf.ForceEndpoint {
		w.Endpoint = conf.Endpoint
	}
	if conf.ForceEndpointPort {
		w.EndpointPort = conf.EndpointPort
	}
	if conf.ForceDNS {
		w.DNS = conf.DNS
	}
	if conf.ForceName {
		w.Name = conf.Name
	}
	if conf.ForceDescription {
		w.Description = conf.Description
	}
	if conf.ForceExtra {
		w.Extra = conf.Extra
	}

	return nil
}

func validateServer(server models.WgConfig) error {
	// assert that the server.ID is a non-zero value - the lowest possible ID
	// is 1
	if server.ID == 0 {
		return fmt.Errorf("invalid configured server ID 0")
	}

	// do some other basic data quality checks
	if server.PublicKey == "" || server.PrivateKey == "" {
		return fmt.Errorf("configured server has empty public/private key")
	}
	if server.IP == "" {
		return fmt.Errorf("configured server has empty IP address")
	}

	return nil
}

func generateConfig(w models.WgConfig, server models.WgConfig) (string, error) {
	persistentKeepAlive := ""
	if w.PersistentKeepAlive > 0 {
		persistentKeepAlive = "\nPersistentKeepAlive = 25"
	}
	extra := ""
	if w.Extra != "" {
		extra = fmt.Sprintf("\n%v", w.Extra)
	}

	return fmt.Sprintf(`[Interface]%v
PrivateKey = %s
Address = %s/32
DNS = %s
MTU = %v

[Peer]
PublicKey = %s
PresharedKey = %s
Endpoint = %s:%v
AllowedIPs = %v%v
`,
		extra,
		w.PrivateKey,
		w.IP,
		w.DNS,
		w.MTU,
		server.PublicKey,
		w.PreSharedKey,
		w.Endpoint,
		w.EndpointPort,
		w.AllowedIPs,
		persistentKeepAlive,
	), nil
}

func generateServerConfig(
	server models.WgConfig,
	conf models.GenerationForm,
	serverPeers []string,
	network *net.IPNet,
) string {
	maskSize, _ := network.Mask.Size()
	extra := ""
	if server.Extra != "" {
		extra = fmt.Sprintf("\n%v", server.Extra)
	}
	config := fmt.Sprintf(`[Interface]%v
PrivateKey = %s
Address = %s/%d
ListenPort = %v
# DNS = %s
MTU = %v
PostUp = iptables -A FORWARD -i %%i -j ACCEPT; iptables -A FORWARD -o %%i -j ACCEPT; iptables -t nat -A POSTROUTING -o %v -j MASQUERADE
PostDown = iptables -D FORWARD -i %%i -j ACCEPT; iptables -D FORWARD -o %%i -j ACCEPT; iptables -t nat -D POSTROUTING -o %v -j MASQUERADE

`,
		extra,
		server.PrivateKey,
		server.IP,
		maskSize,
		conf.EndpointPort,
		conf.DNS,
		conf.MTU,
		conf.ServerInterface,
		conf.ServerInterface,
	)

	for _, serverPeer := range serverPeers {
		config += serverPeer + "\n"
	}

	return config
}

func Generate(
	conf models.GenerationForm,
	serverIP net.IP,
	firstIP net.IP,
	network *net.IPNet,
	showGeneratorProgressDialog *func(),
	hideGeneratorProgressDialog *func(),
	setGeneratorProgressLabel *func(label string),
	setGeneratorProgressValue *func(v float64),
) (err error) {
	if showGeneratorProgressDialog == nil ||
		hideGeneratorProgressDialog == nil ||
		setGeneratorProgressLabel == nil ||
		setGeneratorProgressValue == nil {
		return fmt.Errorf("generate: progress bar dialog nil pointers")
	}

	// First, generate keypairs for every possible IP address within the range,
	// and while doing this, take note of which of them corresponds to the
	// IP address within the range that equals the server's IP address.

	// 1. Generate all possible WgConfig values, ensuring that old name/
	//    description values are preserved, if possible.

	// query the db for an existing wgconfig
	if database.DB == nil {
		return fmt.Errorf(
			constants.ErrorMessageNoDB,
		)
	}

	// reset the database if requested
	if conf.ResetAll {
		(*setGeneratorProgressLabel)("Resetting database...")
		resetResult := database.DB.Exec("DELETE FROM wg_configs WHERE id > 0")
		if resetResult.Error != nil {
			return fmt.Errorf(
				"failed to delete existing configs: %v",
				resetResult.Error.Error(),
			)
		}
	}

	// edge case: All values in the database above the generated IP range
	// need to be cleared out - in particular, they need to have their
	// IsServer flag cleared. Start by doing this first.
	(*setGeneratorProgressLabel)("Old server flag flip...")
	(*setGeneratorProgressValue)(5)
	resetResult := database.DB.Model(&models.WgConfig{}).Where(
		"is_server = ?", true,
	).Update(
		"is_server", false,
	)
	if resetResult.Error != nil {
		return fmt.Errorf(
			"failed to prepare pre-existing configs: %v",
			resetResult.Error.Error(),
		)
	}

	// take note of the total number of wireguard configs to generate so we
	// can accurately render a progress bar readout.
	wgs := helpers.EstimateNetworkSize(network)
	generatingMany := wgs > 512

	ip := firstIP
	i := uint(1)
	for {
		// skip over IP addresses that end with .0
		// ip4 := ip.To4()
		// if ip4[3] == 0 {
		// 	ip = helpers.NextIP(ip)
		// 	i++
		// }

		if !network.Contains(ip) {
			break
		}

		// performance slowdowns start to occur when updating the progress
		// dialog for large numbers
		if (!generatingMany) || (generatingMany && i%50 == 0) {
			(*setGeneratorProgressLabel)(
				fmt.Sprintf(
					"[%v/%v]: Generating keys and configuring peers: Step 1/2",
					i,
					wgs,
				),
			)

			// At the start of the loop we are 5% done with the full 100%.
			// We want to get an additional 90% done, leaving
			// us with 5% remaining. 90% / 2 full iterations of the IP range = 45%,
			// meaning that (1 / wgs) / 45%.
			// p1 := float64(i) / float64(wgs)
			// p2 := (float64(i) / float64(wgs)) * (90 / 2.0)
			// log.Printf("i=%v, wgs=%v, %v %v", float64(i), float64(wgs), p1, p2)
			(*setGeneratorProgressValue)(
				5 + (float64(i)/float64(wgs))*(90/2.0),
			)
		}

		w := models.WgConfig{}

		result := database.DB.First(&w, "id = ?", i)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return fmt.Errorf(
				"failed to look up config %v: %v",
				i,
				result.Error.Error(),
			)
		}

		// update values. Note that in general, if a value in the form is
		// left blank, the original value will be preserved where possible.
		w.ID = i
		w.IP = ip.String()

		err = applySoftRules(conf, &w)
		if err != nil {
			return err
		}
		err = applyForcedRules(conf, &w)
		if err != nil {
			return err
		}
		err = applyKeyRules(conf, &w)
		if err != nil {
			return err
		}

		// detect if this is the server's IP address
		if ip.Equal(serverIP) {
			// found the server ip, set IsServer to true
			w.IsServer = true
		} else {
			// set everything else to false for completeness
			w.IsServer = false
		}

		// save the newly updated wg config (peer) to the database
		saved := database.DB.Save(&w)
		if saved.Error != nil {
			return fmt.Errorf(
				"failed to save config %v: %v",
				i,
				saved.Error.Error(),
			)
		}

		ip = helpers.NextIP(ip)
		i++
	}

	// 2. Now that the entire sequence of configs has been updated, we can
	//    proceed to generate all peer configs. Iterate through the whole
	//    list again (ignoring the server), generating configs for each peer.
	//    Each peer needs to know the server's public key and IP address, so
	//    start by finding the server:
	server := models.WgConfig{}
	serverResult := database.DB.First(&server, "is_server = ?", true)
	if serverResult.Error != nil && serverResult.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf(
			"failed to look up server: %v",
			serverResult.Error.Error(),
		)
	}

	err = validateServer(server)
	if err != nil {
		return err
	}

	// proceed with the remainder of step 2:
	// Generate all peer configs, now that we have the server IP.
	serverPeers := []string{}
	ip = firstIP // reset the IP to the original ipRange
	i = uint(1)  // reset this too
	for {
		// skip over IP addresses that end with .0
		// ip4 := ip.To4()
		// if ip4[3] == 0 {
		// 	ip = helpers.NextIP(ip)
		// 	i++
		// }

		if !network.Contains(ip) {
			break
		}

		// performance slowdowns start to occur when updating the progress
		// dialog for large numbers
		if (!generatingMany) || (generatingMany && i%50 == 0) {
			(*setGeneratorProgressLabel)(
				fmt.Sprintf(
					"[%v/%v]: Generating peer configs: Step 2/2",
					i,
					wgs,
				),
			)

			// At the start of the loop we are 5% done with the full 100%.
			// We want to get an additional 90% done, leaving
			// us with 5% remaining. 90% / 2 full iterations of the IP range = 45%,
			// meaning that (1 / wgs) / 45%.
			(*setGeneratorProgressValue)(
				5 + 45 + float64(i)/float64(wgs)*(90/2.0),
			)
		}

		w := models.WgConfig{}

		// retrieve the peer for this ID
		result := database.DB.First(&w, "id = ?", i)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return fmt.Errorf(
				"failed to look up config for peer conf gen %v: %v",
				i,
				result.Error.Error(),
			)
		}

		// assert that the IP addresses are equal
		if ip.String() != w.IP {
			return fmt.Errorf(
				"ip address mismatch: %v vs. %v",
				ip.String(),
				w.IP,
			)
		}

		// ignore the server when it shows up in this list
		if w.IP == conf.Server || w.IsServer {
			ip = helpers.NextIP(ip)
			i++
			continue
		}

		// generate the peer config for this peer
		peerConf, err := generateConfig(w, server)
		if err != nil {
			return fmt.Errorf(
				"error generating config for client %v: %v\n",
				i,
				err,
			)
		}

		w.Config = peerConf

		serverPeer := fmt.Sprintf(
			"[Peer]\nPublicKey = %s\nAllowedIPs = %s/32\nPresharedKey = %s\n",
			w.PublicKey,
			w.IP,
			w.PreSharedKey,
		)
		serverPeers = append(serverPeers, serverPeer)

		// save the config to the database
		saved := database.DB.Save(&w)
		if saved.Error != nil {
			return fmt.Errorf(
				"failed to save config %v: %v",
				i,
				saved.Error.Error(),
			)
		}

		ip = helpers.NextIP(ip)
		i++
	}

	(*setGeneratorProgressLabel)(
		fmt.Sprintf(
			"[%v/%v]: Generation step 2/2 complete",
			i,
			wgs,
		),
	)
	(*setGeneratorProgressValue)(95)

	// 3. Finally, update the server config.
	(*setGeneratorProgressLabel)("Saving final server config")
	(*setGeneratorProgressValue)(99)
	server.Config = generateServerConfig(server, conf, serverPeers, network)
	saved := database.DB.Save(&server)
	if saved.Error != nil {
		return fmt.Errorf(
			"failed to save server config %v: %v",
			i,
			saved.Error.Error(),
		)
	}

	(*setGeneratorProgressLabel)("Done")
	(*setGeneratorProgressValue)(100)

	return nil
}
