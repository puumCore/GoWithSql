package main

import (
	"database/sql"
	"fmt"
	"go_with_sql/iam/auth"
	"go_with_sql/repository"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

func gRPCServer() {
	listen, err := net.Listen("tcp", ":9090")
	checkError(err)

	grpcServer := grpc.NewServer()
	auth.RegisterDefaultAuthenticationServiceServer(grpcServer, &server{})
	log.Println("Server listening on port :9090")
	if err := grpcServer.Serve(listen); err != nil {
		checkError(err)
	}
}

type server struct {
	auth.UnimplementedDefaultAuthenticationServiceServer
}

func (s *server) RequestPasswordReset(ctx context.Context, req *auth.UserReq) (*auth.StandardResponse, error) {
	// Implementation
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *server) Authenticate(ctx context.Context, req *auth.AuthReq) (*auth.StandardResponse, error) {
	fmt.Printf(">> AuthReq > %s\n", req)

	connection, err := repository.GetDbConnection()
	checkError(err)
	defer func(connection *sql.DB) {
		err := connection.Close()
		repository.CheckError(err)
	}(connection)

	user := repository.GetUserByUsername(connection, req.Username)
	validated := repository.ValidatePassword(user.Password, req.Password)
	if validated {
		return &auth.StandardResponse{Code: 200, Message: "<<Login successful"}, nil
	} else {
		return &auth.StandardResponse{Code: 400, Message: "<<Invalid Credentials"}, nil
	}
}
