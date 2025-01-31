package cli

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/dowlandaiello/GoP2P/common"
	commonProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/common"
	databaseProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/database"
	environmentProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/environment"
	handlerProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/handler"
	nodeProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/node"
	protoProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/protobuf"
	shardProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/shard"
	upnpProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/upnp"
	"github.com/fatih/color"
)

// Terminal - absctract container holding set of variable with values (runtime only)
type Terminal struct {
	Variables []Variable
}

// Variable - container holding variable values
type Variable struct {
	VariableName string      `json:"name"`
	VariableData interface{} `json:"data"`
	VariableType string      `json:"type"`
}

// NewTerminal - attempts to start io handler for term commands
func NewTerminal(rpcPort uint, rpcAddress string) error {
	reader := bufio.NewScanner(os.Stdin) // Init reader

	transport := &http.Transport{ // Init transport
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	for {
		common.Print("\n> ") // Print prompt

		reader.Scan() // Scan

		input := reader.Text() // Fetch string input

		input = strings.TrimSuffix(input, "\n") // Trim newline

		err := handleStringOnly(input) // Handle nil-call

		if err != nil {
			receiver, methodname, params, err := common.ParseStringMethodCall(input) // Attempt to parse as method call

			if err != nil { // Check for errors
				common.Println(err.Error()) // Log found error

				continue
			}

			handleCommand(receiver, methodname, params, rpcPort, rpcAddress, transport) // Handle command
		}
	}
}

func handleCommand(receiver string, methodname string, params []string, rpcPort uint, rpcAddress string, transport *http.Transport) {
	nodeClient := nodeProto.NewNodeProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                      // Init node client
	handlerClient := handlerProto.NewHandlerProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})             // Init handler client
	environmentClient := environmentProto.NewEnvironmentProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport}) // Init environment client
	upnpClient := upnpProto.NewUpnpProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                      // Init upnp client
	databaseClient := databaseProto.NewDatabaseProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})          // Init database client
	commonClient := commonProto.NewCommonProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                // Init common client
	shardClient := shardProto.NewShardProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                   // Init shard client
	protoClient := protoProto.NewProtoProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})                   // Init proto client

	switch receiver {
	case "node":
		err := handleNode(&nodeClient, methodname, params) // Handle node

		if err != nil { // Check for errors
			common.Println("\n" + err.Error()) // Log found error
		}
	case "handler":
		err := handleHandler(&handlerClient, methodname, params) // Handle handler

		if err != nil { // Check for errors
			common.Println("\n" + err.Error()) // Log found error
		}
	case "environment":
		err := handleEnvironment(&environmentClient, methodname, params) // Handle environment

		if err != nil { // Check for errors
			common.Println("\n" + err.Error()) // Log found error
		}
	case "upnp":
		err := handleUpnp(&upnpClient, methodname, params) // Handle upnp

		if err != nil { // Check for errors
			common.Println("\n" + err.Error()) // Log found error
		}
	case "database":
		err := handleDatabase(&databaseClient, methodname, params) // Handle database

		if err != nil { // Check for errors
			common.Println("\n" + err.Error()) // Log found error
		}
	case "common":
		err := handleCommon(&commonClient, methodname, params) // Handle common

		if err != nil { // Check for errors
			common.Println("\n" + err.Error()) // Log found error
		}
	case "shard":
		handleShard(&shardClient, methodname, params) // Handle shard
	case "proto":
		handleProto(&protoClient, methodname, params) // Handle proto
	}
}

func handleNode(nodeClient *nodeProto.Node, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NewNode":
		if len(params) != 2 { // Check for insufficient parameters
			return errors.New("invalid parameters (requires string, int)") // Return error
		}

		boolVal, _ := strconv.ParseBool(params[1]) // Parse isBootstrap

		reflectParams = append(reflectParams, reflect.ValueOf(&nodeProto.GeneralRequest{Address: params[0], IsBootstrap: boolVal})) // Append params
	case "StartListener":
		intVal, _ := strconv.Atoi(params[0]) // Get int val

		reflectParams = append(reflectParams, reflect.ValueOf(&nodeProto.GeneralRequest{Port: uint32(intVal)})) // Append params
	case "WriteToMemory", "ReadFromMemory":
		reflectParams = append(reflectParams, reflect.ValueOf(&nodeProto.GeneralRequest{Path: params[0]})) // Append params
	case "LogNode":
		reflectParams = append(reflectParams, reflect.ValueOf(&nodeProto.GeneralRequest{})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NewNode(), StartListener(), LogNode() WriteToMemory(), ReadFromMemory()") // Return error
	}

	result := reflect.ValueOf(*nodeClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*nodeProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	common.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

func handleHandler(handlerClient *handlerProto.Handler, methodname string, params []string) error {
	if len(params) == 0 { // Check for nil parameters
		return errors.New("invalid parameters (requires at least 1 parameter)") // Return error
	}

	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "StartHandler":
		port, _ := strconv.Atoi(params[0]) // Parse port

		reflectParams = append(reflectParams, reflect.ValueOf(&handlerProto.GeneralRequest{Port: uint32(port)})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: StartHandler()") // Return error
	}

	result := reflect.ValueOf(*handlerClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*handlerProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	common.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

func handleEnvironment(environmentClient *environmentProto.Environment, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NewEnvironment", "LogEnvironment":
		reflectParams = append(reflectParams, reflect.ValueOf(&environmentProto.GeneralRequest{})) // Append empty request
	case "QueryType":
		if len(params) != 1 { // Check for errors
			return errors.New("invalid parameters (requires string)") // Return found error
		}

		queryTypeVal := params[0] // Fetch queryTypeVal

		reflectParams = append(reflectParams, reflect.ValueOf(&environmentProto.GeneralRequest{VariableType: queryTypeVal})) // Append querytype request
	case "QueryValue":
		if len(params) != 1 { // Check for errors
			return errors.New("invalid parameters (requires string)") // Return found error
		}

		queryValueVal := params[0] // Fetch query val

		reflectParams = append(reflectParams, reflect.ValueOf(&environmentProto.GeneralRequest{Value: queryValueVal})) // Append queryval request
	case "NewVariable":
		if len(params) != 2 { // Check for errors
			return errors.New("invalid parameters (requires string, string)") // Return found error
		}

		variablePathVal := params[0] // Fetch variable data path
		variableTypeVal := params[1] // Fetch variable type

		reflectParams = append(reflectParams, reflect.ValueOf(&environmentProto.GeneralRequest{Path: variablePathVal, VariableType: variableTypeVal})) // Append path request
	case "AddVariable", "WriteToMemory", "ReadFromMemory":
		if len(params) != 1 { // Check for errors
			return errors.New("invalid parameters (requires string)") // Return found error
		}

		pathVal := params[0] // Fetch variable data path

		reflectParams = append(reflectParams, reflect.ValueOf(&environmentProto.GeneralRequest{Path: pathVal})) // Append path request
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NewEnvironment(), LogEnvironment(), QueryType(), QueryValue(), NewVariable(), AddVariable(), WriteToMemory(), ReadFromMemory()") // Return error
	}

	result := reflect.ValueOf(*environmentClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*environmentProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	common.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

func handleUpnp(upnpClient *upnpProto.Upnp, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "GetGateway":
		reflectParams = append(reflectParams, reflect.ValueOf(&upnpProto.GeneralRequest{})) // Append params
	case "ForwardPortSilent", "ForwardPort", "RemoveForwarding":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires uint32)") // Return error
		}

		port, err := strconv.Atoi(params[0]) // Convert to int

		if err != nil { // Check for errors
			return err // Return found error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&upnpProto.GeneralRequest{PortNumber: uint32(port)})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: GetGateway(), ForwardPortSilent(), ForwardPort(), RemoveForwarding()") // Return error
	}

	result := reflect.ValueOf(*upnpClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*upnpProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	common.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

func handleDatabase(databaseClient *databaseProto.Database, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NewDatabase":
		if len(params) != 4 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string, uint32, uint32, string)") // Return error
		}

		acceptableTimeout, _ := strconv.Atoi(params[2]) // Fetch acceptable timeout

		networkID, _ := strconv.Atoi(params[1]) // Fetch network id

		reflectParams = append(reflectParams, reflect.ValueOf(&databaseProto.GeneralRequest{NetworkName: params[0], NetworkID: uint32(networkID), AcceptableTimeout: uint32(acceptableTimeout), PrivateKey: params[len(params)-1]})) // Append params
	case "AddNode", "UpdateRemoteDatabase", "LogDatabase":
		if len(params) != 1 { // Check for valid parameters
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&databaseProto.GeneralRequest{NetworkName: params[0]})) // Append nil params
	case "JoinDatabase", "FetchRemoteDatabase":
		if len(params) != 3 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string, uint32, string)") // Return error
		}

		intVal, _ := strconv.Atoi(params[1]) // Convert port to uint

		reflectParams = append(reflectParams, reflect.ValueOf(&databaseProto.GeneralRequest{Address: params[0], Port: uint32(intVal), NetworkName: params[2]})) // Append params
	case "WriteToMemory", "ReadFromMemory", "RemoveNode", "QueryForAddress":
		if len(params) != 2 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string, string)") // Return error
		}

		address := params[0] // Fetch removal address

		reflectParams = append(reflectParams, reflect.ValueOf(&databaseProto.GeneralRequest{Address: address, NetworkName: params[1]})) // Append params
	case "FromBytes":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires []byte)") // Return error
		}

		byteVal := []byte(params[0]) // Fetch byte val

		reflectParams = append(reflectParams, reflect.ValueOf(&databaseProto.GeneralRequest{ByteVal: byteVal})) // Append params
	case "SendDatabaseMessage":
		if len(params) != 6 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string, uint32, string, string, string, uint)") // Return error
		}

		uintVal, _ := strconv.Atoi(params[len(params)-1]) // Convert to uint
		portIntVal, _ := strconv.Atoi(params[1])          // Convert to uint

		reflectParams = append(reflectParams, reflect.ValueOf(&databaseProto.GeneralRequest{NetworkName: params[0], Port: uint32(portIntVal), PrivateKey: params[2], UintVal: uint32(uintVal), StringVals: params[3:5]})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NewDatabase(), LogDatabase(), AddNode(), UpdateRemoteDatabase(), JoinDatabase(), FetchRemoteDatabase(), RemoveNode(), QueryForAddress(), WriteToMemory(), ReadFromMemory(), FromBytes(), SendDatabaseMessage()") // Return error
	}

	result := reflect.ValueOf(*databaseClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*databaseProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	common.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

func handleCommon(commonClient *commonProto.Common, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "SeedAddress":
		if len(params) != 3 { // Check for invalid parameters
			return errors.New("invalid parameters (requires []string, string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&commonProto.GeneralRequest{Inputs: params[0 : len(params)-1], SecondInput: params[len(params)-1]})) // Append params
	case "ParseStringMethodCall", "ParseStringParams", "StringStripReceiverCall", "StringStripParentheses", "StringFetchCallReceiver", "CheckAddress":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&commonProto.GeneralRequest{Input: params[0]})) // Append params
	case "ConvertStringToReflectValues":
		if len(params) < 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&commonProto.GeneralRequest{Inputs: params})) // Append params
	case "Sha3":
		if len(params) != 1 { // Check for invalid parameters
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&commonProto.GeneralRequest{ByteInput: []byte(params[0])})) // Append params
	case "SendBytes":
		if len(params) != 2 { // Check for invalid parameters
			return errors.New("invalid parameters (requires []byte, string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&commonProto.GeneralRequest{ByteInput: []byte(params[0]), Input: params[1]})) // Append params
	case "GetExtIPAddrWithUPnP", "GetExtIPAddrWithoutUPnP", "GetCurrentTime", "GetCurrentDir":
		reflectParams = append(reflectParams, reflect.ValueOf(&commonProto.GeneralRequest{})) // Append empty params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: ParseStringMethodCall(), ParseStringParams(), StringStripReceiverCall(), StringStripParentheses(), StringFetchCallReceiver(), CheckAddress(), ConvertStringToReflectValues(), Sha3(), SendBytes(), GetExtIPAddrWithUPnP(), GetExtIPAddrWithoutUPnP(), GetCurrentTime(), GetCurrentDir()") // Return error
	}

	result := reflect.ValueOf(*commonClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*commonProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	common.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

func handleShard(shardClient *shardProto.Shard, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NewShard", "LogShard":
		if len(params) != 2 { // Check for invalid parameters length
			return errors.New("invalid parameters (requires string, string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&shardProto.GeneralRequest{NetworkName: params[0], Address: params[1]})) // Append params
	case "NewShardWithNodes":
		if len(params) < 2 { // Check for invalid parameters length
			return errors.New("invalid parameters (requires string, []string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&shardProto.GeneralRequest{NetworkName: params[0], Addresses: params[1:]})) // Append params
	case "Shard":
		if len(params) != 3 { // Check for invalid parameters length
			return errors.New("invalid parameters (requires string, string, uint32)") // Return error
		}

		exponent, _ := strconv.Atoi(params[2]) // Fetch int val

		reflectParams = append(reflectParams, reflect.ValueOf(&shardProto.GeneralRequest{NetworkName: params[0], Address: params[1], Exponent: uint32(exponent)})) // Append params
	case "CalculateQuadraticExponent":
		if len(params) != 1 { // Check for invalid parameters length
			return errors.New("invalid parameters (requires uint32)") // Return error
		}

		exponent, _ := strconv.Atoi(params[0]) // Fetch int val

		reflectParams = append(reflectParams, reflect.ValueOf(&shardProto.GeneralRequest{Exponent: uint32(exponent)})) // Append params
	case "SendBytesShardResult", "SendBytesShard":
		if len(params) < 3 { // Check for invalid parameters length
			return errors.New("invalid parameters (requires string, uint32, []byte)") // Return error
		}

		port, _ := strconv.Atoi(params[1]) // Fetch int val

		reflectParams = append(reflectParams, reflect.ValueOf(&shardProto.GeneralRequest{Address: params[0], Port: uint32(port), Bytes: []byte(strings.Join(params[2:], " "))})) // Append params
	case "QueryForAddress":
		if len(params) < 3 { // Check for invalid parameters length
			return errors.New("invalid parameters (require string, string, string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&shardProto.GeneralRequest{NetworkName: params[0], Addresses: params[1:]})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NewShard(), NewShardWithNodes(), Shard(), QueryForAddress(), LogShard(), CalculateQuadraticExponent(), SendBytesShardResult(), SendBytesShard()") // Return error
	}

	result := reflect.ValueOf(*shardClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*shardProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		common.Println(result[1].Interface().(error)) // Log
	} else {
		common.Println(response.Message) // Log response
	}

	return nil // No error occurred, return nil
}

func handleProto(protoClient *protoProto.Proto, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname {
	case "NewProtobufGuide":
		if len(params) != 2 { // Check for invalid parameters length
			return errors.New("invalid parameters (requires string, string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&protoProto.GeneralRequest{ProtoID: params[0], Path: params[1]})) // Append params
	case "ReadGuideFromMemory", "WriteToMemory":
		if len(params) != 1 { // Check for invalid parameters length
			return errors.New("invalid parameters (requires string)") // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&protoProto.GeneralRequest{Path: params[0]})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: NewProtobufGuide(), ReadGuideFromMemory(), WriteToMemory()") // Return error
	}

	result := reflect.ValueOf(*protoClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*protoProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		common.Println(result[1].Interface().(error)) // Log
	} else {
		common.Println(response.Message) // Log response
	}

	return nil // No error occurred, return nil
}

// AddVariable - attempt to append specified variable to terminal variable list
func (term *Terminal) AddVariable(variableName string, variableData interface{}, variableType string) error {
	variable := Variable{VariableName: variableName, VariableData: variableData, VariableType: variableType}

	if reflect.ValueOf(term).IsNil() { // Check for nil variable
		return errors.New("nil terminal found") // Return error
	}

	if len(term.Variables) == 0 { // Check for uninitialized variable array
		term.Variables = []Variable{variable} // Initialize with variable

		return nil // No error occurred, return nil
	}

	term.Variables = append(term.Variables, variable) // Append to array

	return nil // No error occurred, return nil
}

// ReplaceVariable - attempt to replace value at index with specified variable
func (term *Terminal) ReplaceVariable(variableIndex int, variableData interface{}) error {
	if reflect.ValueOf(term).IsNil() { // Check for nil variable
		return errors.New("nil terminal found") // Return error
	}

	if len(term.Variables) == 0 { // Check for uninitialized variable array
		return errors.New("empty terminal environment") // Return found error
	}

	(*term).Variables[variableIndex].VariableData = variableData // Replace value

	return nil
}

// QueryType - attempt to fetch index of variable with matching type
func (term *Terminal) QueryType(variableType string) (uint, error) {
	if variableType == "" { // Check for nil parameter
		return 0, errors.New("invalid type") // Return found error
	}

	if len(term.Variables) == 0 { // Check that terminal environment is not nil
		return 0, errors.New("empty terminal environment") // Return found error
	}

	for x := 0; x != len(term.Variables); x++ { // Declare, increment iterator
		if term.Variables[x].VariableType == variableType { // Check for match
			return uint(x), nil // Return result
		}
	}

	return 0, errors.New("couldn't find matching variable") // Return error
}

// hasVariableSet - checks if specified command sets a variable
func hasVariableSet(command string) bool {
	if strings.HasPrefix(strings.ToLower(command), "var") { // Check for prefix
		return true
	}

	return false
}

func handleStringOnly(input string) error {
	if input == "despacito" || input == "ree" {
		for x := 0; x != 10000; x++ {
			color.Red("Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree")       // Log string
			color.Yellow("Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree")    // Log string
			color.HiMagenta("Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree") // Log string
			color.Green("Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree")     // Log string
			color.Cyan("Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree")      // Log string
			color.Blue("Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree Ree")      // Log string
		}

		return nil // No error occurred, return nil
	}

	return errors.New("not string call")
}
