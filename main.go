package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dowlandaiello/GoP2P/cli"
	"github.com/dowlandaiello/GoP2P/common"
	commonServer "github.com/dowlandaiello/GoP2P/internal/rpc/common"
	"github.com/dowlandaiello/GoP2P/internal/rpc/database"
	"github.com/dowlandaiello/GoP2P/internal/rpc/environment"
	handlerServer "github.com/dowlandaiello/GoP2P/internal/rpc/handler"
	nodeServer "github.com/dowlandaiello/GoP2P/internal/rpc/node"
	commonProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/common"
	databaseProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/database"
	environmentProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/environment"
	handlerProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/handler"
	nodeProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/node"
	protoProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/protobuf"
	shardProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/shard"
	upnpProto "github.com/dowlandaiello/GoP2P/internal/rpc/proto/upnp"
	protoServer "github.com/dowlandaiello/GoP2P/internal/rpc/protobuf"
	shardServer "github.com/dowlandaiello/GoP2P/internal/rpc/shard"
	upnpServer "github.com/dowlandaiello/GoP2P/internal/rpc/upnp"
	"github.com/dowlandaiello/GoP2P/types/handler"
	"github.com/dowlandaiello/GoP2P/types/node"
	"github.com/dowlandaiello/GoP2P/upnp"
	"github.com/fatih/color"
)

var (
	terminalFlag   = flag.Bool("terminal", false, "launch GoP2P in terminal mode")                                                                                    // Init term flag
	upnpFlag       = flag.Bool("no-upnp", false, "launch GoP2P without automatic UPnP port forwarding")                                                               // Init upnp flag
	rpcPortFlag    = flag.Int("rpc-port", 8080, "launch GoP2P with specified RPC port")                                                                               // Init RPC port flag
	noColorFlag    = flag.Bool("no-color", false, "disables GoP2P terminal colored output")                                                                           // Init color flag
	forwardRPCFlag = flag.Bool("forward-rpc", false, "enables forwarding of GoP2P RPC terminal ports")                                                                // Init forward RPC flag
	rpcAddrFlag    = flag.String("rpc-address", fmt.Sprintf("localhost:%s", strconv.Itoa(*rpcPortFlag)), "connects to remote RPC terminal (default: localhost:8080)") // Init remote rpc addr flag
	silentMode     = flag.Bool("s", false, "launches gop2p in silent mode (silences prints)")                                                                         // Init silent flag
)

func main() {
	flag.Parse() // Parse flags

	common.Silent = *silentMode // Set silent

	if !*upnpFlag { // Check for UPnP
		if *forwardRPCFlag {
			go upnp.ForwardPortSilent(uint(*rpcPortFlag)) // Forward RPC port
		}

		go upnp.ForwardPortSilent(3000) // Forward port 3000
	}

	if *noColorFlag { // Check for no colors
		color.NoColor = true // Disable colors
	}

	if strings.Contains(*rpcAddrFlag, "localhost") { // Check for default RPC address
		startRPCServer() // Start RPC server
	}

	if *terminalFlag { // Check for terminal
		*rpcAddrFlag = strings.Split(*rpcAddrFlag, ":")[0] // Remove port

		cli.NewTerminal(uint(*rpcPortFlag), *rpcAddrFlag) // Initialize terminal
	}

	startNode() // Attempt to start GoP2P in node mode
}

// startRPCServer - start RPC server
func startRPCServer() {
	err := common.GenerateTLSCertificates("gop2pTerm") // Generate certs

	if err != nil { // Check for errors
		panic(err) // Panic
	}

	nodeHandler := nodeProto.NewNodeServer(&nodeServer.Server{}, nil)                       // Init handler
	handlerHandler := handlerProto.NewHandlerServer(&handlerServer.Server{}, nil)           // Init handler
	environmentHandler := environmentProto.NewEnvironmentServer(&environment.Server{}, nil) // Init handler
	upnpHandler := upnpProto.NewUpnpServer(&upnpServer.Server{}, nil)                       // Init handler
	databaseHandler := databaseProto.NewDatabaseServer(&database.Server{}, nil)             // Init handler
	commonHandler := commonProto.NewCommonServer(&commonServer.Server{}, nil)               // Init handler
	shardHandler := shardProto.NewShardServer(&shardServer.Server{}, nil)                   // Init handler
	protoHandler := protoProto.NewProtoServer(&protoServer.Server{}, nil)                   // Init handler

	mux := http.NewServeMux() // Init mux

	mux.Handle(nodeProto.NodePathPrefix, nodeHandler)                      // Start mux node handler
	mux.Handle(handlerProto.HandlerPathPrefix, handlerHandler)             // Start mux handler handler
	mux.Handle(environmentProto.EnvironmentPathPrefix, environmentHandler) // Start mux environment handler
	mux.Handle(upnpProto.UpnpPathPrefix, upnpHandler)                      // Start mux upnp handler
	mux.Handle(databaseProto.DatabasePathPrefix, databaseHandler)          // Start mux database handler
	mux.Handle(commonProto.CommonPathPrefix, commonHandler)                // Start mux common handler
	mux.Handle(shardProto.ShardPathPrefix, shardHandler)                   // Start mux shard handler
	mux.Handle(protoProto.ProtoPathPrefix, protoHandler)                   // Start mux proto handler

	go http.ListenAndServeTLS(":"+strconv.Itoa(*rpcPortFlag), "gop2pTermCert.pem", "gop2pTermKey.pem", mux) // Start server
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
- Fix readme (or lack thereof)
- Add -v flag (silence common.Println)
- Fix protonet unit tests
*/
