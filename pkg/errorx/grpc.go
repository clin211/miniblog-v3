// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package errorx

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCError 将 errorx 错误转换为 gRPC 错误
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// 检查是否是 errorx 错误
	if e, ok := err.(*Errno); ok {
		return status.Error(getGRPCCode(e.Code), e.Message)
	}

	// 如果不是 errorx 错误，返回内部错误
	return status.Error(codes.Internal, err.Error())
}

// getGRPCCode 根据 errorx 错误码返回对应的 gRPC 错误码
func getGRPCCode(errorxCode int) codes.Code {
	switch errorxCode {
	case 0:
		return codes.OK
	case 400001, 400002: // ErrBind, ErrInvalidParameter
		return codes.InvalidArgument
	case 401001, 401002, 401003, 401103: // ErrSignToken, ErrTokenInvalid, ErrUnauthorized, ErrPasswordIncorrect
		return codes.Unauthenticated
	case 403104: // ErrUserDisabled
		return codes.PermissionDenied
	case 404001, 404102: // ErrResourceNotFound, ErrUserNotFound
		return codes.NotFound
	case 409101: // ErrUserAlreadyExists
		return codes.AlreadyExists
	case 500001: // InternalServerError
		return codes.Internal
	default:
		// 根据 HTTP 状态码推断 gRPC 错误码
		httpCode := errorxCode / 1000
		switch httpCode {
		case 400:
			return codes.InvalidArgument
		case 401:
			return codes.Unauthenticated
		case 403:
			return codes.PermissionDenied
		case 404:
			return codes.NotFound
		case 409:
			return codes.AlreadyExists
		case 500:
			return codes.Internal
		default:
			return codes.Unknown
		}
	}
}
