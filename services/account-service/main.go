package main

import (
	"log"
	"net"

	acdb "github.com/RAF-SI-2025/EXBanka-4-Backend/services/account-service/db"
	"github.com/RAF-SI-2025/EXBanka-4-Backend/services/account-service/handlers"
	pb "github.com/RAF-SI-2025/EXBanka-4-Backend/shared/pb/account"
	"google.golang.org/grpc"
)

func main() {
	database, err := acdb.Connect("postgres://account_user:account_pass@localhost:5436/account_db?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to account_db: %v", err)
	}
	defer database.Close()

	clientDB, err := acdb.Connect("postgres://client_user:client_pass@localhost:5435/client_db?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to client_db: %v", err)
	}
	defer clientDB.Close()

	exchangeDB, err := acdb.Connect("postgres://exchange_user:exchange_pass@localhost:5438/exchange_db?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to exchange_db: %v", err)
	}
	defer exchangeDB.Close()

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen on :50054: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAccountServiceServer(grpcServer, &handlers.AccountServer{
		DB:         database,
		ClientDB:   clientDB,
		ExchangeDB: exchangeDB,
	})

	log.Println("account-service gRPC server listening on :50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
