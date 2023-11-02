package main

import (
	"database/sql"
	"errors"
	"log"
	"net"

	"github.com/IcaroSilvaFK/grpc-go/internal/database"
	"github.com/IcaroSilvaFK/grpc-go/internal/pb"
	"github.com/IcaroSilvaFK/grpc-go/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")

	if !errors.Is(err, nil) {
		panic(err)
	}

	catDb := database.NewCategory(db)

	catService := service.NewCategoryService(catDb)

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	pb.RegisterCategoryServiceServer(grpcServer, catService)

	lis, err := net.Listen("tcp", ":50051")

	if !errors.Is(err, nil) {
		panic(err)
	}

	log.Println("Server running on port 50051")

	if err = grpcServer.Serve(lis); !errors.Is(err, nil) {
		panic(err)
	}

}
