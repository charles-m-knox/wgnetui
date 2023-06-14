package ui

import (
	"fmt"
	"log"
	"strconv"

	"wgnetui/constants"
	"wgnetui/database"
	"wgnetui/helpers"
	"wgnetui/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var SelectedDevice models.WgConfig

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

// GetDevicesView returns a Fyne container that holds a view for assigning
// devices and names, as a list.
func GetDevicesView(w fyne.Window) (*container.Split, error) {
	devices := []models.WgConfig{}
	result := database.DB.Find(&devices)
	if result.Error != nil {
		return nil, fmt.Errorf(
			"failed to get devices view: %v",
			result.Error.Error(),
		)
	}

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("device-01")
	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Device description")
	extraEntry := widget.NewMultiLineEntry()
	extraEntry.SetPlaceHolder("# extra wireguard [Interface] section config lines can go here\nPostUp = foo")
	configEntry := widget.NewMultiLineEntry()
	configEntry.SetPlaceHolder("[Interface]\n# ...")
	configEntry.SetMinRowsVisible(13)
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
		dialog.ShowError(err, w)
		return nil, err
	}

	deviceQR := canvas.NewImageFromImage(initialQR)
	deviceQR.FillMode = canvas.ImageFillOriginal

	var list *widget.List

	deviceEditorForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry, HintText: constants.HelpTextSelectedDeviceName},
			{Text: "Config", Widget: configEntry, HintText: ""},
			{Text: "Description", Widget: descriptionEntry, HintText: ""},
			{Text: "Extra", Widget: extraEntry, HintText: ""},
			{Text: "IP", Widget: ipEntry, HintText: ""},
			{Text: "AllowedIPs", Widget: allowedIPsEntry, HintText: ""},
			{Text: "PersistentKeepAlive", Widget: persistentKeepAliveEntry, HintText: ""},
			{Text: "Mtu", Widget: mtuEntry, HintText: ""},
			{Text: "Endpoint", Widget: endpointEntry, HintText: ""},
			{Text: "EndpointPort", Widget: endpointPortEntry, HintText: ""},
			{Text: "Dns", Widget: dnsEntry, HintText: ""},
			{Text: "PrivateKey", Widget: privateKeyEntry, HintText: ""},
			{Text: "PublicKey", Widget: publicKeyEntry, HintText: ""},
			{Text: "PreSharedKey", Widget: preSharedKeyEntry, HintText: ""},
		},
		SubmitText: "Save",
		OnSubmit: func() {
			if database.DB == nil {
				dialog.ShowError(
					fmt.Errorf(constants.ErrorMessageNoDB),
					w,
				)
				return
			}

			parsedPersistentKeepAlive, err := strconv.ParseUint(persistentKeepAliveEntry.Text, 10, 64)
			if err != nil {
				dialog.ShowError(
					fmt.Errorf("The provided PersistentKeepAlive value was not valid: %v", err.Error()),
					w,
				)
				return
			}

			parsedEndpointPortEntry, err := strconv.ParseUint(endpointPortEntry.Text, 10, 16)
			if err != nil {
				dialog.ShowError(
					fmt.Errorf("The provided EndpointPort value was not valid: %v", err.Error()),
					w,
				)
				return
			}

			parsedMTUEntry, err := strconv.ParseUint(mtuEntry.Text, 10, 16)
			if err != nil {
				dialog.ShowError(
					fmt.Errorf("The provided MTU value was not valid: %v", err.Error()),
					w,
				)
				return
			}

			// update the device in the DB
			SelectedDevice.Name = nameEntry.Text
			SelectedDevice.Config = configEntry.Text
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
			result := database.DB.Save(SelectedDevice)
			if result.Error != nil {
				dialog.ShowError(
					fmt.Errorf("Failed to save this device: %v", result.Error.Error()),
					w,
				)
				return
			}

			if list != nil {
				err = refreshDevices(&devices)
				if err != nil {
					dialog.ShowError(
						fmt.Errorf(
							"Failed to refresh devices list: %v", err.Error(),
						),
						w,
					)
					return
				}
				list.Refresh()
			}

			dialog.ShowInformation(
				"Generated",
				fmt.Sprintf("Device %v was saved successfully.", SelectedDevice.IP),
				w,
			)
		},
	}

	listUpdateFn := func(i widget.ListItemID, o fyne.CanvasObject) {
		if devices[i].IsServer {
			o.(*widget.Label).SetText(
				fmt.Sprintf("%v. %v [%v] [server]", i+1, devices[i].Name, devices[i].IP),
			)
			// o.(*widget.Label).SetIcon(theme.ComputerIcon())
		} else {
			o.(*widget.Label).SetText(
				fmt.Sprintf("%v. %v [%v]", i+1, devices[i].Name, devices[i].IP),
			)
			// o.(*widget.Label).Importance = widget.LowImportance
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
		log.Printf("unselected %v", id)
	}
	list.OnSelected = func(id int) {
		log.Printf("selected %v", id)
		if database.DB == nil {
			dialog.ShowError(
				fmt.Errorf(constants.ErrorMessageNoDB),
				w,
			)
			return
		}
		newSelectedDevice := models.WgConfig{}
		result := database.DB.First(&newSelectedDevice, id+1)
		if result.Error != nil {
			dialog.ShowError(
				fmt.Errorf(constants.ErrorMessageNoDevice),
				w,
			)
			return
		}

		SelectedDevice = newSelectedDevice

		// update the form to show the selected device
		nameEntry.SetText(SelectedDevice.Name)
		configEntry.SetText(SelectedDevice.Config)
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

		qrc, err := helpers.GetQR(SelectedDevice.Config)
		if err != nil {
			err = fmt.Errorf("Failed to get QR code: %v", err.Error())
			dialog.ShowError(err, w)
		}
		deviceQR.Image = qrc
		deviceQR.FillMode = canvas.ImageFillContain
		deviceQR.Refresh()
	}

	rightSide := container.NewVScroll(container.NewVBox(
		deviceQR,
		deviceEditorForm,
	))
	c := container.NewHSplit(
		container.NewScroll(list),
		rightSide,
	)

	return c, nil
}
