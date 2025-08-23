package errorx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestFromGRPCError(t *testing.T) {
	tests := []struct {
		name         string
		grpcError    error
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "nil error",
			grpcError:    nil,
			expectedCode: 0,
			expectedMsg:  "",
		},
		{
			name:         "InvalidArgument error",
			grpcError:    status.Error(codes.InvalidArgument, "参数错误"),
			expectedCode: ErrInvalidParameter.Code,
			expectedMsg:  "参数错误",
		},
		{
			name:         "Unauthenticated error",
			grpcError:    status.Error(codes.Unauthenticated, "认证失败"),
			expectedCode: ErrUnauthorized.Code,
			expectedMsg:  "认证失败",
		},
		{
			name:         "PermissionDenied error",
			grpcError:    status.Error(codes.PermissionDenied, "权限不足"),
			expectedCode: ErrUserDisabled.Code,
			expectedMsg:  "权限不足",
		},
		{
			name:         "NotFound error",
			grpcError:    status.Error(codes.NotFound, "用户不存在"),
			expectedCode: ErrUserNotFound.Code,
			expectedMsg:  "用户不存在",
		},
		{
			name:         "AlreadyExists error",
			grpcError:    status.Error(codes.AlreadyExists, "手机号已存在"),
			expectedCode: ErrUserAlreadyExists.Code,
			expectedMsg:  "手机号已存在",
		},
		{
			name:         "Internal error",
			grpcError:    status.Error(codes.Internal, "内部错误"),
			expectedCode: InternalServerError.Code,
			expectedMsg:  "内部错误",
		},
		{
			name:         "Unknown error",
			grpcError:    status.Error(codes.Unknown, "未知错误"),
			expectedCode: InternalServerError.Code,
			expectedMsg:  "未知错误",
		},
		{
			name:         "non-gRPC error",
			grpcError:    assert.AnError,
			expectedCode: InternalServerError.Code,
			expectedMsg:  assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromGRPCError(tt.grpcError)

			if tt.grpcError == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)

			if e, ok := result.(*Errno); ok {
				assert.Equal(t, tt.expectedCode, e.Code)
				assert.Equal(t, tt.expectedMsg, e.Message)
			} else {
				t.Fatal("Expected *Errno type")
			}
		})
	}
}
