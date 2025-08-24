// Copyright 2025 长林啊 &lt;767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v3.git.

package logic

import (
	"context"
	"database/sql"
	"regexp"

	"github.com/clin211/miniblog-v3/apps/user/models"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"
	"github.com/clin211/miniblog-v3/pkg/encrypt"
	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/clin211/miniblog-v3/pkg/rid"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Register 用户注册
func (l *RegisterLogic) Register(in *rpc.RegisterRequest) (*rpc.RegisterResponse, error) {
	// 1. 参数验证
	if err := l.validateRegisterRequest(in); err != nil {
		return nil, errorx.ToGRPCError(err)
	}

	// 2. 检查用户唯一性
	if err := l.checkUserUniqueness(in); err != nil {
		return nil, errorx.ToGRPCError(err)
	}

	// 3. 密码加密
	hashedPassword, err := encrypt.Encrypt(in.Password)
	if err != nil {
		logx.Errorf("密码加密失败: %v", err)
		return nil, errorx.ToGRPCError(errorx.InternalServerError.SetMessage("密码加密失败"))
	}

	// 4. 生成用户ID
	userId := rid.UserID.New()

	// 5. 使用模型构造用户实体并插入
	user := &models.Users{
		UserId:              userId,
		Username:            in.Username,
		Password:            hashedPassword,
		Email:               in.Email,
		Phone:               in.Phone,
		Age:                 int64(in.Age),
		Gender:              int64(in.Gender),
		Avatar:              in.Avatar,
		RegisterSource:      int64(in.RegisterSource),
		WechatOpenid:        l.getWechatOpenid(in.WechatOpenid),
		Status:              1, // 1-正常状态
		EmailVerified:       0, // 0-未验证
		PhoneVerified:       0, // 0-未验证
		IsRisk:              0, // 0-非风险用户
		FailedLoginAttempts: 0, // 0-失败登录次数
	}

	if _, err := l.svcCtx.UserModel.Insert(l.ctx, user); err != nil {
		return nil, errorx.ToGRPCError(errorx.InternalServerError.SetMessage("用户创建失败"))
	}

	// 6. 记录注册日志
	logx.Infof("用户注册成功: %s, %s, %s", userId, in.Username, in.Email)

	return &rpc.RegisterResponse{
		UserId: userId,
	}, nil
}

// validateRegisterRequest 验证注册请求参数
func (l *RegisterLogic) validateRegisterRequest(in *rpc.RegisterRequest) error {
	// 验证用户名
	if in.Username == "" {
		return errorx.ErrInvalidParameter.SetMessage("用户名不能为空")
	}
	if len(in.Username) < 3 || len(in.Username) > 20 {
		return errorx.ErrInvalidParameter.SetMessage("用户名长度必须在3-20个字符之间")
	}
	// 用户名只能包含字母、数字、下划线
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(in.Username) {
		return errorx.ErrInvalidParameter.SetMessage("用户名只能包含字母、数字、下划线")
	}

	// 验证密码
	if in.Password == "" {
		return errorx.ErrInvalidParameter.SetMessage("密码不能为空")
	}
	if len(in.Password) < 6 || len(in.Password) > 32 {
		return errorx.ErrInvalidParameter.SetMessage("密码长度必须在6-32个字符之间")
	}
	// 密码必须包含字母和数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(in.Password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(in.Password)
	if !hasLetter || !hasDigit {
		return errorx.ErrInvalidParameter.SetMessage("密码必须包含字母和数字")
	}

	// 验证邮箱
	if in.Email == "" {
		return errorx.ErrInvalidParameter.SetMessage("邮箱不能为空")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(in.Email) {
		return errorx.ErrInvalidParameter.SetMessage("邮箱格式不正确")
	}

	// 验证手机号
	if in.Phone == "" {
		return errorx.ErrInvalidParameter.SetMessage("手机号不能为空")
	}
	if len(in.Phone) != 11 {
		return errorx.ErrInvalidParameter.SetMessage("手机号必须是11位数字")
	}
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(in.Phone) {
		return errorx.ErrInvalidParameter.SetMessage("手机号格式不正确")
	}

	// 验证年龄
	if in.Age < 1 || in.Age > 120 {
		return errorx.ErrInvalidParameter.SetMessage("年龄必须在1-120岁之间")
	}

	// 验证性别
	if in.Gender < 0 || in.Gender > 3 {
		return errorx.ErrInvalidParameter.SetMessage("性别值必须在0-3之间")
	}

	// 验证注册来源
	if in.RegisterSource < 1 || in.RegisterSource > 6 {
		return errorx.ErrInvalidParameter.SetMessage("注册来源值必须在1-6之间")
	}

	return nil
}

// checkUserUniqueness 检查用户唯一性
func (l *RegisterLogic) checkUserUniqueness(in *rpc.RegisterRequest) error {
	type uniqueCheck struct {
		name    string
		query   func() error
		onExist error
	}

	// 统一处理查询结果：nil=存在、ErrNotFound=不存在、其他=查询异常
	existsOrErr := func(name string, query func() error, onExist error) error {
		if err := query(); err != nil {
			if err == sqlx.ErrNotFound {
				return nil
			}
			logx.Errorf("检查%s唯一性失败: %v", name, err)
			return errorx.InternalServerError.SetMessage("%s", "检查"+name+"唯一性失败")
		}
		return onExist
	}

	checks := []uniqueCheck{
		{
			name: "用户名",
			query: func() error {
				_, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
				return err
			},
			onExist: errorx.ErrUserAlreadyExists.SetMessage("用户名已存在"),
		},
		{
			name: "邮箱",
			query: func() error {
				_, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
				return err
			},
			onExist: errorx.ErrUserAlreadyExists.SetMessage("邮箱已存在"),
		},
		{
			name: "手机号",
			query: func() error {
				_, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
				return err
			},
			onExist: errorx.ErrUserAlreadyExists.SetMessage("手机号已存在"),
		},
	}

	// 可选检查：仅当 WechatOpenid 非空时加入
	if in.WechatOpenid != "" {
		checks = append(checks, uniqueCheck{
			name: "微信OpenID",
			query: func() error {
				wechatOpenid := sql.NullString{String: in.WechatOpenid, Valid: true}
				_, err := l.svcCtx.UserModel.FindOneByWechatOpenid(l.ctx, wechatOpenid)
				return err
			},
			onExist: errorx.ErrUserAlreadyExists.SetMessage("微信账号已存在"),
		})
	}

	for _, c := range checks {
		if err := existsOrErr(c.name, c.query, c.onExist); err != nil {
			return err
		}
	}

	return nil
}

// getWechatOpenid 处理微信OpenID字段，空字符串返回NULL
func (l *RegisterLogic) getWechatOpenid(wechatOpenid string) sql.NullString {
	if wechatOpenid == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: wechatOpenid, Valid: true}
}
