# CityWalk API

## 统一格式

```json
{ "code": 0, "message": "success", "data": {} }
```

## 认证

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`
- `GET /api/v1/auth/badges`
- `PUT /api/v1/auth/profile`
- `PUT /api/v1/auth/password`

## 路线

- `GET /api/v1/routes`
- `POST /api/v1/routes`
- `GET /api/v1/routes/{id}`
- `PUT /api/v1/routes/{id}`
- `DELETE /api/v1/routes/{id}`
- `PUT /api/v1/routes/{id}/status`
- `POST /api/v1/routes/{id}/favorite`
- `DELETE /api/v1/routes/{id}/favorite`
- `POST /api/v1/routes/{id}/rate`

## 故事 / 地标 / 活动 / 推荐

- `GET /api/v1/stories`
- `GET /api/v1/landmarks`
- `GET /api/v1/events`
- `GET /api/v1/recommend`
- `GET /api/health`
