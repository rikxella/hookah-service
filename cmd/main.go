package main

import (
	"log"
	"net"
	"path/filepath"

	"hookah-service/internal/model"
	pb "hookah-service/internal/proto"
	"hookah-service/internal/server"
	"google.golang.org/grpc"
)

func main() {
	tobaccos, err := model.LoadFromCSV(filepath.Join("data", "tobaccos.csv"))
	if err != nil {
		log.Fatalf("Ошибка загрузки данных: %v", err)
	}

	index := model.BuildIndex(tobaccos)

	grpcServer := grpc.NewServer()
	pb.RegisterTobaccoSearchServiceServer(grpcServer, server.NewTobaccoServer(index))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	log.Println("Сервер запущен на порту :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка работы сервера: %v", err)
	}
}