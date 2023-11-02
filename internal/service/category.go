package service

import (
	"context"
	"errors"

	"github.com/IcaroSilvaFK/grpc-go/internal/database"
	"github.com/IcaroSilvaFK/grpc-go/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB *database.Category
}

func NewCategoryService(categoryDB *database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {

	cat, err := c.CategoryDB.Create(in.Name, in.Description)

	if !errors.Is(err, nil) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	catRes := &pb.Category{
		Id:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
	}

	return &pb.CategoryResponse{
		Category: catRes,
	}, nil
}
