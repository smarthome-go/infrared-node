package hardware

import (
	"fmt"

	"github.com/MikMuellerDev/smarthome-hw-ir/core/config"
	"github.com/MikMuellerDev/smarthome-hw-ir/core/log"
	"github.com/MikMuellerDev/smarthome-hw-ir/core/smarthome"
)

// Returns the matched Homescript code as a string and a boolean that indicates if a homescript could be matched
func matchCode(code string) {
	for _, option := range config.GetConfig().Actions {
		if option.TriggerCode == code {
			log.Debug(fmt.Sprintf("Code '%s' matched to action '%s'", option.TriggerCode, option.Name))
			output, err := smarthome.SendHomescript(option.ActionHomescript)
			if err != nil {
				log.Error("Homescript execution failed: ", err.Error())
				return
			}
			log.Debug(fmt.Sprintf("Action '%s' completed successfully. Output:\n%s", option.Name, output))
		}
	}
}
