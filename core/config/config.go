package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/smarthome-go/infrared-node/core/log"
)

var Version string

type Config struct {
	Smarthome Smarthome `json:"smarthome"`
	Hardware  Hardware  `json:"hardware"`
	Actions   []Action  `json:"actions"`
}

// Documentation of following parameters: github.com/smarthome-go/infrared
type Hardware struct {
	Enabled bool  `json:"enabled"`
	Pin     uint8 `json:"pin"` // The BCM pin to which an infrared receiver is attached
}

type Smarthome struct {
	URL         string      `json:"url"`
	TokenAuth   bool        `json:"tokenAuth"`
	Credentials Credentials `json:"credentials"`
	// Specifies how long the SDK waits before abandoning a HMS Homescript request
	HmsTimeout uint `json:"hmsTimeout"`
}

type Credentials struct {
	Username string `json:"user"`
	Password string `json:"password"`
	// Or if an authentication token should be used
	Token string `json:"token"`
}

// Specifies what to do when a code is matched
type Action struct {
	Name       string `json:"name"`       // A friendly name for easy recognition
	Code       string `json:"code"`       // The received infrared code
	Homescript string `json:"homescript"` // The action which is executed when the trigger matches
}

// The path were the config file is located
const configPath = "./config.json"

// A dry-run of the `RadConfigFile()` method used in the Health test
func ProbeConfigFile() error {
	// Read file from <configPath> on disk
	// If this file does not exist, return an error
	content, err := os.ReadFile(configPath)
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
	content, err := os.ReadFile(configPath)
	if err != nil {
		configTemp, errCreate := createNewConfigFile()
		if errCreate != nil {
			log.Error("Failed to read config file: ", err.Error())
			log.Fatal("Failed to initialize config: could not read or create a config file: ", errCreate.Error())
			return Config{}, err
		}
		log.Info(fmt.Sprintf("Created a new configuration file at `%s`", configPath))
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
	// Check if the user entered dubious values in the configuration file
	if (configFile.Smarthome.Credentials.Username != "" || configFile.Smarthome.Credentials.Password != "") && configFile.Smarthome.TokenAuth {
		log.Warn("Username and / or password not empty whilst using token authentication: ignoring username and password")
	}
	return configFile, nil
}

// Creates an empty config file, can return an error if it fails
func createNewConfigFile() (Config, error) {
	config := Config{
		Hardware: Hardware{
			Enabled: false,
			Pin:     0,
		},
		Smarthome: Smarthome{
			URL: "http://smarthome.box",
			Credentials: Credentials{
				Token:    "your-token-here",
				Username: "",
				Password: "",
			},
			HmsTimeout: 10,
		},
		Actions: []Action{
			{
				Name:       "demo",
				Code:       "2a00aaa95",
				Homescript: "switch('sx', on)",
			},
		},
	}
	fileContent, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		log.Error("Failed to create config file: creating file content from JSON failed: ", err.Error())
		return Config{}, err
	}
	if err = os.WriteFile("./config.json", fileContent, 0644); err != nil {
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
	err = os.WriteFile("./config.json", configJson, 0644)
	if err != nil {
		log.Fatal("Error writing new token hash to config.json: ", err.Error())
		return err
	}
	log.Debug("Written to config.json")
	return nil
}
