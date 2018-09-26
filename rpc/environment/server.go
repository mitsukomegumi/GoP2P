package environment

import (
	"context"
	"fmt"

	"github.com/mitsukomegumi/GoP2P/common"
	environmentProto "github.com/mitsukomegumi/GoP2P/rpc/proto/environment"
	"github.com/mitsukomegumi/GoP2P/types/environment"
)

// Server - GoP2P RPC server
type Server struct{}

// NewEnvironment - environment.NewEnvironment RPC handler
func (server *Server) NewEnvironment(ctx context.Context, req *environmentProto.GeneralRequest) (*environmentProto.GeneralResponse, error) {
	currentDir, err := common.GetCurrentDir() // Fetch current directory

	if err != nil { // Check for errors
		return &environmentProto.GeneralResponse{}, err // Return found error
	}

	env, err := environment.ReadEnvironmentFromMemory(currentDir) // Attempt to read environment from memory

	if err != nil { // Check for errors
		env, err = environment.NewEnvironment() // Init environment

		if err != nil { // Check for errors
			return &environmentProto.GeneralResponse{}, err // Return found error
		}
	}

	env.WriteToMemory(currentDir) // Save for persistency

	return &environmentProto.GeneralResponse{Message: fmt.Sprintf("\nInitialized environment %v", env)}, nil // Return response
}

// QueryType - environment.QueryType RPC handler
func (server *Server) QueryType(ctx context.Context, req *environmentProto.GeneralRequest) (*environmentProto.GeneralResponse, error) {
	currentDir, err := common.GetCurrentDir() // Fetch current directory

	if err != nil { // Check for errors
		return &environmentProto.GeneralResponse{}, err // Return found error
	}

	env, err := environment.ReadEnvironmentFromMemory(currentDir) // Attempt to read environment from memory

	if err != nil { // Check for errors
		env, err = environment.NewEnvironment() // Init environment

		if err != nil { // Check for errors
			return &environmentProto.GeneralResponse{}, err // Return found error
		}

		env.WriteToMemory(currentDir) // Save for persistency
	}

	foundVariable, err := env.QueryType(req.VariableType) // Query type

	if err != nil { // Check for errors
		return &environmentProto.GeneralResponse{}, err // Return found error
	}

	return &environmentProto.GeneralResponse{Message: fmt.Sprintf("\nFound variable %v", foundVariable)}, nil // No error occurred, return output
}

// QueryValue - environment.QueryValue RPC handler
func (server *Server) QueryValue(ctx context.Context, req *environmentProto.GeneralRequest) (*environmentProto.GeneralResponse, error) {
	currentDir, err := common.GetCurrentDir() // Fetch current directory

	if err != nil { // Check for errors
		return &environmentProto.GeneralResponse{}, err // Return found error
	}

	env, err := environment.ReadEnvironmentFromMemory(currentDir) // Attempt to read environment from memory

	if err != nil { // Check for errors
		env, err = environment.NewEnvironment() // Init environment

		if err != nil { // Check for errors
			return &environmentProto.GeneralResponse{}, err // Return found error
		}

		env.WriteToMemory(currentDir) // Save for persistency
	}

	foundVariable, err := env.QueryValue(req.Value) // Query type

	if err != nil { // Check for errors
		return &environmentProto.GeneralResponse{}, err // Return found error
	}

	return &environmentProto.GeneralResponse{Message: fmt.Sprintf("\nFound variable %v", foundVariable)}, nil // No error occurred, return output
}
