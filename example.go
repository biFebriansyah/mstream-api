package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type testPayload struct {
	Key     string      `json:"Key"`
	Payload interface{} `json:"payload"`
}

func example() {
	message := map[string]string{"uuid": "uuid", "location": "file"}
	toJSON("hello", message)
}

func toJSON(key string, payload interface{}) {
	payloads := testPayload{
		Key:     key,
		Payload: payload,
	}

	data, err := json.Marshal(payloads)
	if err != nil {
		log.Println(err)
	}

	ShowData(data)
}

func ShowData(body []byte) {
	var msg testPayload
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Println(err)
	}

	// Handle Payload as map[string]interface{}
	payloadMap, ok := msg.Payload.(map[string]interface{})
	if ok {
		fmt.Println("Payload as map:", payloadMap)
	}

	// Handle Payload as string if applicable
	payloadString, ok := msg.Payload.(string)
	if ok {
		fmt.Println("Payload as string:", payloadString)
	}

	// Print the raw data type
	fmt.Printf("Payload type: %T\n", msg.Payload)
	fmt.Printf("Payload: %+v\n", msg.Payload)

}
