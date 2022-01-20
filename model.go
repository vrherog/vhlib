package vhlib

import (
	"fmt"
	"strings"
)

type Connection struct {
	ConnectionId       string `xml:"connectionId,attr"`
	Secure             bool   `xml:"secure,attr"`
	ServerMajor        int    `xml:"serverMajor,attr"`
	ServerMinor        int    `xml:"serverMinor,attr"`
	ServerRevision     int    `xml:"serverRevision,attr"`
	RemoteAdmin        bool   `xml:"remoteAdmin,attr"`
	ServerName         string `xml:"serverName,attr"`
	InterfaceName      string `xml:"interfaceName,attr"`
	Hostname           string `xml:"hostname,attr"`
	ServerSerial       string `xml:"serverSerial,attr"`
	LicenseMaxDevices  int    `xml:"license_max_devices"`
	State              int    `xml:"state,attr"`
	ConnectedTime      string `xml:"connectedTime,attr"`
	Host               string `xml:"host,attr"`
	Port               uint16 `xml:"port,attr"`
	HasError           bool   `xml:"error,attr"`
	UUID               string `xml:"uuid,attr"`
	TransportId        string `xml:"transportId,attr"`
	EasyFindEnabled    bool   `xml:"easyFindEnabled,attr"`
	EasyFindAvailable  bool   `xml:"easyFindAvailable,attr"`
	EasyFindId         string `xml:"easyFindId,attr"`
	EasyFindPin        string `xml:"easyFindPin,attr"`
	EasyFindAuthorized int    `xml:"easyFindAuthorized,attr"`
	IP                 string `xml:"ip,attr"`
}

type Device struct {
	Vendor                            string `xml:"vendor,attr"`
	VendorId                          int    `xml:"idVendor,attr"`
	Product                           string `xml:"product,attr"`
	ProductId                         int    `xml:"idProduct,attr"`
	Address                           int    `xml:"address,attr"`
	ConnectionId                      int    `xml:"connectionId,attr"`
	State                             int    `xml:"state,attr"`
	ServerSerial                      string `xml:"serverSerial,attr"`
	ServerName                        string `xml:"serverName,attr"`
	ServerInterfaceName               string `xml:"serverInterfaceName,attr"`
	DeviceSerial                      string `xml:"deviceSerial,attr"`
	ConnectionUUID                    string `xml:"connectionUUID,attr"`
	BoundConnectionUUID               string `xml:"boundConnectionUUID,attr"`
	BoundConnectionIp                 string `xml:"boundConnectionIp,attr"`
	BoundConnectionIp6                string `xml:"boundConnectionIp6,attr"`
	BoundClientHostname               string `xml:"boundClientHostname,attr"`
	Nickname                          string `xml:"nickname,attr"`
	ClientId                          string `xml:"clientId,attr"`
	NumConfigurations                 int    `xml:"numConfigurations,attr"`
	NumInterfacesInFirstConfiguration int    `xml:"numInterfacesInFirstConfiguration,attr"`
	FirstInterfaceClass               int    `xml:"firstInterfaceClass,attr"`
	FirstInterfaceSubClass            int    `xml:"firstInterfaceSubClass,attr"`
	FirstInterfaceProtocol            int    `xml:"firstInterfaceProtocol,attr"`
	HideClientInfo                    bool   `xml:"hideClientInfo,attr"`
	AutoUse                           string `xml:"autoUse,attr"`
}

func (d *Device) GetDeviceId() string {
	var flag string
	if strings.TrimSpace(d.DeviceSerial) == `` {
		flag = fmt.Sprintf(`%s:%d-%d:%d-%d:%d:%d:%d`, d.ServerSerial, d.Address, d.VendorId, d.ProductId,
			d.FirstInterfaceClass, d.FirstInterfaceSubClass, d.FirstInterfaceProtocol, d.NumInterfacesInFirstConfiguration)
	} else {
		flag = fmt.Sprintf(`%d:%d-%s`, d.VendorId, d.ProductId, d.DeviceSerial)
	}
	return md5Hash(flag)
}

type UseInfo struct {
	UserId   string `xml:"userId,attr"`
	Username string `xml:"username,attr"`
	Address  string `xml:"address,attr"`
	Hostname string `xml:"hostname,attr"`
}

type DeviceInfo struct {
	Address      string   `xml:"address,attr"`
	VendorId     int      `xml:"vendorId,attr"`
	Vendor       string   `xml:"vendor,attr"`
	ProductId    int      `xml:"productId,attr"`
	Product      string   `xml:"product,attr"`
	SerialNumber string   `xml:"serialNumber,attr"`
	InUseBy      *UseInfo `xml:"inUseBy,attr"`
	IsConnected  bool     `xml:"-"`
}

type ServerInfo struct {
	Name         string        `xml:"name,attr"`
	Version      string        `xml:"version,attr"`
	State        string        `xml:"state,attr"`
	Address      string        `xml:"address,attr"`
	Port         int           `xml:"port,attr"`
	ConnectedFor int           `xml:"connectedFor,attr"`
	MaxDevices   int           `xml:"maxDevices,attr"`
	ConnectedId  int           `xml:"ConnectedId,attr"`
	Interface    string        `xml:"interface,attr"`
	SerialNumber string        `xml:"serialNumber,attr"`
	EasyFind     bool          `xml:"easyFind,attr"`
	Devices      []*DeviceInfo `xml:"devices,attr"`
}

type Server struct {
	Connection *Connection `xml:"connection"`
	Devices    []*Device   `xml:"device"`
}

type ClientState struct {
	Servers []*Server `xml:"server"`
}

type LicenseInfo struct {
	SerialNumber   string `json:"serial_number"`
	IsUnlimited    bool   `json:"is_unlimited"`
	IsUnregistered bool   `json:"is_unregistered"`
	MaxDeviceCount int    `json:"max_device_count"`
}
