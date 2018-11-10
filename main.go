package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/fatih/color"
	"github.com/mitsukomegumi/GoP2P/cli"
	"github.com/mitsukomegumi/GoP2P/common"
	commonServer "github.com/mitsukomegumi/GoP2P/rpc/common"
	"github.com/mitsukomegumi/GoP2P/rpc/database"
	"github.com/mitsukomegumi/GoP2P/rpc/environment"
	handlerServer "github.com/mitsukomegumi/GoP2P/rpc/handler"
	nodeServer "github.com/mitsukomegumi/GoP2P/rpc/node"
	commonProto "github.com/mitsukomegumi/GoP2P/rpc/proto/common"
	databaseProto "github.com/mitsukomegumi/GoP2P/rpc/proto/database"
	environmentProto "github.com/mitsukomegumi/GoP2P/rpc/proto/environment"
	handlerProto "github.com/mitsukomegumi/GoP2P/rpc/proto/handler"
	nodeProto "github.com/mitsukomegumi/GoP2P/rpc/proto/node"
	upnpProto "github.com/mitsukomegumi/GoP2P/rpc/proto/upnp"
	upnpServer "github.com/mitsukomegumi/GoP2P/rpc/upnp"
	"github.com/mitsukomegumi/GoP2P/types/handler"
	"github.com/mitsukomegumi/GoP2P/types/node"
	"github.com/mitsukomegumi/GoP2P/upnp"
)

var (
	terminalFlag = flag.Bool("terminal", false, "launch GoP2P in terminal mode")                      // Init term flag
	upnpFlag     = flag.Bool("no-upnp", false, "launch GoP2P without automatic UPnP port forwarding") // Init upnp flag
	rpcPortFlag  = flag.Int("rpc-port", 8080, "launch GoP2P with specified RPC port")                 // Init RPC port flag
	noColorFlag  = flag.Bool("no-color", false, "disables GoP2P terminal colored output")             // Init color flag
)

func main() {
	flag.Parse() // Parse flags

	if !*upnpFlag { // Check for UPnP
		go upnp.ForwardPortSilent(uint(*rpcPortFlag)) // Forward RPC port
		go upnp.ForwardPortSilent(3000)               // Forward port 3000
	} else if *noColorFlag { // Check for no colors
		color.NoColor = true // Disable colors
	}

	startRPCServer() // Start RPC server

	if *terminalFlag {
		cli.NewTerminal(uint(*rpcPortFlag)) // Initialize terminal
	}

	startNode() // Attempt to start GoP2P in node mode

	go common.Forever() // Prevent main from closing
	select {}           // Prevent main from closing
}

// startRPCServer - start RPC server
func startRPCServer() {
	nodeHandler := nodeProto.NewNodeServer(&nodeServer.Server{}, nil)                       // Init handler
	handlerHandler := handlerProto.NewHandlerServer(&handlerServer.Server{}, nil)           // Init handler
	environmentHandler := environmentProto.NewEnvironmentServer(&environment.Server{}, nil) // Init handler
	upnpHandler := upnpProto.NewUpnpServer(&upnpServer.Server{}, nil)                       // Init handler
	databaseHandler := databaseProto.NewDatabaseServer(&database.Server{}, nil)             // Init handler
	commonHandler := commonProto.NewCommonServer(&commonServer.Server{}, nil)               // Init handler

	mux := http.NewServeMux() // Init mux

	mux.Handle(nodeProto.NodePathPrefix, nodeHandler)                      // Start mux node handler
	mux.Handle(handlerProto.HandlerPathPrefix, handlerHandler)             // Start mux handler handler
	mux.Handle(environmentProto.EnvironmentPathPrefix, environmentHandler) // Start mux environment handler
	mux.Handle(upnpProto.UpnpPathPrefix, upnpHandler)                      // Start mux upnp handler
	mux.Handle(databaseProto.DatabasePathPrefix, databaseHandler)          // Start mux database handler
	mux.Handle(commonProto.CommonPathPrefix, commonHandler)                // Start mux common handler

	go http.ListenAndServe(":"+strconv.Itoa(*rpcPortFlag), mux) // Start server
}

// startNode - attempt to execute attachnode, starthandler commands
func startNode() {
	currentDir, err := common.GetCurrentDir() // Fetch working directory

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	node, err := node.ReadNodeFromMemory(currentDir) // Read node from current dir

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	ln, err := node.StartListener(3000) // Start listener

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	err = handler.StartHandler(node, ln) // Start handler

	if err != nil { // Check for errors
		panic(err) // Panic
	}
}

/* TODO:
 */
