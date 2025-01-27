package generator

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"sync/atomic"

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

	if w.DNS == "" && !w.IsServer { // don't set the dns for the server
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
	if conf.ForceDNS && !w.IsServer { // don't set the DNS for the server
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

func GenerateConfig(w models.WgConfig, server models.WgConfig) (string, error) {
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
	serverPeers string,
	// serverPeers []string,
	network *net.IPNet,
) string {
	maskSize, _ := network.Mask.Size()
	extra := ""
	if server.Extra != "" {
		extra = fmt.Sprintf("\n%v", server.Extra)
	}
	dns := ""
	if server.DNS != "" {
		dns = fmt.Sprintf("\nDNS = %v", server.DNS)
	}

	config := fmt.Sprintf(`[Interface]%v
PrivateKey = %s
Address = %s/%d
ListenPort = %v%v
MTU = %v
PostUp = iptables -A FORWARD -i %%i -j ACCEPT; iptables -A FORWARD -o %%i -j ACCEPT; iptables -t nat -A POSTROUTING -o %v -j MASQUERADE
PostDown = iptables -D FORWARD -i %%i -j ACCEPT; iptables -D FORWARD -o %%i -j ACCEPT; iptables -t nat -D POSTROUTING -o %v -j MASQUERADE

%v
`,
		extra,
		server.PrivateKey,
		server.IP,
		maskSize,
		conf.EndpointPort,
		dns,
		conf.MTU,
		conf.ServerInterface,
		conf.ServerInterface,
		serverPeers,
	)

	// for _, serverPeer := range serverPeers {
	// 	config += serverPeer + "\n"
	// }

	return config
}

type IPAddress struct {
	// the ip address as a string value
	S string
	// the IP address as a net.IP value
	IP net.IP
	// whether this IP is the server IP
	IsServerIP bool
}

func Generate(
	conf models.GenerationForm,
	serverIP net.IP,
	firstIP net.IP,
	network *net.IPNet,
	setProgressLabel *func(l string),
	setProgressValue *func(v float64),
) (err error) {
	// First, generate keypairs for every possible IP address within the range,
	// and while doing this, take note of which of them corresponds to the
	// IP address within the range that equals the server's IP address.

	// 1. Generate all possible WgConfig values, ensuring that old name/
	//    description values are preserved, if possible.

	// query the db for an existing wgconfig
	database.Reconnect() // reading the db between writes can cause failures
	if database.DB == nil {
		return fmt.Errorf(constants.ErrorMessageNoDB)
	}

	if setProgressLabel == nil || setProgressValue == nil {
		return fmt.Errorf(constants.ErrorMessageNoProgressBarDialog)
	}

	// reset the database if requested
	if conf.ResetAll {
		(*setProgressLabel)("Resetting database...")
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
	(*setProgressLabel)("Identifying previously generated server config...")
	(*setProgressValue)(5)
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
	log.Printf("expecting %v ip addresses", wgs)
	// generatingMany := wgs > 512

	// take note of every IP address that we have to operate on
	allIPs := []IPAddress{}
	ip := firstIP
	serverIPIndex := -1

	k := 0
	for {
		ipa := IPAddress{
			S:          ip.String(),
			IP:         ip,
			IsServerIP: false,
		}
		// if the IP address ends with .0 or .255, skip it
		if strings.HasSuffix(ipa.S, ".0") || strings.HasSuffix(ipa.S, ".255") {
			log.Println("skipping ipa:", ipa)
			ip = helpers.NextIP(ip)
			continue
		}
		if !network.Contains(ip) {
			break
		}

		if ipa.IP.Equal(serverIP) {
			ipa.IsServerIP = true
			serverIPIndex = k
		}

		allIPs = append(allIPs, ipa)
		ip = helpers.NextIP(ip)
		k++
	}

	ips := len(allIPs)

	log.Printf("total IP addresses: %v", ips)

	var progress int64 = 0
	var server *models.WgConfig
	var serverPeers strings.Builder
	// serverPeers := []string{}
	mutex := &sync.Mutex{}

	// prepDevice preps a single device, this is useful for processing
	// the server first before everything else. The logic at this step
	// is the same as all other devices though, only the server will behave
	// slightly different in a few spots.
	prepDevice := func(i uint, ip IPAddress) error {
		w := models.WgConfig{}

		result := database.DB.First(&w, "id = ?", i)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				err := fmt.Errorf(
					"did not find record by id %v: %v",
					i,
					result.Error.Error(),
				)
				log.Println(err.Error())
			} else {
				err := fmt.Errorf(
					"failed to look up config %v: %v",
					i,
					result.Error.Error(),
				)
				return err
			}
		}

		// update values. Note that in general, if a value in the form is
		// left blank, the original value will be preserved where possible.
		w.ID = i
		w.IP = ip.S
		w.IsServer = ip.IsServerIP

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

		if server != nil && !w.IsServer {
			// generate the peer config for this peer
			peerConf, err := GenerateConfig(w, *server)
			if err != nil {
				return fmt.Errorf(
					"error generating config for client %v: %v\n",
					i,
					err,
				)
			}

			w.Config = peerConf

			serverPeer := fmt.Sprintf(
				"[Peer]\nPublicKey = %s\nAllowedIPs = %s/32\nPresharedKey = %s\n\n",
				w.PublicKey,
				w.IP,
				w.PreSharedKey,
			)
			mutex.Lock()
			serverPeers.WriteString(serverPeer)
			mutex.Unlock()
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

		return nil
	}

	prepFn := func(wg *sync.WaitGroup, progress *int64, i uint, ip IPAddress) error {
		defer wg.Done()

		err := prepDevice(i, ip)
		if err != nil {
			return err
		}

		atomic.AddInt64(progress, 1)

		msg := fmt.Sprintf(
			"[%v/%v]: Generating keys and configuring peers: Step 1/2",
			atomic.LoadInt64(progress),
			ips,
		)
		// log.Printf(msg)
		(*setProgressLabel)(msg)
		(*setProgressValue)(
			5 + (float64(atomic.LoadInt64(progress))/float64(ips))*(90),
		)
		return nil
	}

	prepFnWrapped := func(wg *sync.WaitGroup, progress *int64, i uint, ip IPAddress) {
		err := prepFn(wg, progress, i, ip)
		if err != nil {
			log.Fatalf("error generating keys and configuring peers: %v", err.Error())
		}
	}

	// generate the server first - this allows more parallel processing to be
	// done in one step
	if serverIPIndex >= 0 && allIPs[serverIPIndex].IsServerIP {
		err = prepDevice(uint(serverIPIndex)+1, allIPs[serverIPIndex])
		if err != nil {
			return fmt.Errorf(
				"failed to write server: %v",
				err.Error(),
			)
		}
	} else {
		return fmt.Errorf(
			"index of server not valid: %v",
			serverIPIndex,
		)
	}

	// reconnect to allow the transaction to finish?
	database.Reconnect()

	// now find the server we just created and assert that it's valid
	server = &models.WgConfig{}
	serverResult := database.DB.First(&server, "is_server = ?", true)
	if serverResult.Error != nil {
		if serverResult.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf(
				"server was not correctly added to db: %v",
				serverResult.Error.Error(),
			)
		}

		return fmt.Errorf(
			"failed to look up server: %v",
			serverResult.Error.Error(),
		)
	}

	err = validateServer(*server)
	if err != nil {
		return err
	}

	// for j, ip := range allIPs {
	chunkSize := 250
	for j := 0; j < ips; j += chunkSize {
		end := j + chunkSize
		// Check if end is out of bounds
		if end > ips {
			end = ips
		}

		var wg sync.WaitGroup
		// for l := range allIPs[j:end] {
		for l := j; l < end; l++ {
			i := uint(l + 1) // the primary key in the db is 1-based index, not 0
			log.Printf("i=%v, l=%v, j=%v, end=%v", i, l, j, end)
			if allIPs[l].IsServerIP {
				continue
			}
			wg.Add(1)
			go prepFnWrapped(&wg, &progress, i, allIPs[l])
		}
		wg.Wait()
	}

	(*setProgressLabel)(
		fmt.Sprintf(
			"[%v/%v]: Generation step 2/2 complete",
			ips,
			ips,
		),
	)
	(*setProgressValue)(95)

	// 3. Finally, update the server config.
	(*setProgressLabel)("Saving final server config")
	(*setProgressValue)(99)
	server.Config = generateServerConfig(*server, conf, serverPeers.String(), network)
	saved := database.DB.Save(&server)
	if saved.Error != nil {
		return fmt.Errorf(
			"failed to save server config: %v",
			saved.Error.Error(),
		)
	}

	(*setProgressLabel)("Done")
	(*setProgressValue)(100)

	return nil
}
