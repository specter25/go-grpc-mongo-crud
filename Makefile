.PHONY: protos

protos:
	protoc -I protos/ protos/blog.proto --go_out=plugins=grpc:protos/blogs