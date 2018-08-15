package cli

import (
	"errors"
	"fmt"

	"github.com/mitsukomegumi/GoP2P/types/database"
	"github.com/mitsukomegumi/GoP2P/types/node"
)

// handleNewDatabaseCommand - handle execution of handleNewDatabase method (wrapper)
func (term *Terminal) handleNewDatabaseCommand() {
	fmt.Println("attempting to initialize new NodeDatabase") // Log begin

	output, err := term.handleNewDatabase() // Attempt to init new db

	if err != nil { // Check for errors
		fmt.Println("-- ERROR -- " + err.Error()) // Log error
	} else {
		fmt.Println(output) // Log success
	}
}

// handleAddNodeCommand - handle execution of handleAddNode method (wrapper)
func (term *Terminal) handleAddNodeCommand(address string) {
	fmt.Println("attempting to add node " + address + " to database") // Log begin

	output, err := term.handleAddNode(address) // Attempt to append

	if err != nil { // Check for errors
		fmt.Println("-- ERROR -- " + err.Error()) // Log error
	} else {
		fmt.Println(output) // Log success
	}
}

// handleRemoveNodeCommand - handle execution of handleRemoveNode method (wrapper)
func (term *Terminal) handleRemoveNodeCommand(address string) {
	fmt.Println("attempting to remove node " + address + " from database") // Log begin

	output, err := term.handleRemoveNode(address) // Attempt to remove

	if err != nil { // Check for errors
		fmt.Println("-- ERROR -- " + err.Error()) // Log error
	} else {
		fmt.Println(output) // Log success
	}
}

// handleAttachDatabaseCommand - handle execution of handleAttachDatabase method (wrapper)
func (term *Terminal) handleAttachDatabaseCommand() {
	fmt.Println("attempting to attach to NodeDatabase") // Log begin

	output, err := term.handleAttachDatabase() // Attempt to attach to db

	if err != nil { // Check for errors
		fmt.Println("-- ERROR -- " + err.Error()) // Log error
	} else {
		fmt.Println(output) // Log success
	}
}

// handleWriteDatabaseToMemoryCommand - handle execution of handleWritDatabaseToMemory method (wrapper)
func (term *Terminal) handleWriteDatabaseToMemoryCommand() {
	fmt.Println("attempting to write database to memory") // Log begin

	output, err := term.handleWriteDatabaseToMemory() // Attempt to write db

	if err != nil { // Check for errors
		fmt.Println("-- ERROR -- " + err.Error()) // Log error
	} else {
		fmt.Println(output) // Log success
	}
}

// handleNewDatabase - attempt to initialize new NodeDatabase
func (term *Terminal) handleNewDatabase() (string, error) {
	foundNode := node.Node{} // Create placeholder

	for x := 0; x != len(term.Variables); x++ { // Iterate through array
		if term.VariableTypes[x] == "Node" { // Verify element is node
			foundNode = term.Variables[x].(node.Node) // Set to value

			break
		}
	}

	if foundNode.Address == "" { // Check for errors
		return "", errors.New("node not attached") // Log found error
	}

	db, err := database.NewDatabase(&foundNode, 5) // Attempt to create new database

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	term.AddVariable(db, "NodeDatabase") // Add new database

	err = db.WriteToMemory(foundNode.Environment) // Attempt to write to memory

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return "-- SUCCESS -- created new nodedatabase with address " + foundNode.Address, nil // No error occurred, return success
}

// handleAddNode - attempt to append current node to NodeDatabase
func (term *Terminal) handleAddNode(address string) (string, error) {
	if address != "" {
		return term.handleAddSpecificNode(address)
	}

	return term.handleAddCurrentNode()
}

// handleRemoveNode - attempt to remove node from database
func (term *Terminal) handleRemoveNode(address string) (string, error) {
	if address != "" {
		return term.handleRemoveSpecificNode(address)
	}

	return term.handleRemoveCurrentNode()
}

func (term *Terminal) handleAddSpecificNode(address string) (string, error) {
	db, err := term.findDatabase() // Attempt to attach to database

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	_, err = db.QueryForAddress(address)

	if err == nil {
		return "", errors.New("node already added to database")
	}

	newNode, err := node.NewNode(address, false) // Attempt to init node with specified address

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	err = db.AddNode(&newNode) // Attempt to add node

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	err = db.WriteToMemory(term.Variables[0].(node.Node).Environment) // Serialize

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return "-- SUCCESS -- added node with address " + address + " to attached node database", nil // Return success
}

func (term *Terminal) handleRemoveSpecificNode(address string) (string, error) {
	db, err := term.findDatabase()

	if err != nil {
		return "", err
	}

	err = db.RemoveNode(address)

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	err = db.WriteToMemory(term.Variables[0].(node.Node).Environment) // Serialize

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return "-- SUCCESS -- removed node with address " + address + " from attached node database", nil // Return success
}

// handleAddCurrentNode - attempt to add current node to attached node database
func (term *Terminal) handleAddCurrentNode() (string, error) {
	foundNode := node.Node{} // Create placeholder

	for x := 0; x != len(term.Variables); x++ { // Iterate through array
		if term.VariableTypes[x] == "Node" { // Verify element is node
			foundNode = term.Variables[x].(node.Node) // Set to value

			break
		}
	}

	if foundNode.Address == "" { // Check for errors
		return "", errors.New("node not attached") // Log found error
	}

	db, err := term.findDatabase() // Attach to database

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	_, qErr := db.QueryForAddress(foundNode.Address) // Check for already existing node

	if qErr != nil { // Check for already existing node
		err := db.AddNode(&foundNode) // Attempt to add node

		if err != nil { // Check for errors
			return "", err // Return found error
		}
	} else { // Node already exists, return error
		return "", errors.New("node already exists in attached database") // Return found error
	}

	err = db.WriteToMemory(term.Variables[0].(node.Node).Environment) // Serialize

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return "-- SUCCESS -- appended node with address " + foundNode.Address + " to NodeDatabase", nil // No error occurred, return success
}

// handleAddCurrentNode - attempt to add current node to attached node database
func (term *Terminal) handleRemoveCurrentNode() (string, error) {
	foundNode := node.Node{} // Create placeholder

	for x := 0; x != len(term.Variables); x++ { // Iterate through array
		if term.VariableTypes[x] == "Node" { // Verify element is node
			foundNode = term.Variables[x].(node.Node) // Set to value

			break
		}
	}

	if foundNode.Address == "" { // Check for errors
		return "", errors.New("node not attached") // Log found error
	}

	db, err := term.findDatabase()

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	_, qErr := db.QueryForAddress(foundNode.Address) // Check for already existing node

	if qErr != nil { // Check for already existing node
		return "", errors.New("node does not exist in attached database") // Node doesn't exist, return error
	}

	err = db.RemoveNode(foundNode.Address) // Attempt to remove node

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	err = db.WriteToMemory(term.Variables[0].(node.Node).Environment) // Serialize

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return "-- SUCCESS -- removed node with address " + foundNode.Address + " from NodeDatabase", nil // No error occurred, return success
}

// handleAttachDatabase - handle execution of database reading, write to term mem
func (term *Terminal) handleAttachDatabase() (string, error) {
	foundNode := node.Node{} // Create placeholder

	for x := 0; x != len(term.Variables); x++ { // Iterate through array
		if term.VariableTypes[x] == "Node" { // Verify element is node
			foundNode = term.Variables[x].(node.Node) // Set to value

			break
		}
	}

	if foundNode.Address == "" { // Check for errors
		return "", errors.New("node not attached") // Log found error
	}

	db, err := database.ReadDatabaseFromMemory(foundNode.Environment) // Attempt to read database from node environment memory

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	err = term.AddVariable(*db, "NodeDatabase") // Save for persistency

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return "-- SUCCESS -- attached to nodedatabase with bootstrap address " + (*db.Nodes)[0].Address, nil
}

// handleWritDatabaseToMemory - handle execution of NodeDatabase writeToMemory() method
func (term *Terminal) handleWriteDatabaseToMemory() (string, error) {
	foundNode := node.Node{} // Create placeholder

	for x := 0; x != len(term.Variables); x++ { // Iterate through array
		if term.VariableTypes[x] == "Node" { // Verify element is node
			foundNode = term.Variables[x].(node.Node) // Set to value

			break
		}
	}

	if foundNode.Address == "" { // Check for errors
		return "", errors.New("node not attached") // Log found error
	}

	db, err := term.findDatabase() // Attempt to attach to database

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	err = db.WriteToMemory(foundNode.Environment) // Attempt to write to memory

	if err != nil { // Check for errors
		return "", err // Return found error
	}

	return "-- SUCCESS -- wrote nodedatabase with address " + foundNode.Address + " to memory", nil // No error occurred, return success
}

func (term *Terminal) findDatabase() (*database.NodeDatabase, error) {
	foundNode := node.Node{} // Create placeholder

	for x := 0; x != len(term.Variables); x++ { // Iterate through array
		if term.VariableTypes[x] == "Node" { // Verify element is node
			foundNode = term.Variables[x].(node.Node) // Set to value

			break
		}
	}

	if foundNode.Address == "" { // Check for errors
		return &database.NodeDatabase{}, errors.New("node not attached") // Log found error
	}

	db, err := database.ReadDatabaseFromMemory(foundNode.Environment) // Attempt to read database from node environment memory

	if err != nil { // Check for errors
		return &database.NodeDatabase{}, err // Return found error
	}

	return db, nil // No error occurred, return found database
}
