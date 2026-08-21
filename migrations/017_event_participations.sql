CREATE TABLE IF NOT EXISTS event_participations (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  event_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  name VARCHAR(64) NOT NULL,
  remark VARCHAR(255),
  status TINYINT DEFAULT 0 COMMENT '0:待审核 1:已通过 2:已签到 3:已取消',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_event_user (event_id, user_id),
  INDEX idx_event (event_id),
  INDEX idx_user (user_id),
  CONSTRAINT fk_event_participations_event FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
  CONSTRAINT fk_event_participations_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
