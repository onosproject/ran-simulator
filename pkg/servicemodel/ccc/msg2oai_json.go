package ccc

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"os"
)

// json for OAI slicing

type Message struct {
	RrmPolicyRatio []RrmPolicyRatio `json:"rrmPolicyRatio"`
}

type RrmPolicyRatio struct {
	SST            uint8   `json:"sST"`
	SDFlag         float64 `json:"sD_flag"`
	SD             uint32  `json:"sD"`
	MinRatio       uint8   `json:"min_ratio"`
	MaxRatio       uint8   `json:"max_ratio"`
	DedicatedRatio uint8   `json:"dedicated_ratio"`
}

func validateMessage(msg *Message) (bool, error) {
	// Load the JSON schema from the file
	schemaBytes, err := os.ReadFile("/home/garrobo/go/src/github.com/onosproject/ran-simulator/pkg/servicemodel/ccc/policy_schema.json")
	if err != nil {
		return false, fmt.Errorf("failed to read the schema file: %v", err)
	}

	// Create a new JSON schema loader
	loader := gojsonschema.NewBytesLoader(schemaBytes)

	// Load the message data as a JSON document
	dataBytes, err := json.Marshal(msg)
	if err != nil {
		return false, fmt.Errorf("failed to marshal the message: %v", err)
	}
	documentLoader := gojsonschema.NewBytesLoader(dataBytes)

	// Load and compile the schema
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return false, fmt.Errorf("failed to load the schema: %v", err)
	}

	// Validate the document against the schema
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return false, fmt.Errorf("failed to validate the message against the schema: %v", err)
	}

	if result.Valid() {
		return true, nil
	}

	// If validation fails, collect and print the validation errors
	var errors []string
	for _, desc := range result.Errors() {
		errors = append(errors, desc.String())
	}
	log.Error("[CCC]: PolicyRatio Invalid and Discarded!!!")
	return false, fmt.Errorf("validation errors: %v", errors)
}

func Msg2Json(msg Message) {

	rrmPolicyFilePath := "/home/garrobo/oai-slicing-intel/rrmPolicy.json"

	// Validate the message against the schema
	valid, err := validateMessage(&msg)
	if err != nil {
		log.Error("Message validation failed:", err)
		return
	}

	if !valid {
		log.Error("Received message does not follow the required schema.")
		return
	}

	oldJsonData, err := os.ReadFile(rrmPolicyFilePath)
	if err != nil {
		// Convert the JSON message back to a byte slice
		jsonData, err := json.MarshalIndent(msg, "", "  ")
		if err != nil {
			log.Error("Failed to marshal the message:", err)
			return
		}

		// Write the JSON data to a file named rrmPolicy.json
		err = os.WriteFile(rrmPolicyFilePath, jsonData, 0644)
		if err != nil {
			log.Error("Failed to write the JSON file:", err)
			return
		}

		log.Infof("[CCC]: Successfully converted and saved the message as rrmPolicyRatio.json")
	} else {
		// Unmarshal JSON data into a Message
		var oldMsgJson Message
		err = json.Unmarshal(oldJsonData, &oldMsgJson)
		if err != nil {
			log.Error("Failed to Unmarshal the message:", err)
			return
		}

		update_flag := 0
		for i := 0; i < len(oldMsgJson.RrmPolicyRatio); i++ {
			// Assume only one item in msg.RrmPolicyRatio
			if (msg.RrmPolicyRatio[0].SST == oldMsgJson.RrmPolicyRatio[i].SST) && (msg.RrmPolicyRatio[0].SD == oldMsgJson.RrmPolicyRatio[i].SD) {
				// update policy ratio
				oldMsgJson.RrmPolicyRatio[i].SDFlag = msg.RrmPolicyRatio[0].SDFlag
				oldMsgJson.RrmPolicyRatio[i].MinRatio = msg.RrmPolicyRatio[0].MinRatio
				oldMsgJson.RrmPolicyRatio[i].MaxRatio = msg.RrmPolicyRatio[0].MaxRatio
				oldMsgJson.RrmPolicyRatio[i].DedicatedRatio = msg.RrmPolicyRatio[0].DedicatedRatio
				update_flag = 1
				break
			}
		}
		if update_flag == 0 {
			// add new policy ratio
			oldMsgJson.RrmPolicyRatio = append(oldMsgJson.RrmPolicyRatio, msg.RrmPolicyRatio[0])
		}

		// Convert the JSON message back to a byte slice
		newJsonData, err := json.MarshalIndent(oldMsgJson, "", "  ")
		if err != nil {
			log.Error("Failed to marshal the message:", err)
			return
		}

		// Write the JSON data to a file named rrmPolicy.json
		err = os.WriteFile(rrmPolicyFilePath, newJsonData, 0644)
		if err != nil {
			log.Error("Failed to write the JSON file:", err)
			return
		}

		log.Infof("[CCC]: Successfully converted and saved the message as rrmPolicyRatio.json")
	}
}