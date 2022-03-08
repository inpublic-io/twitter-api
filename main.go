package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/inpublic-io/twitter-api/clients"
	"github.com/inpublic-io/twitter-api/grpc"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/net/context"
)

func main() {
	// subscribe for os interrupt signal (ctrl+c)
	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)

	// initalizes twitter client
	err := clients.InitializeTwitterClient()
	if err != nil {
		log.Fatalf("failed to initialize twitter client: %v", err)
	}

	// initalizes influxdb client
	clients.InitializeInfluxClient()

	// creates an grpc server instance
	grpcServer := grpc.NewServer()

	grpcWebPort, ok := os.LookupEnv("GRPC_WEB_SERVER_PORT")
	if !ok {
		grpcWebPort = "8080"
	}

	corsOrigin, ok := os.LookupEnv("GRPC_WEB_SERVER_CORS_ORIGIN")
	if !ok {
		corsOrigin = "*"
	}

	// creates an grpc-web server with the wrapped grpc server
	grpcWebServer := &http.Server{
		Addr: ":" + grpcWebPort,
		Handler: &grpcWebHandler{
			corsOrigin:  corsOrigin,
			wrappedGrpc: grpcweb.WrapServer(grpcServer),
		},
	}

	go func() {
		fmt.Printf("grpc web server listening at %v\n", grpcWebServer.Addr)
		err = grpcWebServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server unexpectedly closed: %+v\n", err)
		}

		fmt.Printf("failed to serve grpc-web: %+v\n", err)
	}()
	defer func() {
		// tries to gracefully stop grpc-web server
		fmt.Printf("trying to gracefully shutdown grpc-web server\n")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		err = grpcWebServer.Shutdown(ctx)
		if err != nil {
			log.Fatalf("test %+v\n", err)
		}
	}()

	grpcServerPort, ok := os.LookupEnv("GRPC_SERVER_PORT")
	if !ok {
		grpcServerPort = "50051"
	}

	listener, err := net.Listen("tcp", ":"+grpcServerPort)
	if err != nil {
		log.Fatalf("failed to bind to port %s: %v", grpcServerPort, err)
	}

	go func() {
		fmt.Printf("grpc server listening at %v\n", listener.Addr())
		if err := grpcServer.Serve(listener); err != nil {
			fmt.Printf("failed to serve grpc server: %+v\n", err)
		}
	}()
	defer func() {
		// tries to gracefully stop grpc server
		fmt.Printf("trying to gracefully shutdown grpc server\n")
		grpcServer.GracefulStop()
	}()

	<-interruptSignal
	fmt.Printf("stopping server...\n")
}

type grpcWebHandler struct {
	corsOrigin  string
	wrappedGrpc *grpcweb.WrappedGrpcServer
}

func (handler *grpcWebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", handler.corsOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,x-grpc-web")
	handler.wrappedGrpc.ServeHTTP(w, r)
}
