package constants

const (
	PlaceholderMarkdown = "# Placeholder markdown\n\nThis is some placeholder text."

	ErrorMessageNoDB           = "The config database is not accessible. Please restart this application. If this continues, please verify that the database file exists and is readable."
	ErrorMessageNoDevice       = "The selected device could not be found in the database."
	ErrorMessageFailGenKeyPair = "Failed to generate a keypair."

	HelpTextSelectedDeviceName = "Name of this device. Can by any string."

	HelpTextCIDR           = "IPv4 CIDR range for the mesh network, such as 192.168.1.0/24"           // This application will generate one peer config for each IP address within this range."
	HelpTextMTU            = "The Maximum Transmission Unit for every connection (suggest 1280-1500)" // Consider starting with 1280 and incrementing until you see failures."
	HelpTextPort           = "The port that the server will listen on, typically 51820"
	HelpTextServer         = "IPv4 address of the Wireguard server, within the provided CIDR range" // A peer connection added to this server's Wireguard config for every IP address within the specified CIDR range."
	HelpTextEndpoint       = "The hostname/IP of the server that each peer will connect to"         // , such as example.com or 5.5.5.5."
	HelpTextRegenerateKeys = "Instead of reusing keys between config changes, generate new ones"    // "Wireguard keys that are generated by this application will be reused, even when the CIDR value changes. Check this box to forcefully regenerate new keys instead of reusing old ones."
	HelpTextDNS            = "DNS server peers will use, such as 1.1.1.1"
	HelpTextPSK            = "Wireguard Pre-Shared Key (PSK) value ('wg genpsk' on command line)"

	DefaultAllowedIPs          = "0.0.0.0/0"
	DefaultPersistentKeepAlive = uint(25)
	DefaultEndpointPort        = uint16(51820)
	DefaultMTU                 = uint16(1280)

	TabAbout     = "About"
	TabGenerator = "Generator"
	TabDevices   = "Devices"
)
