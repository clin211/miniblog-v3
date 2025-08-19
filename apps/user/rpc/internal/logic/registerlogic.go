package logic

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/clin211/miniblog-v3/apps/user/models"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"
	"github.com/clin211/miniblog-v3/pkg/encrypt"
	"github.com/clin211/miniblog-v3/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
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
		return nil, err
	}

	// 2. 检查用户唯一性
	if err := l.checkUserUniqueness(in); err != nil {
		return nil, err
	}

	// 3. 密码加密
	hashedPassword, err := encrypt.Encrypt(in.Password)
	if err != nil {
		logx.Errorf("密码加密失败: %v", err)
		return nil, errorx.InternalServerError.SetMessage("密码加密失败")
	}

	// 4. 生成用户ID
	userId := l.generateUserId()

	// 5. 创建用户记录
	user := &models.Users{
		UserId:         userId,
		Username:       in.Username,
		Password:       hashedPassword,
		Email:          in.Email,
		Phone:          in.Phone,
		Age:            int64(in.Age),
		Gender:         int64(in.Gender),
		Avatar:         in.Avatar,
		RegisterSource: int64(in.RegisterSource),
		WechatOpenid:   in.WechatOpenid,
		Status:         1, // 1-正常状态
		EmailVerified:  0, // 0-未验证
		PhoneVerified:  0, // 0-未验证
		IsRisk:         0, // 0-非风险用户
		FailedLoginAttempts: 0, // 0-失败登录次数
	}

	// 6. 保存到数据库
	if _, err := l.svcCtx.UserModel.Insert(l.ctx, user); err != nil {
		logx.Errorf("用户创建失败: %v", err)
		return nil, errorx.InternalServerError.SetMessage("用户创建失败")
	}

	// 7. 记录注册日志
	logx.Infof("用户注册成功: user_id=%s, username=%s, email=%s", userId, in.Username, in.Email)

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
	// 检查用户名是否已存在
	_, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err == nil {
		return errorx.ErrUserAlreadyExists.SetMessage("用户名已存在")
	}
	if err != errorx.ErrResourceNotFound {
		logx.Errorf("检查用户名唯一性失败: %v", err)
		return errorx.InternalServerError.SetMessage("检查用户名唯一性失败")
	}

	// 检查邮箱是否已存在
	_, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err == nil {
		return errorx.ErrUserAlreadyExists.SetMessage("邮箱已存在")
	}
	if err != errorx.ErrResourceNotFound {
		logx.Errorf("检查邮箱唯一性失败: %v", err)
		return errorx.InternalServerError.SetMessage("检查邮箱唯一性失败")
	}

	// 检查手机号是否已存在
	_, err = l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
	if err == nil {
		return errorx.ErrUserAlreadyExists.SetMessage("手机号已存在")
	}
	if err != errorx.ErrResourceNotFound {
		logx.Errorf("检查手机号唯一性失败: %v", err)
		return errorx.InternalServerError.SetMessage("检查手机号唯一性失败")
	}

	// 如果有微信OpenID，检查是否已存在
	if in.WechatOpenid != "" {
		_, err = l.svcCtx.UserModel.FindOneByWechatOpenid(l.ctx, in.WechatOpenid)
		if err == nil {
			return errorx.ErrUserAlreadyExists.SetMessage("微信账号已存在")
		}
		if err != errorx.ErrResourceNotFound {
			logx.Errorf("检查微信OpenID唯一性失败: %v", err)
			return errorx.InternalServerError.SetMessage("检查微信OpenID唯一性失败")
		}
	}

	return nil
}

// generateUserId 生成用户ID
func (l *RegisterLogic) generateUserId() string {
	// 生成时间戳 + 随机数的用户ID
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("user_%d_%04d", timestamp, randomNum)
}
