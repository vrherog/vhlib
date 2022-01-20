package vhlib

const (
	DefaultTimeout = 1000

	DefaultPipeServer = "."
	DefaultPipeName   = "vhclient"

	DefaultInputStreamFile  = "/tmp/vhclient"
	DefaultOutputStreamFile = "/tmp/vhclient_response"

	CommandList = "LIST\n" // List devices

	CommandGetClientState = "GET CLIENT STATE\n" // Get the detailed full client state as an XML Document

	CommandUseDevice = "USE,%s" // Use a device:    "USE,<address>[,password]"

	CommandStopUsing = "STOP USING,%s" // Stop using a device:    "STOP USING,<address>"

	CommandStopUsingAll = "STOP USING ALL\n" // Stop using all devices on all clients

	CommandStopUsingAllLocal = "STOP USING ALL LOCAL\n" // Stop using all devices just for this client

	CommandDeviceInfo = "DEVICE INFO,%s" // Device Information:    "DEVICE INFO,<address>"

	CommandServerInfo = "SERVER INFO,%s" // Server Information:    "SERVER INFO,<server name>"

	CommandDeviceRename = "DEVICE RENAME,%s,%s" // Set device nickname:    "DEVICE RENAME,<address>,<nickname>"

	CommandServerRename = "SERVER RENAME,%s,%s" // Rename server:    "SERVER RENAME,<hubaddress:port>,<new name>"

	CommandAutoUseAll = "AUTO USE ALL\n" // Turn auto-use all devices on

	CommandAutoUseHub = "AUTO USE HUB,%s"
	// Turn Auto-use all devices on this hub on/off:    "AUTO USE HUB,<server name>"

	CommandAutoUsePort = "AUTO USE PORT,%s"
	// Turn Auto-use any device on this port on/off:    "AUTO USE PORT,<address>"

	CommandAutoUseDevice = "AUTO USE DEVICE,%s"
	// Turn Auto-use this device on any port on/off:    "AUTO USE DEVICE,<address>"

	CommandAutoUseDevicePort = "AUTO USE DEVICE PORT,%s"
	// Turn Auto-use this device on this port on/off:    "AUTO USE DEVICE PORT,<address>"

	CommandAutoUseClearAll = "AUTO USE CLEAR ALL\n" // Clear all auto-use settings

	CommandManualHubAdd = "MANUAL HUB ADD,%s"
	// Specify server to connect to:    "MANUAL HUB ADD,<address>[:port] | <EasyFind address>"

	CommandManualHubRemove = "MANUAL HUB REMOVE,%s"
	// Remove a manually specified hub:    "MANUAL HUB REMOVE,<address>[:port] | <EasyFind address>"

	CommandManualHubRemoveAll = "MANUAL HUB REMOVE ALL\n" // Remove all manually specified hubs

	CommandAddReverse = "ADD REVERSE,%s,%s"
	// Add a reverse client to the server:    "ADD REVERSE,<server serial>,<client address[:port]>"

	CommandRemoveReverse = "REMOVE REVERSE,%s,%s"
	// Remove a reverse client from the server:    "REMOVE REVERSE,<server serial>,<client address[:port]>"

	CommandListReverse = "LIST REVERSE,%s" // List all reverse clients:    "LIST REVERSE,<server serial>"

	CommandManualHubList = "MANUAL HUB LIST\n" // List manually specified hubs

	CommandListLicenses = "LIST LICENSES\n" // List licenses

	CommandLicenseServer = "LICENSE SERVER,%s" // License server:    "LICENSE SERVER,<license key>"

	CommandClearLog = "CLEAR LOG\n" // Clear client log

	CommandCustomEvent = "CUSTOM EVENT,%s,%s"
	// Set a custom device event:    "CUSTOM EVENT,<address>,<event>"

	CommandAutoFind = "AUTOFIND\n" // Turn auto-find on

	CommandReverse = "REVERSE" // Turn reverse lookup on

	CommandSslReverse = "SSLREVERSE" // Turn reverse SSL lookup on

	CommandExit = "EXIT" // Shutdown the client

	ResultOk = "OK"
)
