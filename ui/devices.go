package ui

import (
	"fmt"
	"log"
	"strconv"

	"wgnetui/constants"
	"wgnetui/database"
	"wgnetui/generator"
	"wgnetui/helpers"
	"wgnetui/models"
	"wgnetui/queries"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var SelectedDevice models.WgConfig

// Seems like TypedShortcuts still need some love. Accepting the fact that
// a focused entry will override window focus events.
//
// https://github.com/fyne-io/fyne/issues/3038
//

// func NewEntryWithShortcuts() *entryWithShortcuts {
// 	e := &entryWithShortcuts{
// 		widget.NewEntry(),
// 	}
// 	e.ExtendBaseWidget(e)
// 	e.Enable()
// 	return e
// }

// type entryWithShortcuts struct {
// 	*widget.Entry
// }

// // Used for enabling Ctrl+S support when an entry is focused, as well as
// // anything else.
// // https://developer.fyne.io/explore/shortcuts#adding-shortcuts-to-an-entry
// func (m *entryWithShortcuts) TypedShortcut(s fyne.Shortcut) {
// 	if _, ok := s.(*desktop.CustomShortcut); !ok {
// 		m.Entry.TypedShortcut(s)
// 		return
// 	}

// 	log.Println("Shortcut typed:", s)
// }

// refreshDevices accepts a pointer to a list of devices so that it can update
// the data source for the Fyne list. Make sure to call list.Refresh() for the
// Fyne list itself afterwards, if appropriate, or just use refreshDeviceList
// to do both in one function call.
func refreshDevices(devices *[]models.WgConfig) error {
	result := database.DB.Find(&devices)
	if result.Error != nil {
		return fmt.Errorf(
			"failed to get devices view: %v",
			result.Error.Error(),
		)
	}

	return nil
}

// refreshDevicesList will execute refreshDevices against the provided devices slice,
// and will subsequently call list.Refresh. This is useful for bundling
// the calls to update the list of devices and refreshing the list.
// If either the devices or list pointers are nil, it will not do anything
// aside from leaving a log message.
func refreshDevicesList(devices *[]models.WgConfig, list *widget.List) error {
	if list == nil || devices == nil {
		log.Println("warning: refreshDeviceList has one or more nil pointers")
		return nil
	}

	err := refreshDevices(devices)
	if err != nil {
		return err
	}

	list.Refresh()

	return nil
}

// GetDevicesView returns a Fyne container that holds a view for assigning
// devices and names, as a list.
func GetDevicesView() (*container.Split, error) {
	devices := []models.WgConfig{}

	err := refreshDevices(&devices)
	if err != nil {
		return nil, err
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("")
	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Device description")
	extraEntry := widget.NewMultiLineEntry()
	extraEntry.SetPlaceHolder("# extra wireguard [Interface] section config lines can go here\nPostUp = foo")
	configEntry := widget.NewMultiLineEntry()
	configEntry.SetPlaceHolder("[Interface]\n# ...")
	configEntry.SetMinRowsVisible(10)
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("192.168.100.1")
	ipEntry.Disable()
	allowedIPsEntry := widget.NewEntry()
	allowedIPsEntry.SetPlaceHolder("0.0.0.0/0")
	persistentKeepAliveEntry := widget.NewEntry()
	persistentKeepAliveEntry.SetPlaceHolder("25")
	mtuEntry := widget.NewEntry()
	mtuEntry.SetPlaceHolder("1280")
	endpointEntry := widget.NewEntry()
	endpointEntry.SetPlaceHolder("1.2.3.4")
	endpointPortEntry := widget.NewEntry()
	endpointPortEntry.SetPlaceHolder("51820")
	dnsEntry := widget.NewEntry()
	dnsEntry.SetPlaceHolder("192.168.100.1")
	privateKeyEntry := widget.NewPasswordEntry()
	privateKeyEntry.SetPlaceHolder("")
	publicKeyEntry := widget.NewEntry()
	publicKeyEntry.SetPlaceHolder("")
	preSharedKeyEntry := widget.NewPasswordEntry()
	preSharedKeyEntry.SetPlaceHolder("")

	initialQR, err := helpers.GetQR("Please choose a device first.")
	if err != nil {
		err = fmt.Errorf("Failed to initialize QR render: %v", err.Error())
		dialog.ShowError(err, *W)
		return nil, err
	}

	deviceQR := canvas.NewImageFromImage(initialQR)
	deviceQR.FillMode = canvas.ImageFillOriginal

	exportConfigToFile := widget.NewButtonWithIcon(
		constants.ExportConfigToFileButtonLabel,
		theme.DocumentSaveIcon(),
		func() {
			t := "peer"
			if SelectedDevice.IsServer {
				t = "server"
			}
			SaveToFileDialog(
				SelectedDevice.Config,
				fmt.Sprintf(
					"wgnetui_%v_%v.conf",
					SelectedDevice.IP,
					t,
				),
				[]string{".conf"},
			)
		},
	)
	// exportConfigToFile.Importance = widget.HighImportance
	showQRCodeButton := widget.NewButtonWithIcon(
		constants.ToggleQRCodeButtonLabel,
		theme.VisibilityIcon(),
		func() {
			if deviceQR == nil {
				return
			}

			if deviceQR.Visible() {
				deviceQR.Hide()
			} else {
				deviceQR.Show()
			}
		},
	)

	var list *widget.List

	selectDeviceByID := func(id int) {
		// log.Printf("selected %v", id)
		if database.DB == nil {
			dialog.ShowError(
				fmt.Errorf(constants.ErrorMessageNoDB),
				*W,
			)
			return
		}
		newSelectedDevice := models.WgConfig{}
		result := database.DB.First(&newSelectedDevice, id+1)
		if result.Error != nil {
			dialog.ShowError(
				fmt.Errorf(constants.ErrorMessageNoDevice),
				*W,
			)
			return
		}

		SelectedDevice = newSelectedDevice

		// update the form to show the selected device
		nameEntry.SetText(SelectedDevice.Name)
		descriptionEntry.SetText(SelectedDevice.Description)
		extraEntry.SetText(SelectedDevice.Extra)
		ipEntry.SetText(SelectedDevice.IP)
		allowedIPsEntry.SetText(SelectedDevice.AllowedIPs)
		persistentKeepAliveEntry.SetText(fmt.Sprintf("%v", SelectedDevice.PersistentKeepAlive))
		mtuEntry.SetText(fmt.Sprintf("%v", SelectedDevice.MTU))
		endpointEntry.SetText(SelectedDevice.Endpoint)
		endpointPortEntry.SetText(fmt.Sprintf("%v", SelectedDevice.EndpointPort))
		dnsEntry.SetText(SelectedDevice.DNS)
		privateKeyEntry.SetText(SelectedDevice.PrivateKey)
		publicKeyEntry.SetText(SelectedDevice.PublicKey)
		preSharedKeyEntry.SetText(SelectedDevice.PreSharedKey)

		if SelectedDevice.IsServer {
			// descriptionEntry.Disable()
			// nameEntry.Disable()
			configEntry.SetText(constants.MessageServerConfigNotShown)
			configEntry.Disable()
			// extraEntry.Disable()
			allowedIPsEntry.Disable()
			persistentKeepAliveEntry.Disable()
			// mtuEntry.Disable()
			endpointEntry.Disable()
			endpointPortEntry.Disable()
			// dnsEntry.Disable()
			// privateKeyEntry.Disable()
			// publicKeyEntry.Disable()
			preSharedKeyEntry.Disable()

			// hide the qr for servers and replace it with a button to
			// export to file
			deviceQR.Hide()
			showQRCodeButton.Disable()
		} else {
			// nameEntry.Enable()
			configEntry.SetText(SelectedDevice.Config)
			configEntry.Enable()
			// descriptionEntry.Enable()
			// extraEntry.Enable()
			allowedIPsEntry.Enable()
			persistentKeepAliveEntry.Enable()
			// mtuEntry.Enable()
			endpointEntry.Enable()
			endpointPortEntry.Enable()
			// dnsEntry.Enable()
			// privateKeyEntry.Enable()
			// publicKeyEntry.Enable()
			preSharedKeyEntry.Enable()
			showQRCodeButton.Enable()

			// update the qr only for non-servers
			qrc, err := helpers.GetQR(SelectedDevice.Config)
			if err != nil {
				err = fmt.Errorf("Failed to get QR code: %v", err.Error())
				dialog.ShowError(err, *W)
			}
			deviceQR.Image = qrc
			deviceQR.FillMode = canvas.ImageFillContain
			deviceQR.Hide()
			deviceQR.Refresh()
		}
	}

	saveForm := func() {
		if database.DB == nil {
			dialog.ShowError(
				fmt.Errorf(constants.ErrorMessageNoDB),
				*W,
			)
			return
		}

		parsedPersistentKeepAlive, err := strconv.ParseUint(persistentKeepAliveEntry.Text, 10, 64)
		if err != nil {
			dialog.ShowError(
				fmt.Errorf("The provided PersistentKeepAlive value was not valid: %v", err.Error()),
				*W,
			)
			return
		}

		parsedEndpointPortEntry, err := strconv.ParseUint(endpointPortEntry.Text, 10, 16)
		if err != nil {
			dialog.ShowError(
				fmt.Errorf("The provided EndpointPort value was not valid: %v", err.Error()),
				*W,
			)
			return
		}

		parsedMTUEntry, err := strconv.ParseUint(mtuEntry.Text, 10, 16)
		if err != nil {
			dialog.ShowError(
				fmt.Errorf("The provided MTU value was not valid: %v", err.Error()),
				*W,
			)
			return
		}

		server, err := queries.GetServer()
		if err != nil {
			dialog.ShowError(
				fmt.Errorf("Failed to retrieve server when saving device: %v", err.Error()),
				*W,
			)
			return
		}

		// update the device in the DB
		SelectedDevice.Name = nameEntry.Text
		// cannot edit a server config directly, it's disabled
		if !SelectedDevice.IsServer {
			SelectedDevice.Config = configEntry.Text
		}
		SelectedDevice.Name = nameEntry.Text
		SelectedDevice.Description = descriptionEntry.Text
		SelectedDevice.Extra = extraEntry.Text
		SelectedDevice.IP = ipEntry.Text
		SelectedDevice.AllowedIPs = allowedIPsEntry.Text
		SelectedDevice.PersistentKeepAlive = uint(parsedPersistentKeepAlive)
		SelectedDevice.MTU = uint16(parsedMTUEntry)
		SelectedDevice.Endpoint = endpointEntry.Text
		SelectedDevice.EndpointPort = uint16(parsedEndpointPortEntry)
		SelectedDevice.DNS = dnsEntry.Text
		SelectedDevice.PrivateKey = privateKeyEntry.Text
		SelectedDevice.PublicKey = publicKeyEntry.Text
		SelectedDevice.PreSharedKey = preSharedKeyEntry.Text

		saveDevice := func() {
			result := database.DB.Save(SelectedDevice)
			if result.Error != nil {
				dialog.ShowError(
					fmt.Errorf("Failed to save this device: %v", result.Error.Error()),
					*W,
				)
				return
			}

			if list != nil {
				err = refreshDevicesList(&devices, list)
				if err != nil {
					dialog.ShowError(
						fmt.Errorf(
							"Failed to refresh devices list: %v", err.Error(),
						),
						*W,
					)
					return
				}
				selectDeviceByID(int(SelectedDevice.ID - 1))
			}

			dialog.ShowInformation(
				"Saved Device",
				fmt.Sprintf("Device %v was saved successfully.", SelectedDevice.IP),
				*W,
			)
		}

		// we don't want to auto-generate the server's config, so just save it
		// without a prompt
		if SelectedDevice.IsServer {
			saveDevice()
			return
		}

		// ask if the user wants to regenerate the config or keep their current
		// config
		dialog.ShowConfirm(
			"Update config?",
			"If you've updated any values, the config for this device may need to be regenerated. This may wipe out any customizations you've made to it.",
			func(confirmed bool) {
				if confirmed {
					newConfig, err := generator.GenerateConfig(
						SelectedDevice,
						server,
					)
					if err != nil {
						dialog.ShowError(
							fmt.Errorf("Failed to generate config for this device: %v", err.Error()),
							*W,
						)
						return
					}

					SelectedDevice.Config = newConfig
				}

				saveDevice()
			},
			*W,
		)
	}

	deviceEditorForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry, HintText: constants.HelpTextSelectedDeviceName},
			{Text: "Config", Widget: configEntry, HintText: ""},
			{Text: "Description", Widget: descriptionEntry, HintText: ""},
			{Text: "Extra", Widget: extraEntry, HintText: ""},
			{Text: "IP", Widget: ipEntry, HintText: ""},
			{Text: "AllowedIPs", Widget: allowedIPsEntry, HintText: ""},
			{Text: "PersistentKeepAlive", Widget: persistentKeepAliveEntry, HintText: ""},
			{Text: "MTU", Widget: mtuEntry, HintText: ""},
			{Text: "Endpoint", Widget: endpointEntry, HintText: ""},
			{Text: "EndpointPort", Widget: endpointPortEntry, HintText: ""},
			{Text: "DNS", Widget: dnsEntry, HintText: ""},
			{Text: "PrivateKey", Widget: privateKeyEntry, HintText: ""},
			{Text: "PublicKey", Widget: publicKeyEntry, HintText: ""},
			{Text: "PreSharedKey", Widget: preSharedKeyEntry, HintText: ""},
		},
		SubmitText: "Save",
		OnSubmit:   saveForm,
	}

	listUpdateFn := func(i widget.ListItemID, o fyne.CanvasObject) {
		if devices[i].IsServer {
			o.(*widget.Label).SetText(
				fmt.Sprintf("[%v] %v [server]", devices[i].IP, devices[i].Name),
			)
		} else {
			o.(*widget.Label).SetText(
				fmt.Sprintf("[%v] %v ", devices[i].IP, devices[i].Name),
			)
		}
	}

	list = widget.NewList(
		func() int {
			return len(devices)
		},
		func() fyne.CanvasObject {
			b := widget.NewLabel("")
			b.Alignment = fyne.TextAlignLeading
			return b
		},
		listUpdateFn,
	)

	list.OnUnselected = func(id int) {
		// log.Printf("unselected %v", id)
	}
	list.OnSelected = selectDeviceByID

	selectedDeviceView := container.NewVScroll(
		container.NewVBox(
			showQRCodeButton,
			deviceQR,
			exportConfigToFile,
			container.NewBorder(
				GetPadding(0, 10),
				GetPadding(0, 10),
				GetPadding(10, 0),
				GetPadding(10, 0),
				deviceEditorForm,
			),
		),
	)

	c := container.NewHSplit(
		selectedDeviceView,
		container.NewScroll(list),
	)

	// add a keyboard shortcut Ctrl+R to allow reloading of the device list
	devicesCtrlRShortcut := func() {
		// if ActiveTab != constants.TabDevices {
		// 	return
		// }

		err = refreshDevicesList(&devices, list)
		if err != nil {
			dialog.ShowError(
				fmt.Errorf(
					"Failed to refresh devices list: %v", err.Error(),
				),
				*W,
			)
			return
		}

		// if LoadDevicesView != nil {
		// 	(*LoadDevicesView)()
		// 	dialog.ShowInformation("Refreshed", "Refreshed devices successfully.", *W)
		// }
		dialog.ShowInformation("Refreshed", "Refreshed devices successfully.", *W)
	}
	CtrlRShortcuts[constants.TabDevices] = &devicesCtrlRShortcut

	// add a keyboard shortcut Ctrl+S to allow saving of the device list
	devicesCtrlSShortcut := func() {
		if ActiveTab != constants.TabDevices {
			return
		}

		if SelectedDevice.ID == 0 {
			return
		}

		saveForm()
	}
	CtrlSShortcuts[constants.TabDevices] = &devicesCtrlSShortcut

	// select the first device on load
	if len(devices) > 0 {
		selectDeviceByID(0)
	}

	return c, nil
}
