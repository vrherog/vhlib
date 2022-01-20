package VirtualHereLibrary

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type IClient interface {
	Query(command string) (string, error)
}

type Client struct {
	client IClient

	AutoFindHubs      bool
	AutoUseAllDevices bool
	ReverseLookup     bool
	ReverseSslLookup  bool
}

var reListServer = regexp.MustCompile(`(.+)\((.+):(\d+)\)`)
var reListDevice = regexp.MustCompile(`\s+\-\-\>\s+(.+)\((.+)\.(\d+)\)`)

func (c *Client) List() (result []*ServerInfo, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandList); err == nil {
		var servers = make([]*ServerInfo, 0)
		var server *ServerInfo
		var lines = strings.FieldsFunc(buffer, func(r rune) bool {
			return r == 10 || r == 13
		})
		for i, line := range lines {
			if i < 2 || strings.TrimSpace(line) == `` {
				continue
			}
			var groups = reListServer.FindStringSubmatch(line)
			if len(groups) > 3 {
				server = &ServerInfo{
					Name:    groups[1],
					Address: groups[2],
					Devices: make([]*DeviceInfo, 0),
				}
				if port, e := strconv.Atoi(groups[3]); e == nil {
					server.Port = port
				}
				servers = append(servers, server)
			} else {
				groups = reListDevice.FindStringSubmatch(line)
				if len(groups) > 3 {
					if server != nil {
						var device = DeviceInfo{
							Vendor:  groups[1],
							Address: fmt.Sprintf(`%s:%s`, groups[2], groups[3]),
						}
						server.Devices = append(server.Devices, &device)
					}
				} else {
					if strings.HasPrefix(lines[i], "Auto-Find currently") {
						c.AutoFindHubs = strings.TrimSpace(lines[i][20:]) == "on"
					} else if strings.HasPrefix(lines[i], "Auto-Use All currently") {
						c.AutoUseAllDevices = strings.TrimSpace(lines[i][22:]) == "on"
					} else if strings.HasPrefix(lines[i], "Reverse Lookup currently") {
						c.ReverseLookup = strings.TrimSpace(lines[i][25:]) == "on"
					} else if strings.HasPrefix(lines[i], "Reverse SSL Lookup currently") {
						c.ReverseSslLookup = strings.TrimSpace(lines[i][29:]) == "on"
					} else if strings.HasPrefix(lines[i], "VirtualHere Client is running as a service") {
					}
				}
			}
		}
		if len(servers) > 0 {
			result = make([]*ServerInfo, len(servers))
			for _, item := range servers {
				if server, err = c.GetServerInfo(item.Name); err == nil {
					for _, d := range item.Devices {
						var device *DeviceInfo
						if device, err = c.GetDeviceInfo(d.Address); err == nil {
							server.Devices = append(server.Devices, device)
						}
					}
					result = append(result, server)
				}
			}
		}
	}
	return
}

func (c *Client) GetClientState() (result *ClientState, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandGetClientState)); err == nil {
		if err = xml.Unmarshal([]byte(buffer), &result); err == nil {
			for _, server := range result.Servers {
				if server.Connection != nil {
					server.Connection.UUID = parseBytesString(server.Connection.UUID)
					server.Connection.TransportId = parseBytesString(server.Connection.TransportId)
					server.Connection.EasyFindId = parseBytesString(server.Connection.EasyFindId)
					server.Connection.EasyFindPin = parseBytesString(server.Connection.EasyFindPin)
				}
				for _, device := range server.Devices {
					device.ConnectionUUID = parseBytesString(device.ConnectionUUID)
					device.BoundConnectionUUID = parseBytesString(device.BoundConnectionUUID)
					device.BoundConnectionIp = parseIpAddress(device.BoundConnectionIp)
					device.BoundConnectionIp6 = parseIpAddress(device.BoundConnectionIp6)
				}
			}
		}
	}
	return
}

func (c *Client) GetServerInfo(name string) (result *ServerInfo, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandServerInfo, name)); err == nil {
		result = &ServerInfo{}
		splitMultiLineString(buffer, func(n int, line string) {
			var parts = strings.Split(line, `:`)
			if len(parts) == 2 {
				var value = strings.TrimSpace(parts[1])
				switch strings.TrimSpace(parts[0]) {
				case "NAME":
					result.Name = value
				case "VERSION":
					result.Version = value
				case "STATE":
					result.State = value
				case "ADDRESS":
					result.Address = value
				case "PORT":
					if v, e := strconv.Atoi(value); e == nil {
						result.Port = v
					}
				case "CONNECTED FOR":
					if v, e := strconv.Atoi(strings.Replace(value, `sec`, ``, 1)); e == nil {
						result.ConnectedFor = v
					}
				case "MAX DEVICES":
					if v, e := strconv.Atoi(value); e == nil {
						result.MaxDevices = v
					}
				case "CONNECTION ID":
					if v, e := strconv.Atoi(value); e == nil {
						result.ConnectedId = v
					}
				case "INTERFACE":
					result.Interface = value
				case "SERIAL NUMBER":
					result.SerialNumber = value
				case "EASYFIND":
					result.EasyFind = value != "not enabled"
				}
			}
		})
	}
	return
}

var reDeviceUseInfo = regexp.MustCompile(`(.+)\((.+)\)\s*AT\s*(.+)\s*\((.+)\)`)

func (c Client) GetDeviceInfo(address string) (result *DeviceInfo, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandDeviceInfo, address)); err == nil {
		result = &DeviceInfo{}
		splitMultiLineString(buffer, func(n int, line string) {
			var parts = strings.Split(line, `:`)
			if len(parts) == 2 {
				var value = strings.TrimSpace(parts[1])
				switch strings.TrimSpace(parts[0]) {
				case "ADDRESS":
					result.Address = value
				case "VENDOR":
					result.Vendor = value
				case "VENDOR ID":
					if v, e := strconv.ParseInt(value, 16, 16); e == nil {
						result.VendorId = int(v)
					}
				case "PRODUCT":
					result.Product = value
				case "PRODUCT ID":
					if v, e := strconv.ParseInt(value, 16, 16); e == nil {
						result.ProductId = int(v)
					}
				case "SERIAL":
					result.SerialNumber = value
				case "IN USE BY":
					if value == "YOU" {
						result.IsConnected = true
					} else {
						var groups = reDeviceUseInfo.FindStringSubmatch(value)
						if len(groups) > 4 {
							result.InUseBy = &UseInfo{
								Username: strings.TrimSpace(groups[1]),
								UserId:   strings.TrimSpace(groups[2]),
								Address:  strings.TrimSpace(groups[3]),
								Hostname: strings.TrimSpace(groups[4]),
							}
						}
					}
				}
			}
		})
	}
	return
}

func (c *Client) UseDevice(address, password string) (ok bool, err error) {
	var buffer string
	if password != `` {
		address = fmt.Sprintf(`%s,%s`, address, password)
	}
	if buffer, err = c.client.Query(fmt.Sprintf(CommandUseDevice, address)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) StopUsing(address string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandStopUsing, address)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) StopUsingAll() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandStopUsingAll); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) StopUsingAllLocal() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandStopUsingAllLocal); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ListLicenses() (result map[string]LicenseInfo, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandListLicenses); err == nil {
		result = make(map[string]LicenseInfo, 0)
		splitMultiLineString(buffer, func(n int, line string) {
			var parts = strings.Split(line, `,`)
			if len(parts) > 2 && strings.HasPrefix(parts[1], `s/n=`) {
				var license LicenseInfo
				license.SerialNumber = parts[1][4:]
				var devices = strings.Fields(parts[2])
				if len(devices) > 1 {
					license.IsUnlimited = devices[0] == `unlimited`
					if !license.IsUnlimited {
						license.IsUnregistered = devices[0] == `1`
						if license.IsUnregistered {
							if count, e := strconv.Atoi(devices[0]); e == nil {
								license.MaxDeviceCount = count
							}
						}
					}
				}
				result[parts[0]] = license
			}
		})
	}
	return
}

func (c *Client) LicenseServer(licenseKey string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandLicenseServer, licenseKey)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) DeviceRename(address, nickname string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandDeviceRename, address, nickname)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ServerRename(name, address string, port uint16) (ok bool, err error) {
	var buffer string
	if port > 0 {
		address = fmt.Sprintf(`%s:%d`, address, port)
	}
	if buffer, err = c.client.Query(fmt.Sprintf(CommandServerRename, address, name)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) AutoFind() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandAutoFind); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) AutoUsePort(address string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandAutoUsePort, address)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) AutoUseDevice(address string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandAutoUseDevice, address)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) AutoUseDevicePort(address string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandAutoUseDevicePort, address)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) AutoUseAll() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandAutoUseAll); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) AutoUseClearAll() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandAutoUseClearAll); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) AutoUseHub(address string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandAutoUseHub, address)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ManualHubList() (result []string, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandManualHubList); err == nil {
		result = make([]string, 0)
		splitMultiLineString(buffer, func(n int, line string) {
			result = append(result, line)
		})
	}
	return
}

func (c *Client) ManualHubAdd(address string, port uint16) (ok bool, err error) {
	var buffer string
	if port > 0 {
		address = fmt.Sprintf(`%s:%d`, address, port)
	}
	if buffer, err = c.client.Query(fmt.Sprintf(CommandManualHubAdd, address)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ManualHubAddEasyFind(easyFindAddress string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandManualHubAdd, easyFindAddress)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ManualHubRemove(address string, port uint16) (ok bool, err error) {
	var buffer string
	if port > 0 {
		address = fmt.Sprintf(`%s:%d`, address, port)
	}
	if buffer, err = c.client.Query(fmt.Sprintf(CommandManualHubRemove, address)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ManualHubRemoveEasyFind(easyFindAddress string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandManualHubRemove, easyFindAddress)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ManualHubRemoveAll() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandManualHubRemoveAll); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ListReverse(serverSerial string) (result []string, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandListReverse, serverSerial)); err == nil {
		result = make([]string, 0)
		splitMultiLineString(buffer, func(n int, line string) {
			result = append(result, line)
		})
	}
	return
}

func (c *Client) AddReverse(serverSerial, clientAddress string, port uint16) (ok bool, err error) {
	var buffer string
	if port > 0 {
		clientAddress = fmt.Sprintf(`%s:%d`, clientAddress, port)
	}
	if buffer, err = c.client.Query(fmt.Sprintf(CommandAddReverse, serverSerial, clientAddress)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) RemoveReverse(serverSerial, clientAddress string, port uint16) (ok bool, err error) {
	var buffer string
	if port > 0 {
		clientAddress = fmt.Sprintf(`%s:%d`, clientAddress, port)
	}
	if buffer, err = c.client.Query(fmt.Sprintf(CommandRemoveReverse, serverSerial, clientAddress)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) Reverse() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandReverse); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) SslReverse() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandSslReverse); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) CustomEvent(address, eventName string) (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(fmt.Sprintf(CommandCustomEvent, address, eventName)); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) ClearLog() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandClearLog); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func (c *Client) Exit() (ok bool, err error) {
	var buffer string
	if buffer, err = c.client.Query(CommandExit); err == nil {
		ok = strings.TrimSpace(buffer) == ResultOk
	}
	return
}

func NewClient() (result *Client) {
	result = &Client{}
	switch runtime.GOOS {
	case `windows`:
		result.client = &WindowsClient{}
	case `linux`, `darwin`:
		result.client = &NonWindowsClient{}
	default:
		result = nil
	}
	return
}
