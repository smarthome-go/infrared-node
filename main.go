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
	// Decide whether to use token or password authentication
	var conn *sdk.Connection
	var connErr error
	// Create a new (different) connection based on whether to use token auth
	if conf.Smarthome.TokenAuth {
		conn, connErr = sdk.NewConnection(conf.Smarthome.URL, sdk.AuthMethodQueryToken)
	} else {
		conn, connErr = sdk.NewConnection(conf.Smarthome.URL, sdk.AuthMethodQueryPassword)
	}
	// Handle any connection creation errors here
	if connErr != nil {
		log.Error("Could not initialize SDK: invalid Smarthome configuration: ", connErr.Error())
		os.Exit(1)
	}

	// Attempt to connect to Smarthome
	var loginErr error
	// Use different login methods depending on whether to use a token or not
	if conf.Smarthome.TokenAuth {
		loginErr = conn.TokenLogin(conf.Smarthome.Credentials.Token)
	} else {
		loginErr = conn.UserLogin(conf.Smarthome.Credentials.Username, conf.Smarthome.Credentials.Password)
	}
	if loginErr != nil {
		log.Error("Could establish connection using SDK: invalid Smarthome configuration: ", loginErr.Error())
		os.Exit(1)
	}

	// Test if Homescript code can be executed
	if _, err := conn.RunHomescriptCode(
		"print('test')",
		make(map[string]string, 0),
		time.Second*10,
	); err != nil {
		log.Error("Could not run test Homescript: ", err.Error())
		os.Exit(1)
	}
	log.Trace("Test Homescript was successfully executed. Smarthome configuration is valid.")

	// Do not start the scanner if the hardware is disabled
	if !conf.Hardware.Enabled {
		log.Warn("Hardware is disabled stopping the service...")
		os.Exit(0)
	}

	// Initialize hardware
	log.Debug(fmt.Sprintf("Initializing infrared scanner on port %d", conf.Hardware.Pin))
	scanner, err := hardware.Init(conf.Hardware)
	if err != nil {
		log.Error("Failed to start service: could not initialize hardware: ", err.Error())
		os.Exit(1)
	}
	log.Info(fmt.Sprintf("Smarthome-hw-ir %s is listening on pin %d", Version, conf.Hardware.Pin))
	hardware.StartScan(conn, conf, scanner)
}
