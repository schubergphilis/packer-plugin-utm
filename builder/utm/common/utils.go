package common

import "fmt"

// Map of controller names to their corresponding enum codes
var ControllerEnumMap = map[string]string{
	"none":   "Qdin",
	"ide":    "Qdii",
	"scsi":   "Qdis",
	"sd":     "Qdid",
	"mtd":    "Qdim",
	"floppy": "Qdif",
	"pflash": "Qdip",
	"virtio": "Qdiv",
	"nvme":   "Qdin",
	"usb":    "Qdiu",
}

// Function to get the UTM enum code for a given controller name
func GetControllerEnumCode(controllerName string) (string, error) {
	code, exists := ControllerEnumMap[controllerName]
	if !exists {
		return "", fmt.Errorf("invalid controller name: %s", controllerName)
	}
	return code, nil
}
