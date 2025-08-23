// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package errorx

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FromGRPCError 将 gRPC 错误转换为 errorx 错误
func FromGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// 检查是否是 gRPC 错误
	st, ok := status.FromError(err)
	if !ok {
		// 如果不是 gRPC 错误，返回内部错误
		return InternalServerError.SetMessage("%s", err.Error())
	}

	// 根据 gRPC 错误码转换为 errorx 错误
	switch st.Code() {
	case codes.OK:
		return nil
	case codes.InvalidArgument:
		return ErrInvalidParameter.SetMessage("%s", st.Message())
	case codes.Unauthenticated:
		return ErrUnauthorized.SetMessage("%s", st.Message())
	case codes.PermissionDenied:
		return ErrUserDisabled.SetMessage("%s", st.Message())
	case codes.NotFound:
		return ErrUserNotFound.SetMessage("%s", st.Message())
	case codes.AlreadyExists:
		return ErrUserAlreadyExists.SetMessage("%s", st.Message())
	case codes.Internal:
		return InternalServerError.SetMessage("%s", st.Message())
	default:
		// 其他错误码返回内部错误
		return InternalServerError.SetMessage("%s", st.Message())
	}
}
