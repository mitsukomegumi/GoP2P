package command

import (
	"errors"
	"reflect"

	"github.com/dowlandaiello/GoP2P/types/environment"
)

// Command - absctract container holding command values
type Command struct {
	Command string `json:"command"`

	ModifierSet *ModifierSet `json:"modifiers"`
}

// ModifierSet - abstract containers holding specific parameters for a command
type ModifierSet struct {
	Type     string                `json:"type"`
	Value    interface{}           `json:"value"`
	Variable *environment.Variable `json:"variable"`
}

/*
	BEGIN EXPORTED METHODS:
*/

// NewModifierSet - initialize new instance of modifierset struct with modifiers
func NewModifierSet(modifierType string, modifierValue interface{}, modifierVariable *environment.Variable) *ModifierSet {
	return &ModifierSet{Type: modifierType, Value: modifierValue, Variable: modifierVariable} // Initialize modifier set
}

// NewCommand - attempt to initialize new instance of command struct with specified command, modifiers.
func NewCommand(command string, modifierSet *ModifierSet) (*Command, error) {
	if command == "" { // Check for nil command
		return &Command{}, errors.New("invalid command") // Return found error
	} else if reflect.ValueOf(modifierSet).IsNil() { // Check for nil modifier
		return &Command{}, errors.New("invalid modifier") // Return found error
	}

	return &Command{Command: command, ModifierSet: modifierSet}, nil // Return instance
}
