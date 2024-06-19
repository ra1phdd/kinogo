// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: protos/movies_v1/movies_v1.proto

package movies_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MoviesV1Client is the client API for MoviesV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MoviesV1Client interface {
	GetMovies(ctx context.Context, in *GetMoviesRequest, opts ...grpc.CallOption) (*GetMoviesResponse, error)
	GetMovieById(ctx context.Context, in *GetMoviesByIdRequest, opts ...grpc.CallOption) (*GetMoviesByIdResponse, error)
	GetMoviesByFilter(ctx context.Context, in *GetMoviesByFilterRequest, opts ...grpc.CallOption) (*GetMoviesResponse, error)
	AddMovies(ctx context.Context, in *AddMoviesRequest, opts ...grpc.CallOption) (*AddMoviesResponse, error)
	DeleteMovies(ctx context.Context, in *DeleteMoviesRequest, opts ...grpc.CallOption) (*DeleteMoviesResponse, error)
}

type moviesV1Client struct {
	cc grpc.ClientConnInterface
}

func NewMoviesV1Client(cc grpc.ClientConnInterface) MoviesV1Client {
	return &moviesV1Client{cc}
}

func (c *moviesV1Client) GetMovies(ctx context.Context, in *GetMoviesRequest, opts ...grpc.CallOption) (*GetMoviesResponse, error) {
	out := new(GetMoviesResponse)
	err := c.cc.Invoke(ctx, "/movies_v1.MoviesV1/GetMovies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesV1Client) GetMovieById(ctx context.Context, in *GetMoviesByIdRequest, opts ...grpc.CallOption) (*GetMoviesByIdResponse, error) {
	out := new(GetMoviesByIdResponse)
	err := c.cc.Invoke(ctx, "/movies_v1.MoviesV1/GetMovieById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesV1Client) GetMoviesByFilter(ctx context.Context, in *GetMoviesByFilterRequest, opts ...grpc.CallOption) (*GetMoviesResponse, error) {
	out := new(GetMoviesResponse)
	err := c.cc.Invoke(ctx, "/movies_v1.MoviesV1/GetMoviesByFilter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesV1Client) AddMovies(ctx context.Context, in *AddMoviesRequest, opts ...grpc.CallOption) (*AddMoviesResponse, error) {
	out := new(AddMoviesResponse)
	err := c.cc.Invoke(ctx, "/movies_v1.MoviesV1/AddMovies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moviesV1Client) DeleteMovies(ctx context.Context, in *DeleteMoviesRequest, opts ...grpc.CallOption) (*DeleteMoviesResponse, error) {
	out := new(DeleteMoviesResponse)
	err := c.cc.Invoke(ctx, "/movies_v1.MoviesV1/DeleteMovies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MoviesV1Server is the server API for MoviesV1 service.
// All implementations must embed UnimplementedMoviesV1Server
// for forward compatibility
type MoviesV1Server interface {
	GetMovies(context.Context, *GetMoviesRequest) (*GetMoviesResponse, error)
	GetMovieById(context.Context, *GetMoviesByIdRequest) (*GetMoviesByIdResponse, error)
	GetMoviesByFilter(context.Context, *GetMoviesByFilterRequest) (*GetMoviesResponse, error)
	AddMovies(context.Context, *AddMoviesRequest) (*AddMoviesResponse, error)
	DeleteMovies(context.Context, *DeleteMoviesRequest) (*DeleteMoviesResponse, error)
	mustEmbedUnimplementedMoviesV1Server()
}

// UnimplementedMoviesV1Server must be embedded to have forward compatible implementations.
type UnimplementedMoviesV1Server struct {
}

func (UnimplementedMoviesV1Server) GetMovies(context.Context, *GetMoviesRequest) (*GetMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovies not implemented")
}
func (UnimplementedMoviesV1Server) GetMovieById(context.Context, *GetMoviesByIdRequest) (*GetMoviesByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovieById not implemented")
}
func (UnimplementedMoviesV1Server) GetMoviesByFilter(context.Context, *GetMoviesByFilterRequest) (*GetMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMoviesByFilter not implemented")
}
func (UnimplementedMoviesV1Server) AddMovies(context.Context, *AddMoviesRequest) (*AddMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMovies not implemented")
}
func (UnimplementedMoviesV1Server) DeleteMovies(context.Context, *DeleteMoviesRequest) (*DeleteMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMovies not implemented")
}
func (UnimplementedMoviesV1Server) mustEmbedUnimplementedMoviesV1Server() {}

// UnsafeMoviesV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MoviesV1Server will
// result in compilation errors.
type UnsafeMoviesV1Server interface {
	mustEmbedUnimplementedMoviesV1Server()
}

func RegisterMoviesV1Server(s grpc.ServiceRegistrar, srv MoviesV1Server) {
	s.RegisterService(&MoviesV1_ServiceDesc, srv)
}

func _MoviesV1_GetMovies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMoviesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesV1Server).GetMovies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_v1.MoviesV1/GetMovies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesV1Server).GetMovies(ctx, req.(*GetMoviesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesV1_GetMovieById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMoviesByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesV1Server).GetMovieById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_v1.MoviesV1/GetMovieById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesV1Server).GetMovieById(ctx, req.(*GetMoviesByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesV1_GetMoviesByFilter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMoviesByFilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesV1Server).GetMoviesByFilter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_v1.MoviesV1/GetMoviesByFilter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesV1Server).GetMoviesByFilter(ctx, req.(*GetMoviesByFilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesV1_AddMovies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMoviesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesV1Server).AddMovies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_v1.MoviesV1/AddMovies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesV1Server).AddMovies(ctx, req.(*AddMoviesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MoviesV1_DeleteMovies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMoviesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MoviesV1Server).DeleteMovies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movies_v1.MoviesV1/DeleteMovies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MoviesV1Server).DeleteMovies(ctx, req.(*DeleteMoviesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MoviesV1_ServiceDesc is the grpc.ServiceDesc for MoviesV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MoviesV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "movies_v1.MoviesV1",
	HandlerType: (*MoviesV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMovies",
			Handler:    _MoviesV1_GetMovies_Handler,
		},
		{
			MethodName: "GetMovieById",
			Handler:    _MoviesV1_GetMovieById_Handler,
		},
		{
			MethodName: "GetMoviesByFilter",
			Handler:    _MoviesV1_GetMoviesByFilter_Handler,
		},
		{
			MethodName: "AddMovies",
			Handler:    _MoviesV1_AddMovies_Handler,
		},
		{
			MethodName: "DeleteMovies",
			Handler:    _MoviesV1_DeleteMovies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/movies_v1/movies_v1.proto",
}