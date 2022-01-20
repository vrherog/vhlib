package VirtualHereLibrary

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net"
	"strings"
)

func md5Hash(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func parseBytesString(raw string) (result string) {
	if strings.TrimSpace(raw) != `` {
		if buffer, err := base64.StdEncoding.DecodeString(raw); err == nil {
			buffer = bytes.TrimFunc(buffer, func(r rune) bool {
				return r == 0
			})
			result = strings.TrimSpace(hex.EncodeToString(buffer))
		}
	}
	return
}

func parseIpAddress(raw string) (result string) {
	if strings.TrimSpace(raw) != `` {
		if buffer, err := base64.StdEncoding.DecodeString(raw); err == nil {
			buffer = bytes.TrimFunc(buffer, func(r rune) bool {
				return r == 0
			})
			if len(buffer) > 0 {
				result = net.IP(buffer).String()
			}
		}
	}
	return
}

func splitMultiLineString(content string, lineHandler func(int, string)) {
	var lines = strings.FieldsFunc(content, func(r rune) bool {
		return r == 10 || r == 13
	})
	for i, line := range lines {
		line = strings.TrimFunc(line, func(r rune) bool {
			return r < 33
		})
		if line != `` {
			lineHandler(i, line)
		}
	}
}
