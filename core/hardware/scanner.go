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
func Init(config config.Hardware) (infrared.IfScanner, error) {
	ifScanner := infrared.IfScanner{}
	if err := ifScanner.Setup(config.ScannerDevicePin); err != nil {
		log.Error("Failed to setup scanner: ", err.Error())
		return infrared.IfScanner{}, err
	}
	log.Trace("Successfully initialized receiver")
	return ifScanner, nil
}

// The `StartScan` function is launched as a Go-routine which matches the received codes against the ones in the configuration file
// If a code is matched, the specified Homescript is executed on the Smarthome server
func StartScan(shome *sdk.Connection, conf config.Config, ifScanner infrared.IfScanner) {
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

// Is called when a code has been received
// If a match has been found, the matching Homescript code is executed
func matchCode(shome *sdk.Connection, conf config.Config, code string) {
	for _, option := range conf.Actions {
		if option.TriggerCode == code {
			log.Debug(fmt.Sprintf("Code `%s` has been matched to action `%s`", option.TriggerCode, option.Name))
			res, err := shome.RunHomescriptCode(
				option.ActionHomescript,
				make(map[string]string, 0),
				time.Duration(uint(time.Second)*conf.Smarthome.HmsTimeout),
			)
			if err != nil {
				log.Error("Homescript execution failed: ", err.Error())
				return
			}
			log.Debug(fmt.Sprintf("Homescript execution completed\n Exit-code: %d, Output: `%s`", res.Exitcode, res.Output))
			log.Debug(fmt.Sprintf("Action `%s` completed successfully", option.Name))
		}
	}
}
