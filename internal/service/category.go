package service

import (
	"context"
	"errors"
	"io"

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

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {

	cat, err := c.CategoryDB.Create(in.Name, in.Description)

	if !errors.Is(err, nil) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	catRes := &pb.Category{
		Id:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
	}

	return catRes, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {

	cats, err := c.CategoryDB.FindAll()

	if !errors.Is(err, nil) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	catList := []*pb.Category{}

	for _, cat := range *cats {
		catList = append(catList, &pb.Category{
			Id:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
		})
	}

	return &pb.CategoryList{
		Categories: catList,
	}, nil
}

func (c *CategoryService) GetCategoryById(ctx context.Context, in *pb.GetCategoryId) (*pb.Category, error) {

	cat, err := c.CategoryDB.FindById(in.Id)

	if !errors.Is(err, nil) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Category{
		Id:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
	}, nil

}

func (c *CategoryService) CreateCategoryStream(str pb.CategoryService_CreateCategoryStreamServer) error {

	var cats []*pb.Category

	for {

		cat, err := str.Recv()

		if err == io.EOF {
			return str.SendAndClose(&pb.CategoryList{
				Categories: cats,
			})
		}

		if !errors.Is(err, nil) {
			return status.Error(codes.Internal, err.Error())
		}

		categoryRes, err := c.CategoryDB.Create(cat.Name, cat.Description)

		if !errors.Is(err, nil) {
			return status.Error(codes.Internal, err.Error())
		}

		cats = append(cats, &pb.Category{
			Id:          categoryRes.ID,
			Name:        categoryRes.Name,
			Description: categoryRes.Description,
		})

	}

}

func (c *CategoryService) CreateCategoryBiStream(str pb.CategoryService_CreateCategoryBiStreamServer) error {

	for {

		cat, err := str.Recv()

		if err == io.EOF {
			return nil
		}

		if !errors.Is(err, nil) {
			return status.Error(codes.Internal, err.Error())
		}

		category, err := c.CategoryDB.Create(cat.Name, cat.Description)

		if !errors.Is(err, nil) {
			return status.Error(codes.Internal, err.Error())
		}

		err = str.Send(&pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})

		if !errors.Is(err, nil) {
			return status.Error(codes.Internal, err.Error())
		}

	}

}
