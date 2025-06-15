package server

import (
	"context"
	"strings"

	"hookah-service/internal/model"
	pb "hookah-service/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TobaccoServer struct {
	pb.UnimplementedTobaccoSearchServiceServer
	index *model.SearchIndex
}

func NewTobaccoServer(index *model.SearchIndex) *TobaccoServer {
	return &TobaccoServer{
		index: index,
	}
}

func (s *TobaccoServer) BrandTobacco(
	ctx context.Context,
	req *pb.BrandTobaccoSearch,
) (*pb.TobaccoSearchResponse, error) {
	prefix := strings.ToLower(strings.TrimSpace(req.BrandPrefix))
	if len(prefix) < 2 {
		return nil, status.Error(codes.InvalidArgument, "Минимальная длина префикса - 2 символа")
	}

	var results []string
	for p, brands := range s.index.BrandPrefixes {
		if strings.HasPrefix(p, prefix) {
			results = append(results, brands...)
		}
	}

	return &pb.TobaccoSearchResponse{
		Results: removeDuplicates(results),
	}, nil
}

func (s *TobaccoServer) NameTobacco(
	ctx context.Context,
	req *pb.NameTobaccoSearch,
) (*pb.TobaccoSearchResponse, error) {
	brand := strings.TrimSpace(req.Brand)
	namePrefix := strings.ToLower(strings.TrimSpace(req.NamePrefix))
	if len(namePrefix) < 2 {
		return nil, status.Error(codes.InvalidArgument, "Минимальная длина префикса - 2 символа")
	}

	var results []string
	if namesMap, exists := s.index.BrandToNames[brand]; exists {
		for prefix, tobaccos := range namesMap {
			if strings.HasPrefix(prefix, namePrefix) {
				for _, t := range tobaccos {
					results = append(results, t.FlavorName)
				}
			}
		}
	}

	return &pb.TobaccoSearchResponse{Results: results}, nil
}

func removeDuplicates(items []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(items))
	
	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}