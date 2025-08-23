// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package errorx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToGRPCError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode codes.Code
		expectedMsg  string
	}{
		{
			name:         "nil error",
			err:          nil,
			expectedCode: codes.OK,
			expectedMsg:  "",
		},
		{
			name:         "ErrInvalidParameter",
			err:          ErrInvalidParameter,
			expectedCode: codes.InvalidArgument,
			expectedMsg:  "Parameter verification failed.",
		},
		{
			name:         "ErrTokenInvalid",
			err:          ErrTokenInvalid,
			expectedCode: codes.Unauthenticated,
			expectedMsg:  "Token was invalid.",
		},
		{
			name:         "ErrUnauthorized",
			err:          ErrUnauthorized,
			expectedCode: codes.Unauthenticated,
			expectedMsg:  "Unauthorized.",
		},
		{
			name:         "ErrUserNotFound",
			err:          ErrUserNotFound,
			expectedCode: codes.NotFound,
			expectedMsg:  "用户不存在",
		},
		{
			name:         "ErrUserAlreadyExists",
			err:          ErrUserAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedMsg:  "User already exists.",
		},
		{
			name:         "ErrUserDisabled",
			err:          ErrUserDisabled,
			expectedCode: codes.PermissionDenied,
			expectedMsg:  "User is disabled.",
		},
		{
			name:         "InternalServerError",
			err:          InternalServerError,
			expectedCode: codes.Internal,
			expectedMsg:  "Internal server error.",
		},
		{
			name:         "custom error with message",
			err:          ErrUserNotFound.SetMessage("用户不存在"),
			expectedCode: codes.NotFound,
			expectedMsg:  "用户不存在",
		},
		{
			name:         "non-errorx error",
			err:          assert.AnError,
			expectedCode: codes.Internal,
			expectedMsg:  assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToGRPCError(tt.err)

			if tt.err == nil {
				assert.Nil(t, result)
				return
			}

			assert.NotNil(t, result)

			st, ok := status.FromError(result)
			assert.True(t, ok)

			assert.Equal(t, tt.expectedCode, st.Code())
			assert.Equal(t, tt.expectedMsg, st.Message())
		})
	}
}

func TestGetGRPCCode(t *testing.T) {
	tests := []struct {
		name         string
		errorxCode   int
		expectedCode codes.Code
	}{
		{"OK", 0, codes.OK},
		{"ErrBind", 400001, codes.InvalidArgument},
		{"ErrInvalidParameter", 400002, codes.InvalidArgument},
		{"ErrSignToken", 401001, codes.Unauthenticated},
		{"ErrTokenInvalid", 401002, codes.Unauthenticated},
		{"ErrUnauthorized", 401003, codes.Unauthenticated},
		{"ErrPasswordIncorrect", 401103, codes.Unauthenticated},
		{"ErrUserDisabled", 403104, codes.PermissionDenied},
		{"ErrResourceNotFound", 404001, codes.NotFound},
		{"ErrUserNotFound", 404102, codes.NotFound},
		{"ErrUserAlreadyExists", 409101, codes.AlreadyExists},
		{"InternalServerError", 500001, codes.Internal},
		{"Unknown 400", 400999, codes.InvalidArgument},
		{"Unknown 401", 401999, codes.Unauthenticated},
		{"Unknown 403", 403999, codes.PermissionDenied},
		{"Unknown 404", 404999, codes.NotFound},
		{"Unknown 409", 409999, codes.AlreadyExists},
		{"Unknown 500", 500999, codes.Internal},
		{"Unknown code", 999999, codes.Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getGRPCCode(tt.errorxCode)
			assert.Equal(t, tt.expectedCode, result)
		})
	}
}
