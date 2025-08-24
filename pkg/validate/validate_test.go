// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package validate

import (
	"testing"
	"time"
)

// TestUser 测试用的用户结构体
type TestUser struct {
	ID        string    `json:"id" valid:"required,uuid"`
	Username  string    `json:"username" valid:"required,alphanum,length(3|20)"`
	Email     string    `json:"email" valid:"required,email"`
	Password  string    `json:"password" valid:"required,length(6|50)"`
	Age       int       `json:"age" valid:"range(1|120)"`
	Phone     string    `json:"phone" valid:"matches(^1[3-9]\\d{9}$)"`
	Gender    string    `json:"gender" valid:"in(male|female|other)"`
	Status    string    `json:"status" valid:"in(active|inactive|pending)"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}

// TestValidUser 测试有效的用户数据
func TestValidUser(t *testing.T) {
	user := &TestUser{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "password123",
		Age:      25,
		Phone:    "13800138000",
		Gender:   "male",
		Status:   "active",
	}

	err := ValidateStructWithCustomRules(user)
	if err != nil {
		t.Errorf("有效用户数据验证失败: %v", err)
	}
}

// TestInvalidUser 测试无效的用户数据
func TestInvalidUser(t *testing.T) {
	user := &TestUser{
		ID:       "invalid-uuid",
		Username: "jo", // 太短
		Email:    "invalid-email",
		Password: "123",         // 太短
		Age:      150,           // 超出范围
		Phone:    "12345678901", // 格式错误
		Gender:   "unknown",     // 不在枚举中
		Status:   "invalid",     // 不在枚举中
	}

	err := ValidateStructWithCustomRules(user)
	if err == nil {
		t.Error("无效用户数据应该验证失败")
		return
	}

	// 检查错误数量
	if validationErrors, ok := err.(ValidationErrors); ok {
		if len(validationErrors) < 7 {
			t.Errorf("期望至少7个验证错误，实际得到 %d 个", len(validationErrors))
		}
	}
}

// TestRequiredField 测试必填字段验证
func TestRequiredField(t *testing.T) {
	user := &TestUser{
		Username: "john_doe",
		Email:    "john@example.com",
		// 缺少必填字段 ID 和 Password
	}

	err := ValidateStructWithCustomRules(user)
	if err == nil {
		t.Error("缺少必填字段应该验证失败")
		return
	}

	if validationErrors, ok := err.(ValidationErrors); ok {
		foundID := false
		foundPassword := false

		for _, validationError := range validationErrors {
			if validationError.Field == "ID" {
				foundID = true
			}
			if validationError.Field == "Password" {
				foundPassword = true
			}
		}

		if !foundID {
			t.Error("应该报告ID字段的必填错误")
		}
		if !foundPassword {
			t.Error("应该报告Password字段的必填错误")
		}
	}
}

// TestEmailValidation 测试邮箱验证
func TestEmailValidation(t *testing.T) {
	testCases := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
		{"test@", false},
		{"@example.com", false},
		{"test.example.com", false},
		{"", false}, // 空值应该验证失败（因为有required标签）
	}

	for _, tc := range testCases {
		user := &TestUser{
			ID:       "550e8400-e29b-41d4-a716-446655440000",
			Username: "testuser",
			Email:    tc.email,
			Password: "password123",
		}

		err := ValidateStructWithCustomRules(user)
		if tc.valid && err != nil {
			t.Errorf("邮箱 %s 应该有效，但验证失败: %v", tc.email, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("邮箱 %s 应该无效，但验证通过", tc.email)
		}
	}
}

// TestLengthValidation 测试长度验证
func TestLengthValidation(t *testing.T) {
	testCases := []struct {
		username string
		valid    bool
	}{
		{"john", true},
		{"jo", false},                        // 太短
		{"verylongusername123456789", false}, // 太长
		{"john_doe", true},
		{"john123", true},
		{"", false}, // 空值（必填）
	}

	for _, tc := range testCases {
		user := &TestUser{
			ID:       "550e8400-e29b-41d4-a716-446655440000",
			Username: tc.username,
			Email:    "test@example.com",
			Password: "password123",
		}

		err := ValidateStructWithCustomRules(user)
		if tc.valid && err != nil {
			t.Errorf("用户名 %s 应该有效，但验证失败: %v", tc.username, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("用户名 %s 应该无效，但验证通过", tc.username)
		}
	}
}

// TestRangeValidation 测试范围验证
func TestRangeValidation(t *testing.T) {
	testCases := []struct {
		age   int
		valid bool
	}{
		{1, true},
		{25, true},
		{120, true},
		{121, false}, // 太大
	}

	for _, tc := range testCases {
		user := &TestUser{
			ID:       "550e8400-e29b-41d4-a716-446655440000",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Age:      tc.age,
		}

		err := ValidateStructWithCustomRules(user)
		if tc.valid && err != nil {
			t.Errorf("年龄 %d 应该有效，但验证失败: %v", tc.age, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("年龄 %d 应该无效，但验证通过", tc.age)
		}
	}
}

// TestEnumValidation 测试枚举值验证
func TestEnumValidation(t *testing.T) {
	testCases := []struct {
		gender string
		valid  bool
	}{
		{"male", true},
		{"female", true},
		{"other", true},
		{"unknown", false},
		{"", true}, // 空值跳过验证
	}

	for _, tc := range testCases {
		user := &TestUser{
			ID:       "550e8400-e29b-41d4-a716-446655440000",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Gender:   tc.gender,
		}

		err := ValidateStructWithCustomRules(user)
		if tc.valid && err != nil {
			t.Errorf("性别 %s 应该有效，但验证失败: %v", tc.gender, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("性别 %s 应该无效，但验证通过", tc.gender)
		}
	}
}

// TestRegexValidation 测试正则表达式验证
func TestRegexValidation(t *testing.T) {
	testCases := []struct {
		phone string
		valid bool
	}{
		{"13800138000", true},
		{"13900139000", true},
		{"12345678901", false},  // 格式错误
		{"1380013800", false},   // 太短
		{"138001380000", false}, // 太长
		{"", true},              // 空值跳过验证
	}

	for _, tc := range testCases {
		user := &TestUser{
			ID:       "550e8400-e29b-41d4-a716-446655440000",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Phone:    tc.phone,
		}

		err := ValidateStructWithCustomRules(user)
		if tc.valid && err != nil {
			t.Errorf("手机号 %s 应该有效，但验证失败: %v", tc.phone, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("手机号 %s 应该无效，但验证通过", tc.phone)
		}
	}
}

// TestUUIDValidation 测试UUID验证
func TestUUIDValidation(t *testing.T) {
	testCases := []struct {
		id    string
		valid bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550e8400-e29b-41d4-a716-446655440001", true},
		{"invalid-uuid", false},
		{"550e8400-e29b-41d4-a716", false}, // 不完整
		{"", false},                        // 空值（必填）
	}

	for _, tc := range testCases {
		user := &TestUser{
			ID:       tc.id,
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		err := ValidateStructWithCustomRules(user)
		if tc.valid && err != nil {
			t.Errorf("UUID %s 应该有效，但验证失败: %v", tc.id, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("UUID %s 应该无效，但验证通过", tc.id)
		}
	}
}

// BenchmarkValidation 性能测试
func BenchmarkValidation(b *testing.B) {
	user := &TestUser{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "password123",
		Age:      25,
		Phone:    "13800138000",
		Gender:   "male",
		Status:   "active",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateStructWithCustomRules(user)
	}
}
