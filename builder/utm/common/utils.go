package common

import (
	"fmt"
	"regexp"
	"strings"
)

// Map of controller names to their corresponding enum codes
var ControllerEnumMap = map[string]string{
	"none":   "QdIn",
	"ide":    "QdIi",
	"scsi":   "QdIs",
	"sd":     "QdId",
	"mtd":    "QdIm",
	"floppy": "QdIf",
	"pflash": "QdIp",
	"virtio": "QdIv",
	"nvme":   "QdIn",
	"usb":    "QdIu",
}

// Function to get the UTM enum code for a given controller name
func GetControllerEnumCode(controllerName string) (string, error) {
	code, exists := ControllerEnumMap[controllerName]
	if !exists {
		return "", fmt.Errorf("invalid controller name: %s", controllerName)
	}
	return code, nil
}

func MajorMinorDriverVersion(version string) string {
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.\d+$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(version))

	if len(matches) == 3 {
		return matches[1] + "_" + matches[2]
	}

	return version
}
