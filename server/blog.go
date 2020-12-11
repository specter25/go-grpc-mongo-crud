package server

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-hclog"
	"github.com/specter25/go-grpc-mongo-crud/models"
	blogpb "github.com/specter25/go-grpc-mongo-crud/protos/blog"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogServiceServer struct {
	log    hclog.Logger
	db     *mongo.Client
	blogdb *mongo.Collection
}

var mongoctx = context.Background()

func NewBlogServiceServer(log hclog.Logger, db *mongo.Client, blogdb *mongo.Collection) *BlogServiceServer {
	c := &BlogServiceServer{log, db, blogdb}
	return c
}
func (b *BlogServiceServer) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	blog := req.GetBlog()

	data := models.BlogItem{
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	result, err := b.blogdb.InsertOne(mongoctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error :%v", err),
		)
	}
	//convert object id of the created mongo db object to string to make it compatible with grpc
	oid := result.InsertedID.(primitive.ObjectID)
	blog.Id = oid.Hex()
	return &blogpb.CreateBlogRes{Blog: blog}, nil

}
func (b *BlogServiceServer) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {

	// convert Id from string to objectId
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("could not convert to objectID %s:%v", req.GetId(), err))
	}
	result := b.blogdb.FindOne(ctx, bson.M{"_id": oid})
	data := models.BlogItem{}

	if err = result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("could not find the blog with this object Id %s:%v", req.GetId(), err))
	}

	response := &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: data.AuthorId,
			Title:    data.Title,
			Content:  data.Content,
		},
	}
	return response, nil

}
func (b *BlogServiceServer) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {

	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("could not convert to objectID %v", err))
	}
	update := bson.M{
		"author_id": blog.GetAuthorId(),
		"title":     blog.GetTitle(),
		"content":   blog.GetContent(),
	}
	filter := bson.M{"_id": oid}
	result := b.blogdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	decoded := models.BlogItem{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("could not find the blog with this object Id : %v", err))
	}
	return &blogpb.UpdateBlogRes{
		Blog: &blogpb.Blog{
			Id:       decoded.ID.Hex(),
			AuthorId: decoded.AuthorId,
			Title:    decoded.Title,
			Content:  decoded.Content,
		},
	}, nil

}
func (b *BlogServiceServer) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("could not convert to objectID %s:%v", req.GetId(), err))
	}

	_, err = b.blogdb.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("could not find the blog to delete %s:%v", req.GetId(), err))
	}
	return &blogpb.DeleteBlogRes{
		Success: true,
	}, nil

}
func (s *BlogServiceServer) ListBlog(req *blogpb.ListBlogReq, stream blogpb.BlogService_ListBlogServer) error {

	return nil

}
