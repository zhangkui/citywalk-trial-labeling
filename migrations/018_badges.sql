CREATE TABLE IF NOT EXISTS badges (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) UNIQUE NOT NULL,
  description VARCHAR(255),
  icon VARCHAR(64),
  condition_type VARCHAR(32) NOT NULL COMMENT 'register/create_route/create_story...',
  condition_value INT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
