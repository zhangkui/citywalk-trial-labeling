CREATE TABLE IF NOT EXISTS waypoints (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  route_id BIGINT NOT NULL,
  name VARCHAR(200),
  lat DECIMAL(10,7) NOT NULL,
  lng DECIMAL(10,7) NOT NULL,
  stay_duration INT DEFAULT 0 COMMENT '停留时长(分钟)',
  `order` INT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_route (route_id),
  CONSTRAINT fk_waypoints_route FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE
);
