package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	blogpb "github.com/specter25/go-grpc-mongo-crud/protos/blogs"
)

type BlogServiceServer struct {
	log hclog.Logger
}

func NewBlogServiceServer(log hclog.Logger) *BlogServiceServer {
	c := &BlogServiceServer{log}
	return c
}
func (b *BlogServiceServer) CreateBlog(ctx context.Context, cbr *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {
	return nil, nil

}
func (b *BlogServiceServer) ReadBlog(ctx context.Context, cbr *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {

	return nil, nil

}
func (b *BlogServiceServer) UpdateBlog(ctx context.Context, cbr *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {

	return nil, nil

}
func (b *BlogServiceServer) DeleteBlog(ctx context.Context, cbr *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {

	return nil, nil

}
func (s *BlogServiceServer) ListBlog(req *blogpb.ListBlogReq, stream blogpb.BlogService_ListBlogServer) error {

	return nil

}
