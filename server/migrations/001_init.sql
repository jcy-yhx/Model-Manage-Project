-- ============================================================
-- 企业大模型公共服务平台 (AI Gateway Console)
-- 数据库初始化脚本
-- 版本: V1.2
-- ============================================================

CREATE DATABASE IF NOT EXISTS ai_gateway
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE ai_gateway;

-- ============================================================
-- 1. users 用户表
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id           BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username     VARCHAR(64)  NOT NULL,
    display_name VARCHAR(64)  DEFAULT '',
    password_hash VARCHAR(256) NOT NULL,
    role         VARCHAR(32)  NOT NULL DEFAULT 'student' COMMENT 'admin / teacher / student',
    department   VARCHAR(128) DEFAULT '',
    status       TINYINT      NOT NULL DEFAULT 1 COMMENT '1=正常, 0=禁用',
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='用户表';


-- ============================================================
-- 2. projects 项目表
-- ============================================================
CREATE TABLE IF NOT EXISTS projects (
    id          BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    owner_id    BIGINT       NOT NULL,
    description VARCHAR(512) DEFAULT '',
    status      TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用, 0=停用',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_owner_id (owner_id),
    CONSTRAINT fk_projects_owner FOREIGN KEY (owner_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='项目表';


-- ============================================================
-- 3. api_keys API 密钥表
-- ============================================================
CREATE TABLE IF NOT EXISTS api_keys (
    id           BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    project_id   BIGINT       NOT NULL,
    name         VARCHAR(64)  DEFAULT '' COMMENT '密钥备注名',
    key_hash     VARCHAR(256) NOT NULL COMMENT 'bcrypt 密钥哈希',
    key_prefix   VARCHAR(16)  DEFAULT '' COMMENT '前缀展示（如 sk-xxxx）',
    status       TINYINT      NOT NULL DEFAULT 1 COMMENT '1=有效, 0=禁用',
    last_used_at DATETIME     DEFAULT NULL,
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_project_id (project_id),
    INDEX idx_key_prefix (key_prefix),
    INDEX idx_status (status),
    CONSTRAINT fk_api_keys_project FOREIGN KEY (project_id) REFERENCES projects(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='API 密钥表';


-- ============================================================
-- 4. models 模型表
-- ============================================================
CREATE TABLE IF NOT EXISTS models (
    id           BIGINT        NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(64)   NOT NULL COMMENT '模型标识名（如 qwen2.5-72b）',
    display_name VARCHAR(128)  DEFAULT '' COMMENT '展示名称',
    provider     VARCHAR(32)   DEFAULT '' COMMENT '提供商',
    status       VARCHAR(16)   NOT NULL DEFAULT 'online' COMMENT 'online / offline / degraded',
    input_price  DECIMAL(10,6) NOT NULL DEFAULT 0 COMMENT '输入价格（元/1K tokens）',
    output_price DECIMAL(10,6) NOT NULL DEFAULT 0 COMMENT '输出价格（元/1K tokens）',
    description  VARCHAR(256)  DEFAULT '',
    created_at   DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='模型表';


-- ============================================================
-- 5. api_usage_logs 调用日志表
-- ============================================================
CREATE TABLE IF NOT EXISTS api_usage_logs (
    id                BIGINT        NOT NULL AUTO_INCREMENT PRIMARY KEY,
    project_id        BIGINT        NOT NULL,
    api_key_id        BIGINT        NOT NULL,
    model_id          BIGINT        NOT NULL,
    prompt_tokens     INT           NOT NULL DEFAULT 0 COMMENT '输入 Token 数',
    completion_tokens INT           NOT NULL DEFAULT 0 COMMENT '输出 Token 数',
    total_tokens      INT           NOT NULL DEFAULT 0 COMMENT '总 Token 数',
    cost              DECIMAL(10,6) NOT NULL DEFAULT 0 COMMENT '本次调用费用（元）',
    latency_ms        INT           NOT NULL DEFAULT 0 COMMENT '延迟（毫秒）',
    status            VARCHAR(16)   NOT NULL DEFAULT 'success' COMMENT 'success / fail / timeout',
    error_message     VARCHAR(256)  DEFAULT '' COMMENT '错误信息（失败时）',
    created_at        DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_created_at (created_at),
    INDEX idx_model_id_created_at (model_id, created_at),
    INDEX idx_project_id_created_at (project_id, created_at),
    INDEX idx_api_key_id_created_at (api_key_id, created_at),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='API 调用日志表';


-- ============================================================
-- 6. daily_usage_stats 日汇总表
-- ============================================================
CREATE TABLE IF NOT EXISTS daily_usage_stats (
    id              BIGINT         NOT NULL AUTO_INCREMENT PRIMARY KEY,
    stat_date       DATE           NOT NULL COMMENT '统计日期',
    model_id        BIGINT         DEFAULT NULL COMMENT '关联模型（NULL=全局汇总）',
    total_tokens    BIGINT         NOT NULL DEFAULT 0,
    total_requests  INT            NOT NULL DEFAULT 0,
    success_count   INT            NOT NULL DEFAULT 0,
    fail_count      INT            NOT NULL DEFAULT 0,
    total_cost      DECIMAL(12,6)  NOT NULL DEFAULT 0,
    avg_latency_ms  INT            NOT NULL DEFAULT 0,

    UNIQUE INDEX idx_date_model (stat_date, model_id),
    INDEX idx_stat_date (stat_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
  COMMENT='日汇总统计表';
