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
	// Attempt to connect to Smarthome
	c, err := sdk.NewConnection(conf.Smarthome.SmarthomeUrl, sdk.AuthMethodQuery)
	if err != nil {
		log.Error("Could not initialize SDK: invalid Smarthome configuration: ", err.Error())
		os.Exit(1)
	}
	if err := c.Connect(conf.Smarthome.SmarthomeUser, conf.Smarthome.SmarthomePassword); err != nil {
		log.Error("Could establish connection using SDK: invalid Smarthome configuration: ", err.Error())
		os.Exit(1)
	}
	// Test if Homescript can be executed
	if _, err := c.RunHomescriptCode("print('test')", make(map[string]string, 0), time.Second*10); err != nil {
		log.Error("Could not run test Homescript: ", err.Error())
		os.Exit(1)
	}
	// Do not start the scanner if the hardware is disabled
	if !conf.Hardware.HardwareEnabled {
		log.Warn("Hardware is not enabled, exiting")
		os.Exit(0)
	}
	scanner, err := hardware.Init(conf.Hardware)
	if err != nil {
		log.Error("Failed to start service: could not initialize hardware: ", err.Error())
		os.Exit(1)
	}
	hardware.Scan(c, conf, scanner)
}
