-- ============================================================
-- 企业大模型公共服务平台 (AI Gateway Console)
-- 种子数据脚本
-- 版本: V1.2
-- ============================================================

USE ai_gateway;

-- ============================================================
-- 用户数据（3 个）
-- 原始密码:
--   admin         → admin123
--   teacher_wang  → teacher123
--   student_zhang → student123
-- 哈希方式: bcrypt, cost=12
-- ============================================================
INSERT INTO users (id, username, display_name, password_hash, role, department, status) VALUES
(1, 'admin',          '系统管理员', '$2a$12$Gmcx2/VE4Ggh1Ohy2e0fAu/dOvuRlijAVRVS8r8ZyZUgWrzgiofS.', 'admin',   '信息中心',   1),
(2, 'teacher_wang',   '王老师',     '$2a$12$QS4GZz76rHEMTlF3kMqGTu7SKafa4ZXuOCrb1c3EHj8ebUV3clzJe', 'teacher', '计算机学院', 1),
(3, 'student_zhang',  '张同学',     '$2a$12$VX86XhrUE0dddqvj2kAbPurASEHRTCbjmnfGA5hy8N6XHMFZ3plKi', 'student', '计算机学院', 1)
ON DUPLICATE KEY UPDATE username = VALUES(username);


-- ============================================================
-- 项目数据（1 个）
-- ============================================================
INSERT INTO projects (id, name, owner_id, description, status) VALUES
(1, 'AI 智能问答系统', 1, '面向全校师生的智能问答服务平台', 1)
ON DUPLICATE KEY UPDATE name = VALUES(name);


-- ============================================================
-- API Key 数据（1 个）
-- 原始 Key: sk-example-key-for-demo-20260805
-- 哈希方式: bcrypt, cost=12
-- key_prefix: sk-example-k  (取原始 Key 前 12 字符)
-- ============================================================
INSERT INTO api_keys (id, project_id, name, key_hash, key_prefix, status) VALUES
(1, 1, '默认测试 Key',
 '$2a$12$BMvJ0c5SQyBk/Q9Zb02I0eo.nji1Sy2AeyZ5PYW29/zqfbgCvb9ia',
 'sk-example-k', 1)
ON DUPLICATE KEY UPDATE name = VALUES(name);


-- ============================================================
-- 模型数据（5 个）及定价
-- 价格单位: 元 / 1K tokens
-- ============================================================
INSERT INTO models (id, name, display_name, provider, status, input_price, output_price, description) VALUES
(1, 'qwen2.5-72b',      'Qwen2.5-72B',      'Qwen',    'online', 0.004, 0.012, '通义千问 2.5 72B 参数模型'),
(2, 'deepseek-v3',      'DeepSeek-V3',       'DeepSeek','online', 0.002, 0.008, 'DeepSeek V3 大模型'),
(3, 'minimax-text-01',  'MiniMax-Text-01',   'MiniMax', 'online', 0.015, 0.060, 'MiniMax 文本大模型'),
(4, 'glm-4',            'GLM-4',             'GLM',     'online', 0.050, 0.050, '智谱 GLM-4 大模型'),
(5, 'qwen2.5-coder',    'Qwen2.5-Coder',     'Qwen',    'online', 0.004, 0.012, '通义千问 2.5 编程模型')
ON DUPLICATE KEY UPDATE name = VALUES(name);
