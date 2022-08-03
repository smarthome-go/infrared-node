package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/smarthome-go/infrared-node/core/config"
	"github.com/smarthome-go/infrared-node/core/hardware"
	"github.com/smarthome-go/infrared-node/core/log"
	"github.com/smarthome-go/sdk"
)

const Version = "1.3.0"

func main() {
	if err := log.InitLogger(logrus.TraceLevel); err != nil {
		fmt.Println("Failed to initialize logger: ", err.Error())
		os.Exit(1)
	}
	conf, err := config.ReadConfigFile()
	if err != nil {
		log.Error("Failed to start service: could not read config file: ", err.Error())
		os.Exit(1)
	}

	// Create a connection to Smarthome
	log.Debug(fmt.Sprintf("Initiating connection to Smarthome (`%s@%s`)", conf.Smarthome.SmarthomeUser, conf.Smarthome.SmarthomeUrl))
	smarthomeConnection, err := sdk.NewConnection(
		conf.Smarthome.SmarthomeUrl,
		sdk.AuthMethodQuery,
	)
	if err != nil {
		log.Error("Could not initialize SDK: invalid Smarthome configuration: ", err.Error())
		os.Exit(1)
	}

	// Authenticate
	if err := smarthomeConnection.Connect(
		conf.Smarthome.SmarthomeUser,
		conf.Smarthome.SmarthomePassword,
	); err != nil {
		log.Error("Could establish connection using SDK: invalid Smarthome configuration: ", err.Error())
		os.Exit(1)
	}
	log.Trace(fmt.Sprintf("Successfully established connection to Smarthome (`%s@%s`)", conf.Smarthome.SmarthomeUser, conf.Smarthome.SmarthomeUrl))

	// Test if Homescript can be executed
	log.Debug("Executing test Homescript...")
	if _, err := smarthomeConnection.RunHomescriptCode(
		"print('test')",
		make(map[string]string, 0),
		time.Second*10,
	); err != nil {
		log.Error("Could not run test Homescript: ", err.Error())
		os.Exit(1)
	}
	log.Trace("Test Homescript was successfully executed. Smarthome configuration is valid.")

	// Do not start the scanner if the hardware is disabled
	if !conf.Hardware.HardwareEnabled {
		log.Warn("Hardware is not enabled, exiting")
		os.Exit(0)
	}

	// Initialize hardware
	log.Debug(fmt.Sprintf("Initializing infrared scanner on port %d", conf.Hardware.ScannerDevicePin))
	scanner, err := hardware.Init(conf.Hardware)
	if err != nil {
		log.Error("Failed to start service: could not initialize hardware: ", err.Error())
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("Smarthome-hw-ir %s is listening on pin %d", Version, conf.Hardware.ScannerDevicePin))

	// Start receiving codes
	hardware.StartScan(
		smarthomeConnection,
		conf,
		scanner,
	)
}
