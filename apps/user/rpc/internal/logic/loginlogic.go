package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/clin211/miniblog-v3/apps/user/models"
	"github.com/clin211/miniblog-v3/apps/user/rpc/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/rpc/pb/rpc"
	"github.com/clin211/miniblog-v3/pkg/encrypt"
	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/clin211/miniblog-v3/pkg/token"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewLoginLogic 创建登录逻辑实例
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 用户登录
func (l *LoginLogic) Login(in *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	// 1. 参数验证
	if err := l.validateLoginRequest(in); err != nil {
		return nil, err
	}

	// 2. 检查账户锁定状态
	if l.isAccountLocked(in.Username) {
		return nil, errorx.ErrUnauthorized.SetMessage("账户已被锁定，请30分钟后重试")
	}

	// 3. 查询用户信息
	user, err := l.getUserByUsername(in.Username)
	if err != nil {
		l.recordFailedLogin(in.Username)
		return nil, errorx.ErrPasswordIncorrect.SetMessage("用户名或密码错误")
	}

	// 4. 验证密码
	if err := encrypt.Compare(user.Password, in.Password); err != nil {
		l.recordFailedLogin(in.Username)
		return nil, errorx.ErrPasswordIncorrect.SetMessage("用户名或密码错误")
	}

	// 5. 检查用户状态
	if user.Status != 1 {
		return nil, errorx.ErrUserDisabled.SetMessage("账户已被禁用")
	}

	// 6. 生成 JWT Token
	tokenStr, expireAt, err := token.Sign(user.UserId)
	if err != nil {
		logx.Errorf("生成Token失败: %v", err)
		return nil, errorx.ErrSignToken.SetMessage("生成Token失败")
	}

	// 7. 更新登录信息
	if err := l.updateLoginInfo(user, ""); err != nil {
		logx.Errorf("更新登录信息失败: %v", err)
		// 不返回错误，继续执行
	}

	// 8. 重置失败次数
	l.resetFailedLoginCount(in.Username)

	// 9. 缓存会话信息
	l.cacheSession(tokenStr, user)

	// 10. 记录登录日志
	logx.Infof("用户登录成功: user_id=%s, username=%s", user.UserId, user.Username)

	return &rpc.LoginResponse{
		Token:    tokenStr,
		ExpireAt: expireAt.Format(time.RFC3339),
	}, nil
}

// validateLoginRequest 验证登录请求参数
func (l *LoginLogic) validateLoginRequest(in *rpc.LoginRequest) error {
	// 验证用户名
	if in.Username == "" {
		return errorx.ErrInvalidParameter.SetMessage("用户名不能为空")
	}

	// 验证密码
	if in.Password == "" {
		return errorx.ErrInvalidParameter.SetMessage("密码不能为空")
	}

	return nil
}

// isAccountLocked 检查账户是否被锁定
func (l *LoginLogic) isAccountLocked(username string) bool {
	key := fmt.Sprintf("user:lock:%s", username)
	exists, err := l.svcCtx.Redis.Exists(key)
	if err != nil {
		l.Error("检查账户锁定状态失败", logx.Field("error", err))
		return false
	}
	return exists
}

// getUserByUsername 根据用户名/邮箱/手机号获取用户信息
func (l *LoginLogic) getUserByUsername(username string) (*models.Users, error) {
	// 先尝试通过用户名查找
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, username)
	if err == nil {
		return user, nil
	}

	// 再尝试通过邮箱查找
	user, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, username)
	if err == nil {
		return user, nil
	}

	// 最后尝试通过手机号查找
	user, err = l.svcCtx.UserModel.FindOneByPhone(l.ctx, username)
	if err == nil {
		return user, nil
	}

	// 如果都找不到，返回用户不存在错误
	if err == sqlx.ErrNotFound {
		return nil, errorx.ErrUserNotFound.SetMessage("用户不存在")
	}

	return nil, errorx.InternalServerError.SetMessage("查询用户信息失败")
}

// recordFailedLogin 记录失败登录
func (l *LoginLogic) recordFailedLogin(username string) error {
	key := fmt.Sprintf("user:failed:%s", username)
	count, err := l.svcCtx.Redis.Incr(key)
	if err != nil {
		l.Error("记录失败登录次数失败", logx.Field("error", err))
		return err
	}

	// 设置过期时间
	l.svcCtx.Redis.Expire(key, int(30*time.Minute.Seconds()))

	// 超过5次锁定账户
	if count >= 5 {
		lockKey := fmt.Sprintf("user:lock:%s", username)
		l.svcCtx.Redis.Set(lockKey, "locked")
		l.svcCtx.Redis.Expire(lockKey, int(30*time.Minute.Seconds()))
		l.Error("账户被锁定", logx.Field("username", username))
	}

	return nil
}

// updateLoginInfo 更新登录信息
func (l *LoginLogic) updateLoginInfo(user *models.Users, loginIP string) error {
	now := time.Now()
	user.LastLoginAt = sql.NullTime{Time: now, Valid: true}
	user.LastLoginIp = loginIP
	user.FailedLoginAttempts = 0

	return l.svcCtx.UserModel.Update(l.ctx, user)
}

// resetFailedLoginCount 重置失败登录次数
func (l *LoginLogic) resetFailedLoginCount(username string) {
	key := fmt.Sprintf("user:failed:%s", username)
	l.svcCtx.Redis.Del(key)
}

// cacheSession 缓存会话信息
func (l *LoginLogic) cacheSession(tokenStr string, user *models.Users) {
	key := fmt.Sprintf("user:session:%s", tokenStr)
	sessionData := map[string]string{
		"user_id":     user.UserId,
		"username":    user.Username,
		"login_time":  fmt.Sprintf("%d", time.Now().Unix()),
		"expire_time": fmt.Sprintf("%d", time.Now().Add(24*time.Hour).Unix()),
	}

	// 缓存会话信息
	l.svcCtx.Redis.Hmset(key, sessionData)
	l.svcCtx.Redis.Expire(key, int(24*time.Hour.Seconds()))
}
