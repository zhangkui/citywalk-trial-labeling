CREATE TABLE IF NOT EXISTS comments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  target_type VARCHAR(20) NOT NULL COMMENT 'route/story/landmark',
  target_id BIGINT NOT NULL,
  content TEXT NOT NULL,
  parent_id BIGINT DEFAULT 0,
  created_by BIGINT NOT NULL,
  like_count INT DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME,
  INDEX idx_target (target_type, target_id),
  INDEX idx_parent (parent_id),
  INDEX idx_created_by (created_by),
  INDEX idx_created_at (created_at),
  CONSTRAINT fk_comments_user FOREIGN KEY (created_by) REFERENCES users(id)
);
