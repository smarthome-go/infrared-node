package smarthome

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MikMuellerDev/smarthome-hw-ir/core/config"
	"github.com/MikMuellerDev/smarthome-hw-ir/core/log"
)

type HomescriptRequest struct {
	Code string `json:"code"`
}

type HomescriptResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

// Here, the network calls to the smarthome server are defined
func SendHomescript(code string) error {
	smarthomeConfig := config.GetConfig().Smarthome
	requestBody, err := json.Marshal(HomescriptRequest{
		Code: code,
	})
	if err != nil {
		log.Error("Could not parse Homescript request: ", err.Error())
		return err
	}
	// Create a client with a more realistic timeout of 1 second
	client := http.Client{Timeout: time.Second}
	res, err := client.Post(
		fmt.Sprintf("%s/api/homescript/run/live?username=%s&password=%s",
			smarthomeConfig.SmarthomeUrl,
			smarthomeConfig.SmarthomeUser,
			smarthomeConfig.SmarthomePassword,
		),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		log.Error("Homescript execution request failed: could not request execution: ", err.Error())
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error(fmt.Sprintf("Received non 200 status code: %d", res.StatusCode))
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal("Failed to read response of Homescript request: ", err.Error())
		}
		var parsedBody HomescriptResponse
		if err := json.Unmarshal(body, &parsedBody); err != nil {
			log.Error("Failed to read response of Homescript request: ", err.Error())
			return nil
		}
		log.Error("Execution error: ", parsedBody.Error)
		log.Error("Execution output: ", parsedBody.Output)
		return errors.New("Homescript execution failed: non 200 status code")
	}
	return nil
}
