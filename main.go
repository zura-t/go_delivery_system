package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/zura-t/go_delivery_system/cmd/gapi"
	"github.com/zura-t/go_delivery_system/internal"
	"github.com/zura-t/go_delivery_system/pb"
	_ "github.com/zura-t/simplebank/doc/statik"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := internal.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config file:", err)
	}

	runGatewayServer(config)
}

func runGatewayServer(config internal.Config) {
	server, err := gapi.NewServer(config)
	if err != nil {
		log.Fatalf("can't create server: %s", err)
	}

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterUsersServiceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalf("can't register handler server: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("cannot create statik fs: %s", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatalf("can't create listener: %s", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("can't start HTTP gateway server: %s", err)
	}
}
