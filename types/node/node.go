package node

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dowlandaiello/GoP2P/common"
	"github.com/dowlandaiello/GoP2P/types/environment"
)

// Node - abstract struct containing metadata for a node
type Node struct {
	Address      string                   `json:"IP address"`   // Node's IP address
	Reputation   uint                     `json:"reputation"`   // Node's reputation (used for node finding algorithm)
	LastPingTime time.Time                `json:"ping"`         // Last time that the node was pinged successfully (also used for node finding algorithm)
	IsBootstrap  bool                     `json:"is bootstrap"` // Value used for checking whether or not a specific node is a bootstrap node (again, used for node finding algorithm)
	Environment  *environment.Environment `json:"environment"`  // Used for variable storage and referencing
}

/*
	BEGIN EXPORTED METHODS:
*/

// NewNode - create new instance of node struct, with address specified
func NewNode(address string, isBootstrap bool) (Node, error) {
	environment, err := environment.NewEnvironment() // Create new environment

	if err != nil { // Check for errors
		return Node{}, err // Return error
	}

	if address == "" { // Check for invalid address
		return Node{}, errors.New("invalid init values") // Return error
	}

	node := Node{Address: address, Reputation: 0, IsBootstrap: isBootstrap, Environment: environment} // Creates new node instance with specified address

	err = common.CheckAddress(node.Address) // Verify address

	if err != nil { // If node address is invalid, return error
		return Node{}, err // Returns nil node, error
	}

	node.LastPingTime = common.GetCurrentTime() // Since node address is valid, add current time as last ping time
	node.Reputation += common.NodeAvailableRep

	return node, nil // No error occurred, return nil
}

// StartListener - attempt to listen on specified port, return new listener
func (node *Node) StartListener(port int) (*net.Listener, error) {
	ln, err := tls.Listen("tcp", ":"+strconv.Itoa(port), common.GeneralTLSConfig) // Listen on port

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &ln, nil // No error occurred, return listener
}

// LogNode - serialize and print contents of entire node
func (node *Node) LogNode() error {
	marshaledVal, err := json.MarshalIndent(*node, "", "  ") // Marshal node

	if err != nil { // Check for errors
		return err // Return found error
	}

	common.Println("\n" + string(marshaledVal)) // Log marshaled val

	return nil // No error occurred, return nil
}

// String - convert node to string
func (node *Node) String() string {
	marshaledVal, _ := json.MarshalIndent(*node, "", "  ") // Marshal node

	return string(marshaledVal) // No error occurred, return nil
}

// WriteToMemory - create serialized instance of specified environment in specified path (string)
func (node *Node) WriteToMemory(path string) error {
	os.Remove(path + filepath.FromSlash("/node.gob")) // Overwrite

	err := common.WriteGob(path+filepath.FromSlash("/node.gob"), node) // Attempt to write env to path

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil.
}

// ReadNodeFromMemory - read serialized object of specified node from specified path
func ReadNodeFromMemory(path string) (*Node, error) {
	tempNode := new(Node)

	err := common.ReadGob(path+filepath.FromSlash("/node.gob"), tempNode)

	if err != nil { // Check for errors
		return nil, err // Return error
	}
	return tempNode, nil // No error occurred, return nil error, env
}

/*
	END EXPORTED METHODS:
*/
