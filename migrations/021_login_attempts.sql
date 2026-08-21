CREATE TABLE IF NOT EXISTS login_attempts (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64),
  ip VARCHAR(45),
  success TINYINT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_username (username),
  INDEX idx_ip (ip),
  INDEX idx_created_at (created_at)
);
