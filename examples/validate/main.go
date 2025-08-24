// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package main

import (
	"fmt"
	"time"

	"github.com/clin211/miniblog-v3/pkg/validate"
)

// User 用户结构体示例
type User struct {
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

func main() {
	fmt.Println("=== 自定义参数校验系统演示 ===\n")

	// 示例1: 验证有效的用户数据
	fmt.Println("1. 验证有效的用户数据:")
	validUser := &User{
		ID:       "550e8400-e29b-41d4-a716-446655440000",
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "password123",
		Age:      25,
		Phone:    "13800138000",
		Gender:   "male",
		Status:   "active",
	}

	if err := validate.ValidateStructWithCustomRules(validUser); err != nil {
		fmt.Printf("❌ 验证失败: %v\n", err)
	} else {
		fmt.Println("✅ 验证通过")
	}

	fmt.Println()

	// 示例2: 验证无效的用户数据
	fmt.Println("2. 验证无效的用户数据:")
	invalidUser := &User{
		ID:       "invalid-uuid",
		Username: "jo", // 太短
		Email:    "invalid-email",
		Password: "123",         // 太短
		Age:      150,           // 超出范围
		Phone:    "12345678901", // 格式错误
		Gender:   "unknown",     // 不在枚举中
		Status:   "invalid",     // 不在枚举中
	}

	if err := validate.ValidateStructWithCustomRules(invalidUser); err != nil {
		fmt.Printf("❌ 验证失败:\n")
		if validationErrors, ok := err.(validate.ValidationErrors); ok {
			for i, validationError := range validationErrors {
				fmt.Printf("   %d. 字段: %s, 错误: %s, 值: %v\n",
					i+1, validationError.Field, validationError.Message, validationError.Value)
			}
		}
	} else {
		fmt.Println("✅ 验证通过")
	}

	fmt.Println()

	// 示例3: 展示验证规则
	fmt.Println("3. 支持的验证规则:")
	rules := validate.GetCommonRules()
	for name, rule := range rules {
		fmt.Printf("   - %s: %s\n", name, rule.Description)
	}

	fmt.Println()

	// 示例4: 性能测试
	fmt.Println("4. 性能测试:")
	start := time.Now()
	for i := 0; i < 1000; i++ {
		validate.ValidateStructWithCustomRules(validUser)
	}
	duration := time.Since(start)
	fmt.Printf("   1000次验证耗时: %v\n", duration)
	fmt.Printf("   平均每次验证耗时: %v\n", duration/1000)

	fmt.Println("\n=== 演示完成 ===")
}
