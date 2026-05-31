package main

import (
	"context"
	"log"

	proto "github.com/Ternuraa/DistributedMicroservice/services/listing_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListingServer struct {
	proto.UnimplementedListingServiceServer
}

// ВАЖНО: Имя функции теперь GetListingInfo (строго как в сгенерированном .pb.go)
func (s *ListingServer) GetListingInfo(ctx context.Context, req *proto.ListingRequest) (*proto.ListingResponse, error) {
	log.Printf("🚀 [Listing Service] Получен gRPC запрос на жильё с ID: %s", req.Id)

	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ID не может быть пустым")
	}

	if req.Id == "404" {
		return nil, status.Errorf(codes.NotFound, "Жильё с таким ID не существует")
	}

	return &proto.ListingResponse{
		Id:          req.Id,
		Price:       1500.50,
		IsAvailable: true,
	}, nil
}
