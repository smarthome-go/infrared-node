package smarthome

import "fmt"

// Here, the network calls to the smarthome server are defined
func SendHomescript(code string) (string, error) {
	fmt.Println("Running code: ", code)
	return "", nil
}
