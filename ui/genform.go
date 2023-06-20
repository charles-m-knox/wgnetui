package ui

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"wgnetui/constants"
	"wgnetui/database"
	"wgnetui/generator"
	"wgnetui/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func GetWgGenForm() *fyne.Container {
	gf := models.GenerationForm{}
	gotPreconfigured := false
	if database.DB != nil {
		result := database.DB.Order("created_at desc").Limit(1).First(&gf)
		if result.Error != nil {
			dialog.ShowError(
				fmt.Errorf(
					"Failed to retrieve previously stored configuration: %v",
					result.Error.Error(),
				),
				*W,
			)
		}
		gotPreconfigured = true
	}

	log.Println("opening new form", database.OpenedFilePath)
	cidr := widget.NewEntry()
	cidr.SetPlaceHolder("192.168.100.0/24")
	mtu := widget.NewEntry()
	mtu.SetPlaceHolder("1280")
	allowedIPs := widget.NewEntry()
	allowedIPs.SetPlaceHolder("0.0.0.0/0")
	persistentKeepAlive := widget.NewEntry()
	persistentKeepAlive.SetPlaceHolder("")
	server := widget.NewEntry()
	server.SetPlaceHolder("192.168.100.2")
	endpoint := widget.NewEntry()
	endpoint.SetPlaceHolder("")
	endpointPort := widget.NewEntry()
	endpointPort.SetPlaceHolder("")
	dns := widget.NewEntry()
	dns.SetPlaceHolder("192.168.100.2")

	name := widget.NewEntry()
	name.SetPlaceHolder("my-device-")
	description := widget.NewEntry()
	description.SetPlaceHolder("Whatever you want")
	extra := widget.NewMultiLineEntry()
	extra.SetPlaceHolder("PostUp = foo")
	serverInterface := widget.NewEntry()
	serverInterface.SetPlaceHolder("eth0")

	regenerateKeys := widget.NewCheck("Regenerate all Wireguard keys", nil)
	resetAll := widget.NewCheck("Delete pre-existing data", nil)
	forceAllowedIPs := widget.NewCheck("Ignore old AllowedIPs values", nil)
	forcePersistentKeepAlive := widget.NewCheck("Ignore old PersistentKeepAlive values", nil)
	forceMTU := widget.NewCheck("Ignore old MTU values", nil)
	forceEndpoint := widget.NewCheck("Ignore old Endpoint values", nil)
	forceEndpointPort := widget.NewCheck("Ignore old EndpointPort values", nil)
	forceDNS := widget.NewCheck("Ignore old DNS values", nil)
	forceName := widget.NewCheck("Ignore old Name values", nil)
	forceDescription := widget.NewCheck("Ignore old Description values", nil)
	forceExtra := widget.NewCheck("Ignore old Extra values", nil)

	// leaving here for reference - please delete later
	// pskGenBtn := widget.NewButtonWithIcon(
	// 	"Generate PSK",
	// 	theme.MediaReplayIcon(),
	// 	func() {
	// 		psk.SetText(helpers.GeneratePreSharedKey())
	// 	},
	// )

	if gotPreconfigured {
		log.Printf("got preconfigured: %v", gotPreconfigured)

		cidr.SetText(gf.CIDR)
		mtu.SetText(fmt.Sprintf("%v", gf.MTU))
		allowedIPs.SetText(gf.AllowedIPs)
		persistentKeepAlive.SetText(fmt.Sprintf("%v", gf.PersistentKeepAlive))
		server.SetText(gf.Server)
		endpoint.SetText(gf.Endpoint)
		endpointPort.SetText(fmt.Sprintf("%v", gf.EndpointPort))
		dns.SetText(gf.DNS)
		name.SetText(gf.Name)
		description.SetText(gf.Description)
		extra.SetText(gf.Extra)
		serverInterface.SetText(gf.ServerInterface)
		regenerateKeys.SetChecked(gf.RegenerateKeys)
		resetAll.SetChecked(gf.ResetAll)
		forceAllowedIPs.SetChecked(gf.ForceAllowedIPs)
		forcePersistentKeepAlive.SetChecked(gf.ForcePersistentKeepAlive)
		forceMTU.SetChecked(gf.ForceMTU)
		forceEndpoint.SetChecked(gf.ForceEndpoint)
		forceEndpointPort.SetChecked(gf.ForceEndpointPort)
		forceDNS.SetChecked(gf.ForceDNS)
		forceName.SetChecked(gf.ForceName)
		forceDescription.SetChecked(gf.ForceDescription)
		forceExtra.SetChecked(gf.ForceExtra)
	}

	log.Println("new form:", gf)

	// create the form
	genForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "CIDR", Widget: cidr, HintText: constants.HelpTextCIDR},
			{Text: "DNS", Widget: dns, HintText: constants.HelpTextDNS},
			{Text: "Server", Widget: server, HintText: constants.HelpTextServer},
			{Text: "Server Interface", Widget: serverInterface, HintText: ""},
			{Text: "Endpoint", Widget: endpoint, HintText: constants.HelpTextEndpoint},
			{Text: "Endpoint Port", Widget: endpointPort, HintText: constants.HelpTextPort},
			{Text: "MTU", Widget: mtu, HintText: constants.HelpTextMTU},
			{Text: "AllowedIPs", Widget: allowedIPs, HintText: ""},
			{Text: "PersistentKeepAlive", Widget: persistentKeepAlive, HintText: ""},
			{Text: "Name", Widget: name, HintText: "${id} gets replaced with the numeric ID"},
			{Text: "Description", Widget: description, HintText: "${id} gets replaced with the numeric ID. ${name} gets replaced with the device's name."},
			{Text: "Extra", Widget: extra, HintText: "${id} gets replaced with the numeric ID. ${name} gets replaced with the device's name."},
			{Text: "Regenerate Keys", Widget: regenerateKeys, HintText: constants.HelpTextRegenerateKeys},
			{Text: "Reset All", Widget: resetAll, HintText: ""},
			{Text: "Force AllowedIPs", Widget: forceAllowedIPs, HintText: ""},
			{Text: "Force PersistentKeepAlive", Widget: forcePersistentKeepAlive, HintText: ""},
			{Text: "Force MTU", Widget: forceMTU, HintText: ""},
			{Text: "Force Endpoint", Widget: forceEndpoint, HintText: ""},
			{Text: "Force EndpointPort", Widget: forceEndpointPort, HintText: ""},
			{Text: "Force DNS", Widget: forceDNS, HintText: ""},
			{Text: "Force Name", Widget: forceName, HintText: ""},
			{Text: "Force Description", Widget: forceDescription, HintText: ""},
			{Text: "Force Extra", Widget: forceExtra, HintText: ""},
		},
		SubmitText: "Update",
		OnSubmit: func() {
			if ProgressBarDialog == nil {
				dialog.ShowError(
					fmt.Errorf("Unable to show progress dialog."),
					*W,
				)
				return
			}

			ProgressBarDialog.Fresh(
				"Generating...",
				"Hide",
				"Checking database...",
			)

			ProgressBarDialog.SetProgress(0)
			ProgressBarDialog.Show()
			defer ProgressBarDialog.Hide()

			database.Reconnect() // reading the db between writes can cause failures
			if database.DB == nil {
				dialog.ShowError(
					fmt.Errorf(constants.ErrorMessageNoDB),
					*W,
				)
				return
			}

			ProgressBarDialog.SetLabel("Validating inputs...")
			ProgressBarDialog.SetProgress(1)

			g := models.GenerationForm{}

			valstr := []string{}

			firstIP, cidrNet, err := net.ParseCIDR(cidr.Text)
			if err != nil {
				valstr = append(valstr, fmt.Sprintf("CIDR: %v", err.Error()))
			}

			if cidrNet == nil {
				valstr = append(valstr, "CIDR: Was unable to determine the CIDR subnet")
			}

			// assert that the user provided an aligned CIDR block
			actualCIDR := fmt.Sprintf("%v", cidrNet)
			if actualCIDR != cidr.Text {
				valstr = append(valstr,
					fmt.Sprintf(
						"CIDR: The range you provided (%v) is not a correctly aligned subnet. You must adjust to use the correctly aligned subnet %v instead.",
						cidr.Text,
						actualCIDR,
					),
				)
			}

			parsedMTU64, err := strconv.ParseUint(mtu.Text, 10, 16)
			if err != nil {
				valstr = append(valstr, fmt.Sprintf("MTU: %v", err.Error()))
			}

			parsedEndpointPort64, err := strconv.ParseUint(endpointPort.Text, 10, 16)
			if err != nil {
				valstr = append(valstr, fmt.Sprintf("Port: %v", err.Error()))
			}

			parsedServer := net.ParseIP(server.Text)
			if parsedServer == nil {
				valstr = append(valstr, fmt.Sprintf("Server: Is not an IP address: %v", server.Text))
			}

			if cidrNet != nil && !cidrNet.Contains(parsedServer) {
				valstr = append(valstr, fmt.Sprintf("Server: Must be an IP address within the range %v: %v", cidr.Text, server.Text))
			}

			if endpoint.Text == "" {
				valstr = append(valstr, "Endpoint: Must not be empty, and must be a hostname or IPv4 or IPv6 address")
			}

			parsedDNS := net.ParseIP(dns.Text)
			if parsedDNS == nil {
				if dns.Text != "" {
					valstr = append(valstr, fmt.Sprintf("DNS: Is not an IP address: %v", dns.Text))
				}
			}

			parsedPersistentKeepAlive64, err := strconv.ParseUint(persistentKeepAlive.Text, 10, 64)
			if err != nil {
				valstr = append(valstr, fmt.Sprintf("Port: %v", err.Error()))
			}

			if len(valstr) > 0 {
				ProgressBarDialog.Hide()
				dialog.ShowError(
					fmt.Errorf("One or more values were entered incorrectly.\n\n%v", strings.Join(valstr, "\n")),
					*W,
				)
				return
			}

			ProgressBarDialog.SetLabel("Preparing inputs...")
			ProgressBarDialog.SetProgress(2)

			g.CIDR = cidr.Text // warning: rendering parsedCIDR as a string results in the /24 suffix dropping
			g.DNS = fmt.Sprintf("%v", parsedDNS)
			g.Server = fmt.Sprintf("%v", parsedServer)
			g.ServerInterface = serverInterface.Text
			g.Endpoint = endpoint.Text
			g.EndpointPort = uint16(parsedEndpointPort64)
			g.MTU = uint16(parsedMTU64)
			g.AllowedIPs = allowedIPs.Text
			g.PersistentKeepAlive = uint(parsedPersistentKeepAlive64)
			g.Name = name.Text
			g.Description = description.Text
			g.Extra = extra.Text

			g.RegenerateKeys = regenerateKeys.Checked
			g.ResetAll = resetAll.Checked
			g.ForceAllowedIPs = forceAllowedIPs.Checked
			g.ForcePersistentKeepAlive = forcePersistentKeepAlive.Checked
			g.ForceMTU = forceMTU.Checked
			g.ForceEndpoint = forceEndpoint.Checked
			g.ForceEndpointPort = forceEndpointPort.Checked
			g.ForceDNS = forceDNS.Checked
			g.ForceName = forceName.Checked
			g.ForceDescription = forceDescription.Checked
			g.ForceExtra = forceExtra.Checked

			ProgressBarDialog.SetLabel("Validating inputs...")
			ProgressBarDialog.SetProgress(3)

			result := database.DB.Create(&g)
			if result.Error != nil {
				dialog.ShowError(
					fmt.Errorf("Failed to store this configuration: %v", result.Error.Error()),
					*W,
				)
				return
			}

			ProgressBarDialog.SetLabel("Starting generation...")
			ProgressBarDialog.SetProgress(4)

			// define these wrapper functions to work around Go's quirkiness
			// of being unable to pass &ProgressBarDialog.SetLabel / SetProgress
			// directly to the function generator.Generate()
			sl := func(l string) {
				ProgressBarDialog.SetLabel(l)
			}
			sp := func(v float64) {
				ProgressBarDialog.SetProgress(v)
			}
			// now generate all keys / etc based on the config
			err = generator.Generate(
				g,
				parsedServer,
				firstIP,
				cidrNet,
				&sl,
				&sp,
			)
			if err != nil {
				dialog.ShowError(
					fmt.Errorf(
						"Failed to generate configuration: %v",
						err.Error(),
					),
					*W,
				)
				return
			}

			dialog.ShowInformation(
				"Generated",
				"The configuration was saved successfully.",
				*W,
			)
		},
	}

	return container.NewMax(container.NewScroll(container.NewBorder(
		GetPadding(0, 10),
		GetPadding(0, 10),
		GetPadding(10, 0),
		GetPadding(10, 0),
		container.NewVBox(
			genForm,
		),
	)))
}
