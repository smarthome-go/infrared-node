package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/smarthome-go/infrared-node/core/log"
)

type Config struct {
	Hardware  Hardware  `json:"hardware"`
	Smarthome Smarthome `json:"smarthome"`
	Actions   []Action  `json:"actions"`
}

// Documentation of following parameters: github.com/smarthome-go/infrared
type Hardware struct {
	HardwareEnabled  bool  `json:"hardwareEnabled"`
	ScannerDevicePin uint8 `json:"pin"` // The BCM pin to which a infrared receiver is attached
}

type Smarthome struct {
	SmarthomeUrl  string `json:"url"`
	SmarthomeUser string `json:"user"`
	// The password should later be replaced with an access token
	SmarthomePassword string `json:"password"`
	// Specifies how long the SDK waits before abandoning a HMS request
	HmsTimeout uint `json:"hmsTimeout"`
}

// Specifies what to do when a code is matched
type Action struct {
	TriggerCode      string `json:"trigger"` // The received infrared code
	ActionHomescript string `json:"action"`  // The action which is executed when the trigger matches
	Name             string `json:"name"`
}

// The path were the config file is located
const configPath = "./config.json"

// A dry-run of the `RadConfigFile()` method used in the Health test
func ProbeConfigFile() error {
	// Read file from <configPath> on disk
	// If this file does not exist, return an error
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Error("Failed to read config file: ", err.Error())
		return nil
	}
	// Parse config file to a test struct <Config>
	var configFile Config
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&configFile)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse config file at `%s` into Config struct: %s", configPath, err.Error()))
		return err
	}
	return nil
}

// Reads the config file from disk, if the file does not exist (for example first run), a new one is created
func ReadConfigFile() (Config, error) {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		configTemp, errCreate := createNewConfigFile()
		if errCreate != nil {
			log.Error("Failed to read config file: ", err.Error())
			log.Fatal("Failed to initialize config: could not read or create a config file: ", errCreate.Error())
			return Config{}, err
		}
		log.Info("Failed to read config file: created a new config file")
		return configTemp, nil
	}
	// Parse config file to struct <Config>
	var configFile Config
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&configFile)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse config file at `%s` into Config struct: %s", configPath, err.Error()))
		return Config{}, err
	}
	return configFile, nil
}

// Creates an empty config file, can return an error if it fails
func createNewConfigFile() (Config, error) {
	config := Config{
		Hardware: Hardware{
			HardwareEnabled:  false,
			ScannerDevicePin: 0,
		},
		Smarthome: Smarthome{
			SmarthomeUrl:      "http://smarthome.box",
			SmarthomeUser:     "admin",
			SmarthomePassword: "admin",
			HmsTimeout:        10,
		},
		Actions: []Action{
			{
				TriggerCode:      "2a00aaa95",
				ActionHomescript: "switch('sx', on)",
				Name:             "demo",
			},
		},
	}
	fileContent, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		log.Error("Failed to create config file: creating file content from JSON failed: ", err.Error())
		return Config{}, err
	}
	if err = ioutil.WriteFile("./config.json", fileContent, 0644); err != nil {
		log.Error("Failed to write file to disk: ", err.Error())
		return Config{}, err
	}
	return config, nil
}

// Writes the current state of the global config to the file on disk
func WriteConfig(config Config) error {
	var jsonBlob = []byte(`{}`)
	err := json.Unmarshal(jsonBlob, &config)
	if err != nil {
		log.Fatal("Error during unmarshal: ", err.Error())
		return err
	}
	configJson, _ := json.MarshalIndent(&config, "", "    ")
	err = ioutil.WriteFile("./config.json", configJson, 0644)
	if err != nil {
		log.Fatal("Error writing new token hash to config.json: ", err.Error())
		return err
	}
	log.Debug("Written to config.json")
	return nil
}
