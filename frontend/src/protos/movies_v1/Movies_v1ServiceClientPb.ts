/**
 * @fileoverview gRPC-Web generated client stub for movies_v1
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v4.25.3
// source: protos/movies_v1/movies_v1.proto


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as protos_movies_v1_movies_v1_pb from '../../protos/movies_v1/movies_v1_pb'; // proto import: "protos/movies_v1/movies_v1.proto"


export class MoviesV1Client {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname.replace(/\/+$/, '');
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodDescriptorGetMovies = new grpcWeb.MethodDescriptor(
    '/movies_v1.MoviesV1/GetMovies',
    grpcWeb.MethodType.UNARY,
    protos_movies_v1_movies_v1_pb.GetMoviesRequest,
    protos_movies_v1_movies_v1_pb.GetMoviesResponse,
    (request: protos_movies_v1_movies_v1_pb.GetMoviesRequest) => {
      return request.serializeBinary();
    },
    protos_movies_v1_movies_v1_pb.GetMoviesResponse.deserializeBinary
  );

  getMovies(
    request: protos_movies_v1_movies_v1_pb.GetMoviesRequest,
    metadata?: grpcWeb.Metadata | null): Promise<protos_movies_v1_movies_v1_pb.GetMoviesResponse>;

  getMovies(
    request: protos_movies_v1_movies_v1_pb.GetMoviesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.GetMoviesResponse) => void): grpcWeb.ClientReadableStream<protos_movies_v1_movies_v1_pb.GetMoviesResponse>;

  getMovies(
    request: protos_movies_v1_movies_v1_pb.GetMoviesRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.GetMoviesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/movies_v1.MoviesV1/GetMovies',
        request,
        metadata || {},
        this.methodDescriptorGetMovies,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/movies_v1.MoviesV1/GetMovies',
    request,
    metadata || {},
    this.methodDescriptorGetMovies);
  }

  methodDescriptorGetMovieById = new grpcWeb.MethodDescriptor(
    '/movies_v1.MoviesV1/GetMovieById',
    grpcWeb.MethodType.UNARY,
    protos_movies_v1_movies_v1_pb.GetMoviesByIdRequest,
    protos_movies_v1_movies_v1_pb.GetMoviesByIdResponse,
    (request: protos_movies_v1_movies_v1_pb.GetMoviesByIdRequest) => {
      return request.serializeBinary();
    },
    protos_movies_v1_movies_v1_pb.GetMoviesByIdResponse.deserializeBinary
  );

  getMovieById(
    request: protos_movies_v1_movies_v1_pb.GetMoviesByIdRequest,
    metadata?: grpcWeb.Metadata | null): Promise<protos_movies_v1_movies_v1_pb.GetMoviesByIdResponse>;

  getMovieById(
    request: protos_movies_v1_movies_v1_pb.GetMoviesByIdRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.GetMoviesByIdResponse) => void): grpcWeb.ClientReadableStream<protos_movies_v1_movies_v1_pb.GetMoviesByIdResponse>;

  getMovieById(
    request: protos_movies_v1_movies_v1_pb.GetMoviesByIdRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.GetMoviesByIdResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/movies_v1.MoviesV1/GetMovieById',
        request,
        metadata || {},
        this.methodDescriptorGetMovieById,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/movies_v1.MoviesV1/GetMovieById',
    request,
    metadata || {},
    this.methodDescriptorGetMovieById);
  }

  methodDescriptorGetMoviesByFilter = new grpcWeb.MethodDescriptor(
    '/movies_v1.MoviesV1/GetMoviesByFilter',
    grpcWeb.MethodType.UNARY,
    protos_movies_v1_movies_v1_pb.GetMoviesByFilterRequest,
    protos_movies_v1_movies_v1_pb.GetMoviesResponse,
    (request: protos_movies_v1_movies_v1_pb.GetMoviesByFilterRequest) => {
      return request.serializeBinary();
    },
    protos_movies_v1_movies_v1_pb.GetMoviesResponse.deserializeBinary
  );

  getMoviesByFilter(
    request: protos_movies_v1_movies_v1_pb.GetMoviesByFilterRequest,
    metadata?: grpcWeb.Metadata | null): Promise<protos_movies_v1_movies_v1_pb.GetMoviesResponse>;

  getMoviesByFilter(
    request: protos_movies_v1_movies_v1_pb.GetMoviesByFilterRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.GetMoviesResponse) => void): grpcWeb.ClientReadableStream<protos_movies_v1_movies_v1_pb.GetMoviesResponse>;

  getMoviesByFilter(
    request: protos_movies_v1_movies_v1_pb.GetMoviesByFilterRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.GetMoviesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/movies_v1.MoviesV1/GetMoviesByFilter',
        request,
        metadata || {},
        this.methodDescriptorGetMoviesByFilter,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/movies_v1.MoviesV1/GetMoviesByFilter',
    request,
    metadata || {},
    this.methodDescriptorGetMoviesByFilter);
  }

  methodDescriptorAddMovies = new grpcWeb.MethodDescriptor(
    '/movies_v1.MoviesV1/AddMovies',
    grpcWeb.MethodType.UNARY,
    protos_movies_v1_movies_v1_pb.AddMoviesRequest,
    protos_movies_v1_movies_v1_pb.AddMoviesResponse,
    (request: protos_movies_v1_movies_v1_pb.AddMoviesRequest) => {
      return request.serializeBinary();
    },
    protos_movies_v1_movies_v1_pb.AddMoviesResponse.deserializeBinary
  );

  addMovies(
    request: protos_movies_v1_movies_v1_pb.AddMoviesRequest,
    metadata?: grpcWeb.Metadata | null): Promise<protos_movies_v1_movies_v1_pb.AddMoviesResponse>;

  addMovies(
    request: protos_movies_v1_movies_v1_pb.AddMoviesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.AddMoviesResponse) => void): grpcWeb.ClientReadableStream<protos_movies_v1_movies_v1_pb.AddMoviesResponse>;

  addMovies(
    request: protos_movies_v1_movies_v1_pb.AddMoviesRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.AddMoviesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/movies_v1.MoviesV1/AddMovies',
        request,
        metadata || {},
        this.methodDescriptorAddMovies,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/movies_v1.MoviesV1/AddMovies',
    request,
    metadata || {},
    this.methodDescriptorAddMovies);
  }

  methodDescriptorDeleteMovies = new grpcWeb.MethodDescriptor(
    '/movies_v1.MoviesV1/DeleteMovies',
    grpcWeb.MethodType.UNARY,
    protos_movies_v1_movies_v1_pb.DeleteMoviesRequest,
    protos_movies_v1_movies_v1_pb.DeleteMoviesResponse,
    (request: protos_movies_v1_movies_v1_pb.DeleteMoviesRequest) => {
      return request.serializeBinary();
    },
    protos_movies_v1_movies_v1_pb.DeleteMoviesResponse.deserializeBinary
  );

  deleteMovies(
    request: protos_movies_v1_movies_v1_pb.DeleteMoviesRequest,
    metadata?: grpcWeb.Metadata | null): Promise<protos_movies_v1_movies_v1_pb.DeleteMoviesResponse>;

  deleteMovies(
    request: protos_movies_v1_movies_v1_pb.DeleteMoviesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.DeleteMoviesResponse) => void): grpcWeb.ClientReadableStream<protos_movies_v1_movies_v1_pb.DeleteMoviesResponse>;

  deleteMovies(
    request: protos_movies_v1_movies_v1_pb.DeleteMoviesRequest,
    metadata?: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: protos_movies_v1_movies_v1_pb.DeleteMoviesResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/movies_v1.MoviesV1/DeleteMovies',
        request,
        metadata || {},
        this.methodDescriptorDeleteMovies,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/movies_v1.MoviesV1/DeleteMovies',
    request,
    metadata || {},
    this.methodDescriptorDeleteMovies);
  }

}
