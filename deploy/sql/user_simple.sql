-- 用户表
CREATE TABLE `users` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '自增 ID',
    `user_id` VARCHAR(32) NOT NULL DEFAULT '' COMMENT '用户ID',
    `age` INT DEFAULT 0 COMMENT '年龄',
    `avatar` VARCHAR(255) DEFAULT '' COMMENT '头像URL',
    `username` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '密码',
    `password_updated_at` TIMESTAMP NULL COMMENT '密码更新时间',
    `email` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '邮箱',
    `email_verified` TINYINT DEFAULT 0 COMMENT '邮箱是否已验证；1-已验证,0-未验证',
    `phone` VARCHAR(20) NOT NULL DEFAULT '' COMMENT '手机号',
    `phone_verified` TINYINT DEFAULT 0 COMMENT '手机号是否已验证；1-已验证,0-未验证',
    `gender` TINYINT DEFAULT 0 COMMENT '性别：0-未设置，1-男，2-女，3-其他',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    `failed_login_attempts` INT DEFAULT 0 COMMENT '失败登录次数，超过5次则锁定账户，登录成功后重置',
    `last_login_at` TIMESTAMP NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) DEFAULT '' COMMENT '最后登录IP',
    `is_risk` TINYINT DEFAULT 0 COMMENT '是否为风险用户；1-是,0-否',
    `register_source` TINYINT DEFAULT 1 COMMENT '注册来源：1-web，2-app，3-wechat，4-qq，5-github，6-google',
    `register_ip` VARCHAR(45) DEFAULT '' COMMENT '注册IP',
    `wechat_openid` VARCHAR(100) DEFAULT '' COMMENT '微信OpenID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL COMMENT '删除时间',

    -- 唯一索引（保证数据唯一性，必须保留）
    UNIQUE KEY uk_user_id (`user_id`),
    UNIQUE KEY uk_username (`username`),
    UNIQUE KEY uk_email (`email`),
    UNIQUE KEY uk_phone (`phone`),
    UNIQUE KEY uk_wechat_openid (`wechat_openid`),

    -- 基础查询索引（最常用的）
    INDEX idx_status (`status`),
    INDEX idx_deleted_at (`deleted_at`)
) COMMENT='用户表' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
