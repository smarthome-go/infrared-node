package hardware

import (
	"fmt"
	"time"

	"github.com/smarthome-go/infrared"
	"github.com/smarthome-go/sdk"

	"github.com/smarthome-go/infrared-node/core/config"
	"github.com/smarthome-go/infrared-node/core/log"
)

// Sets up the scanner and prepare it for operation
func Init(hardware config.Hardware) (infrared.IfScanner, error) {
	ifScanner := infrared.IfScanner{}
	if err := ifScanner.Setup(hardware.Pin); err != nil {
		log.Error("Failed to setup scanner: ", err.Error())
		return infrared.IfScanner{}, err
	}
	log.Trace("Successfully initialized receiver")
	return ifScanner, nil
}

// The `scan` function is launched as a Go-routine which matches the received codes against the ones in the config file
// If a code is matched, the specified Homescript string is executed on the Smarthome server
func Scan(shome *sdk.Connection, conf config.Config, ifScanner infrared.IfScanner) {
	for {
		receivedCode, err := ifScanner.Scan()
		if err != nil {
			log.Error("Failed to scan code, exiting: ", err.Error())
			return
		}
		// Match the received code
		fmt.Println("Code received: ", receivedCode)
		go matchCode(shome, conf, receivedCode)
	}
}

// Is called when a infrared code has been detected
// If a match has been found, the matching Homescript code will be executed
func matchCode(shome *sdk.Connection, conf config.Config, code string) {
	for _, option := range conf.Actions {
		if option.Code == code {
			log.Debug(fmt.Sprintf("Code '%s' matched to action '%s'", option.Code, option.Name))
			res, err := shome.RunHomescriptCode(option.Homescript, make(map[string]string, 0), time.Duration(uint(time.Second)*conf.Smarthome.HmsTimeout))
			if err != nil {
				log.Error("Homescript execution failed: ", err.Error())
				return
			}
			log.Debug(fmt.Sprintf("Exit: %d | Output: %s", res.Exitcode, res.Output))
			log.Debug(fmt.Sprintf("Action '%s' completed successfully", option.Name))
		}
	}
}
