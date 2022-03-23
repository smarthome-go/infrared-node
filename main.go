package main

import (
	"github.com/MikMuellerDev/smarthome-hw-ir/core/config"
	"github.com/MikMuellerDev/smarthome-hw-ir/core/hardware"
	"github.com/MikMuellerDev/smarthome-hw-ir/core/log"
	"github.com/sirupsen/logrus"
)

func main() {
	log.InitLogger(logrus.TraceLevel)
	if err := config.ReadConfigFile(); err != nil {
		log.Fatal("Failed to start service: could not read config file: ", err.Error())
	}

	// Do not start the scanner if the hardware is disabled
	if !config.GetConfig().Hardware.HardwareEnabled {
		log.Warn("Hardware is not enabled, exiting")
		return
	}
	if err := hardware.Init(config.GetConfig().Hardware); err != nil {
		log.Fatal("Failed to start service: could not initialize hardware: ", err.Error())
	}
	hardware.Scan()
}
