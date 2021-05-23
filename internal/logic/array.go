package array

import (
	"cerberus-security-laboratories/des-wristband-ui/internal/core"
)

func Remove(slice []core.WristbandProxy, index int) []core.WristbandProxy {
	return append(slice[:index], slice[index + 1:]...)
}
