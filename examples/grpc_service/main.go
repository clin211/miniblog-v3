// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/clin211/miniblog-v3/examples/grpc_service/pb"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

// GetUser 实现单个用户查询
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := req.Id

	// 模拟业务逻辑
	if id == "0" {
		return nil, status.Errorf(codes.InvalidArgument, "用户ID非法: %s", id)
	}

	if id == "404" {
		return nil, status.Errorf(codes.NotFound, "用户不存在: %s", id)
	}

	// 模拟成功返回
	return &pb.GetUserResponse{
		Id:    id,
		Name:  fmt.Sprintf("用户-%s", id),
		Email: fmt.Sprintf("user%s@example.com", id),
	}, nil
}

// BatchGetUsers 实现批量用户查询
func (s *server) BatchGetUsers(ctx context.Context, req *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
	var results []*pb.UserResult

	for _, id := range req.Ids {
		result := &pb.UserResult{Id: id}

		// 尝试获取用户
		user, err := s.GetUser(ctx, &pb.GetUserRequest{Id: id})
		if err != nil {
			// 将错误转换为字符串
			result.Result = &pb.UserResult_Error{Error: err.Error()}
		} else {
			result.Result = &pb.UserResult_User{User: user}
		}

		results = append(results, result)
	}

	return &pb.BatchGetUsersResponse{Results: results}, nil
}

func main() {
	// 监听端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	s := grpc.NewServer()

	// 注册服务
	pb.RegisterUserServiceServer(s, &server{})

	// 启用反射服务（grpcurl 需要）
	reflection.Register(s)

	log.Printf("gRPC 服务器启动在端口 :50051")

	// 启动服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
