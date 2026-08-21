CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT,
  action VARCHAR(64) NOT NULL,
  resource_type VARCHAR(64),
  resource_id BIGINT,
  detail JSON,
  ip VARCHAR(45),
  user_agent VARCHAR(255),
  result TINYINT COMMENT '1:成功 0:失败',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_user (user_id),
  INDEX idx_action (action),
  INDEX idx_created_at (created_at),
  INDEX idx_resource (resource_type, resource_id)
);
