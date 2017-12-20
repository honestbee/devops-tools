// ref: github.com/adammck/terraform-inventory
// see also hashicorp/terraform@master/-/blob/terraform/state.go

package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type (
	// State keeps track of a snapshot state-of-the-world that Terraform
	// can use to keep track of what real world resources it is actually
	// managing
	state struct {
		// Version is the state file protocol version
		Version int `json:"version"`
		// TFVersion is the version of Terraform that wrote this state.
		TFVersion string `json:"terraform_version,omitempty"`
		// Serial is incremented for any operation to detect conflicts
		Serial int64 `json:"serial"`
		// Modules contains all the modules in a breadth-first order
		Modules []moduleState `json:"modules"`
	}
	// ModuleState is used to track all the state relevant to a single module
	moduleState struct {
		Path      []string                 `json:"path"`
		Resources map[string]resourceState `json:"resources"`
		Outputs   map[string]outputState   `json:"outputs"`
	}
	outputState struct {
		// Sensitive describes whether the output is considered sensitive,
		// which may lead to masking the value on screen in some cases.
		Sensitive bool `json:"sensitive"`
		// Type describes the structure of Value. Valid values are "string",
		// "map" and "list"
		Type string `json:"type"`
		// Value contains the value of the output, in the structure described
		// by the Type field.
		Value interface{} `json:"value"`
	}
	resourceState struct {
		// Populated from statefile
		Type string `json:"type"`
		// Primary is the current active instance for this resource
		Primary instanceState `json:"primary"`
	}
	instanceState struct {
		// A unique ID for this resource. This is opaque to Terraform
		// and is only meant as a lookup mechanism for the providers.
		ID string `json:"id"`
		// Attributes are basic information about the resource.
		Attributes map[string]string `json:"attributes,omitempty"`
	}
)

// read populates the state object from a statefile.
func (s *state) read(stateFile io.Reader) error {

	// read statefile contents
	b, err := ioutil.ReadAll(stateFile)
	if err != nil {
		return err
	}

	// parse into struct
	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	return nil
}
