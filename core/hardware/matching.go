package hardware

import (
	"fmt"

	"github.com/smarthome-go/infrared-node/core/config"
	"github.com/smarthome-go/infrared-node/core/log"
	"github.com/smarthome-go/infrared-node/core/smarthome"
)

// Returns the matched Homescript code as a string and a boolean that indicates if a homescript could be matched
func matchCode(code string) {
	for _, option := range config.GetConfig().Actions {
		if option.TriggerCode == code {
			log.Debug(fmt.Sprintf("Code '%s' matched to action '%s'", option.TriggerCode, option.Name))
			if err := smarthome.SendHomescript(option.ActionHomescript); err != nil {
				log.Error("Homescript execution failed: ", err.Error())
				return
			}
			log.Debug(fmt.Sprintf("Action '%s' completed successfully", option.Name))
		}
	}
}
