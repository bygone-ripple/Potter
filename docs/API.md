API 文档

## 基础信息

- **Base URL**: `http://localhost:8080/api`
- **认证方式**: Cookie (`user-session`)
- **Content-Type**: `application/json`

---

## 1. 认证模块 (Auth)

### 1.1 用户登录

**接口**: `POST /auth/login`

**描述**: 用户登录，成功后会设置 session cookie

**请求示例**:
```json
{
  "name": "girlsbandcry",
  "password": "ninaninanina"
}
```

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "girlsbandcry",
    "avatar": "https://localhost:8080/uploads/avatar_1145141919.png"
  }
}
```

---

### 1.2 用户登出

**接口**: `POST /auth/logout`

**描述**: 用户登出，清除 session

**请求**: 无 body

**响应示例** (200 OK):
```json
{
  "success": true
}
```

---

## 2. 用户模块 (User)

### 2.1 用户注册

**接口**: `POST /users/`

**描述**: 注册新用户。用户昵称不允许重复

**请求示例**:
```json
{
  "name": "bocchibocchi",
  "avatar": "https://localhost:8080/uploads/avatar_1145141919.png",
  "password": "kitakita"
}
```

**字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 用户 ID |
| name | string | 是 | 用户昵称 |
| password | string | 是 | 用户密码 |
| avatar | string | 否 | 用户头像url |

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 2,
    "name": "bocchibocchi",
    "avatar": "https://localhost:8080/uploads/avatar_1145141919.png"
  }
}
```

---

### 2.2 获取个人信息（登录状态）

**接口**: `GET /users/me`

**权限**: 需要登录（role >= 1）

**描述**: 返回当前用户的个人信息，若未登录则返回错误

**请求**: 无 body

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 7,
    "name": "cyrene",
    "avatar":"https://localhost:8080/uploads/avatar_1145141919.png"
  }
}
```

**错误响应**
- 未登录:
```json
{
  "success": false,
  "message": "鉴权错误: 您未登录\n",
  "code": 6
}
```

---

### 2.3 更新个人信息

**接口**: `PUT /users/me`

**权限**: 需要登录 (role >= 1)

**描述**: 更新当前用户的个人信息

**请求示例**:
```json
{
  "id": 1,
  "name": "gbc",
  "avatar": "https://localhost:8080/uploads/avatar_1145141919.png",
  "password": "buchikome"
}
```

**字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 用户 ID |
| name | string | 否 | 用户昵称，不填则不更新 |
| avatar | string | 否 | 用户头像url，不填则不更新 |
| password | string | 否 | 用户密码，不填则不更新 |

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "gbc",
    "avatar": "https://localhost:8080/uploads/avatar_1145141919.png"
  }
}
```

---

### 2.4 获取我发布的锅单

**接口**: `GET /users/me/posted-tasks`

**权限**: 需要登录 (role >= 1)

**描述**: 获取当前用户发布的所有锅单列表，锅单只会返回有值的字段

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | integer | 是 | 页码，从 1 开始 |
| limit | integer | 是 | 每页数量 |

**请求示例**:
```
GET /users/me/posted-tasks/?page=1&limit=2
```

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": [
      {
          "id": 13,
          "name": "拯救世界",
          "depart": "tech",
          "description": "技术宅拯救世界",
          "ddl": "2025-12-15T18:30:00+08:00",
          "level": 3,
          "status": 1,
          "uris": [
              "https://localhost:8080/uploads/task_1145141919.png"
          ],
          "posterID": 5
      }
  ]
}
```

---

### 2.5 获取我接取的锅单

**接口**: `GET /users/me/assigned-tasks`

**权限**: 需要登录 (role >= 1)

**描述**: 获取当前用户接取的所有锅单列表，不会返回锅单完整信息

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | integer | 是 | 页码，从 1 开始 |
| limit | integer | 是 | 每页数量 |

**请求示例**:
```
GET /users/me/assigned-tasks/?page=1&limit=2
```

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": [
    {
      "id": 5,
      "name": "修复首页加载缓慢问题",
      "depart": "tech",
      "description": "首页加载时间超过3秒，需要优化数据库查询和前端资源加载",
      "ddl": "2025-11-16T02:00:00+08:00",
      "level": 4,
      "status": 2,
      "uris": [
          "/static/screenshot_1730246400.png",
          "/static/performance_report_1730246401.pdf"
      ],
      "posterID": 1,
      "assigneeID": 1
    }
  ]
}
```

---

## 3. 锅单模块 (Task)

### 3.1 创建锅单

**接口**: `POST /tasks/`

**权限**: 需要登录 (role >= 1)

**描述**: 创建新的锅单，会将 Status 默认设置为 1。锅单名称不允许重复

**请求示例**:
```json
{
  "name": "锅单示例",
  "depart": "tech",
  "description": "这是一个测试锅单，用于演示请求格式",
  "ddl": "2025-12-15T18:30:00+08:00",
  "level": 3,
  "uris": [
    "https://localhost:8080/static/spec_1730246400.pdf",
    "https://localhost:8080/static/design_1730246401.png"
  ],
}
```

**字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 锅单名称 |
| depart | string | 是 | 部门，可选值 `tech`/`video`/`art` |
| description | string | 是 | 锅单描述 |
| ddl | string | 是 | 截止时间，RFC3339 格式 |
| level | integer | 是 | 难度等级，1-5 |
| uris | array | 否 | 附件 URI 数组 |

**响应示例** (201 Created):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "锅单示例",
    "depart": "tech",
    "description": "这是一个测试锅单，用于演示请求格式",
    "ddl": "2025-12-15T18:30:00+08:00",
    "level": 3,
    "status": 1,
    "uris": [
      "https://localhost:8080/static/spec_1730246400.pdf",
      "https://localhost:8080/static/design_1730246401.png"
    ],
    "posterID": 1
  }
}
```

---

### 3.2 获取锅单列表

**接口**: `GET /tasks/`

**权限**: 需要登录 (role >= 1)

**描述**: 分页查询锅单列表，支持多条件筛选

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | integer | 是 | 页码，从 1 开始 |
| limit | integer | 是 | 每页数量 |
| name | string | 否 | 锅单名称模糊搜索 |
| depart | string | 否 | 部门筛选：`tech`/`video`/`art` |
| status | integer | 否 | 状态筛选：1-未接取, 2-进行中, 3-已完成, 4-已弃置 |
| level | integer | 否 | 优先级筛选：1-5 |

**请求示例**:
```
GET /tasks/?page=1&limit=2&depart=tech&status=1
```

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": {
    "total": 25,
    "list": [
      {
        "id": 1,
        "name": "修复首页加载缓慢问题",
        "depart": "tech",
        "description": "首页加载时间超过3秒，需要优化",
        "ddl": "2025-11-15T18:00:00+08:00",
        "level": 4,
        "status": 1,
        "uris": [
          "https://localhost:8080/static/spec_1730246400.pdf",
          "https://localhost:8080/static/design_1730246401.png"
        ],
        "posterID": 1
      },
      {
        "id": 2,
        "name": "实现用户权限系统",
        "depart": "tech",
        "description": "添加角色和权限管理功能",
        "ddl": "2025-12-01T12:00:00+08:00",
        "level": 5,
        "status": 1,
        "posterID": 2,
        "assigneeID": 1
      }
    ]
  }
}
```

---

### 3.3 获取锅单详情

**接口**: `GET /tasks/{taskID}`

**权限**: 需要登录 (role >= 1)

**描述**: 获取指定锅单的完整信息，包括发布者、接取者、评论等

**路径参数**:
- `taskID`: 锅单 ID

**请求示例**:
```
GET /tasks/1
```

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "修复首页加载缓慢问题",
    "depart": "tech",
    "description": "首页加载时间超过3秒，需要优化数据库查询和前端资源加载",
    "ddl": "2025-11-15T18:00:00+08:00",
    "level": 4,
    "status": 2,
    "uris": [
      "https://localhost:8080/static/screenshot_1730246400.png",
      "https://localhost:8080/static/performance_report_1730246401.pdf"
    ],
    "comments": [
      {
        "id": 1,
        "content": "111已经开始处理",
        "time": "2025-11-01T10:00:00+08:00",
        "posterID": 2,
        "poster": {
          "id": 2,
          "name": "bocchibocchi",
          "avatar": "https://localhost:8080/uploads/avatar_1145141919.png"
        }
      }
    ],
    "posterID": 1,
    "poster": {
      "id": 1,
      "name": "girlsbandcry",
      "avatar": "https://localhost:8080/uploads/avatar_1145141919.png"
    },
    "assigneeID": 2,
    "assignee": {
      "id": 2,
      "name": "bocchibocchi",
      "avatar": "https://localhost:8080/uploads/avatar_1145141919.png"
    }
  }
}
```

---

### 3.4 更新锅单信息

**接口**: `PUT /tasks/{taskID}`

**权限**: 需要登录 (role >= 1)，且必须是锅单发布者

**描述**: 更新锅单信息，只有发布者可以修改

**路径参数**:
- `taskID`: 锅单 ID

**请求示例**:
```json
{
  "name": "【紧急】修复首页加载缓慢问题",
  "depart": "tech",
  "description": "首页加载时间超过3秒，需要优化数据库查询和前端资源加载。问题比预期严重，提升优先级。",
  "ddl": "2025-11-15T18:00:00+08:00",
  "level": 5,
  "status": 2,
  "uris": [
    "https://localhost:8080/static/screenshot_1730246400.png",
    "https://localhost:8080/static/performance_report_1730246401.pdf"
  ]
}
```

**字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 锅单名称 |
| depart | string | 是 | 部门，可选值 `tech`/`video`/`art` |
| description | string | 是 | 锅单描述 |
| ddl | string | 是 | 截止时间，RFC3339 格式 |
| level | integer | 是 | 难度等级，1-5 |
| status | integer | 是 | 状态：1-未接取, 2-进行中, 3-已完成, 4-已弃置 |
| uris | array | 是 | 附件 URI 数组，可传空数组 |

**重要说明**:
- ⚠️ **必须传递表中所有字段**（包括不想修改的字段），后端会用传入的完整数据更新锅单
- `posterID` 和 `assigneeID` 字段无法通过此接口修改
- 若要修改接取者，请使用接取/退出锅单的专用接口
- 锅单状态的未接取和进行中，在接取/退出锅单接口会自动更新

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "【紧急】修复首页加载缓慢问题",
    "depart": "tech",
    "description": "首页加载时间超过3秒，需要优化数据库查询和前端资源加载。问题比预期严重，提升优先级。",
    "ddl": "2025-11-15T18:00:00+08:00",
    "level": 5,
    "status": 2,
    "uris": [
      "https://localhost:8080/static/screenshot_1730246400.png",
      "https://localhost:8080/static/performance_report_1730246401.pdf"
    ],
    "posterID": 1,
    "poster": {
      "id": 1,
      "name": "girlsbandcry",
      "avatar":"https://localhost:8080/uploads/avatar_1145141919.png"
    },
    "assigneeID": 2,
    "assignee": {
      "id": 2,
      "name": "bocchibocchi",
      "avatar": "https://localhost:8080/uploads/avatar_1145141919.png"
    }
  }
}
```

---

### 3.5 删除锅单

**接口**: `DELETE /tasks/{taskID}`

**权限**: 需要登录 (role >= 1)，且必须是锅单发布者

**描述**: 删除锅单，只有发布者可以删除。删除后锅单名为被格式化为"name__deleted_timestamp"以避免重复

**路径参数**:
- `taskID`: 锅单 ID

**请求示例**:
```
DELETE /tasks/1
```

**响应** (204 No Content):
无响应体

---

### 3.6 更新锅单接取者

**接口**: `PUT /tasks/{taskID}/assignees`

**权限**: 需要登录 (role >= 1)，且必须是锅单发布者

**描述**: 发布者指定锅单的接取者，用于强制分配锅单

**路径参数**:
- `taskID`: 锅单 ID

**请求示例**:
```json
{
  "assigneeID": 1
}
```

**字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| assigneeID | integer | 是 | 接取者的用户 ID |

**响应示例** (200 OK):
```json
{
    "success": true,
    "data": {
        "id": 13,
        "name": "拯救世界",
        "depart": "tech",
        "description": "技术宅拯救世界",
        "ddl": "2025-12-15T18:30:00+08:00",
        "level": 3,
        "status": 2,
        "uris": [
            "https://localhost:8080/uploads/task_1145141919.png"
        ],
        "posterID": 5,
        "poster": {
            "id": 5,
            "name": "cyrene",
            "avatar":"https://localhost:8080/uploads/avatar_1145141919.png"
        },
        "assigneeID": 1,
        "assignee": {
            "id": 1,
            "name": "girlsbandcry",
            "avatar":"https://localhost:8080/uploads/avatar_1145141919.png"
        }
    }
}
```

**错误响应**:
- 不是发布者: 
```json
{
  "success": false,
  "message": "操作错误: 该锅单不存在或发布者不是您\n",
  "code": 5
}
```

---

### 3.7 接取锅单

**接口**: `POST /tasks/{taskID}/assignees/me`

**权限**: 需要登录 (role >= 1)

**描述**: 当前用户接取指定锅单，锅单状态会自动更新为"进行中"

**路径参数**:
- `taskID`: 锅单 ID

**请求**: 无 body

**请求示例**:
```
POST /tasks/1/assignees/me
```

**响应示例** (201 Created):
```json
{
    "success": true,
    "data": {
        "id": 13,
        "name": "拯救世界",
        "depart": "tech",
        "description": "技术宅拯救世界",
        "ddl": "2025-12-15T18:30:00+08:00",
        "level": 3,
        "status": 2,
        "uris": [
            "https://localhost:8080/uploads/task_1145141919.png"
        ],
        "posterID": 5,
        "poster": {
            "id": 5,
            "name": "cyrene",
            "avatar":"https://localhost:8080/uploads/avatar_1145141919.png"
        },
        "assigneeID": 5,
        "assignee": {
            "id": 5,
            "name": "cyrene",
            "avatar":"https://localhost:8080/uploads/avatar_1145141919.png"
        }
    }
}
```

**错误响应**:
- 锅单已有接取者: 
```json
{
    "success": false,
    "message": "操作错误: 该锅单已有接锅人\n",
    "code": 5
}
```

---

### 3.8 取消接锅

**接口**: `DELETE /tasks/{taskID}/assignees/me`

**权限**: 需要登录 (role >= 1)，且必须是当前锅单的接取者

**描述**: 当前用户退出已接取的锅单，锅单状态会恢复为"未接取"

**路径参数**:
- `taskID`: 锅单 ID

**请求示例**:
```
DELETE /tasks/1/assignees/me
```

**响应** (204 No Content):
无响应体

**错误响应**:
- 不是接取者: 
```json
{
    "success": false,
    "message": "操作错误: 该锅单不存在或接锅人不是您\n",
    "code": 5
}
```

---

### 3.9 发布评论

**接口**: `POST /tasks/{taskID}/comments`

**权限**: 需要登录 (role >= 1)

**描述**: 在指定锅单下发布评论

**路径参数**:
- `taskID`: 锅单 ID

**请求示例**:
```json
{
  "content": "bocchichansaikou"
}
```

**字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| content | string | 是 | 评论内容，允许重复 |

**响应示例** (201 Created):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "content": "bocchichansaikou",
    "time": "2025-11-01T10:30:00+08:00",
    "taskID": 1,
    "posterID": 2,
    "poster": {
      "id": 2,
      "name": "bocchibocchi",
      "avatar":"https://localhost:8080/uploads/avatar_1145141919.png"
    }
  }
}
```

---

## 4. 评论模块 (Comment)

### 4.1 删除评论

**接口**: `DELETE /comments/{commentID}`

**权限**: 需要登录 (role >= 1)

**描述**: 删除评论，允许评论发布者或对应锅单的发布者删除

**路径参数**:
- `commentID`: 评论 ID

**请求示例**:
```
DELETE /comments/1
```

**响应** (204 No Content):
无响应体

**权限说明**:
- 评论发布者可以删除自己的评论
- 锅单发布者可以删除该锅单下的任何评论

---

## 5. 文件上传

### 5.1 上传文件

**接口**: `POST /uploads`

**权限**: 需要登录 (role >= 1)

**描述**: 上传图片、视频等文件到服务器

**Content-Type**: `multipart/form-data`

**表单字段**:
- `file`: 文件（二进制）

**请求示例** (cURL):
```bash
curl -X POST http://localhost:8080/api/uploads \
  -H "Cookie: user-session=your-session-token" \
  -F "file=@/path/to/Furina.png"
```

**响应示例** (200 OK):
```json
{
  "success": true,
  "data": {
    "uri": "/static/Furina_1761825610.png"
  }
}
```

**说明**:
- 文件会保存在服务器的 `uploads` 目录
- 文件名格式: `原文件名_时间戳.扩展名`
- 返回的 URI 可用于锅单的 `uris` 字段
- 可通过 `http://localhost:8080/static/文件名` 访问上传的文件

---

## 6. 错误响应格式

所有错误响应统一格式：

```json
{
  "success": false,
  "message": "错误描述信息",
  "code": 5
}
```

**常见错误码**:
- 1: "内部错误"
- 2: "公开错误"
- 3: "参数错误"
- 4: "系统错误"
- 5: "操作错误"
- 6: "鉴权错误"
- 7: "权限错误"

---

## 7. 数据字典

### 7.1 部门 (depart)

| 值 | 说明 |
|----|------|
| tech | 技术部 |
| video | 视频部 |
| art | 美术部 |

### 7.2 锅单状态 (status)

| 值 | 说明 |
|----|------|
| 1 | 未接取 |
| 2 | 进行中 |
| 3 | 已完成 |
| 4 | 已弃置 |

### 7.3 难度 (level)

| 值 | 说明 |
|----|------|
| 1 | 轻松 |
| 2 | 较易 |
| 3 | 中等 |
| 4 | 较难 |
| 5 | 困难 |

---

## 8. 注意事项

1. **时间格式**: 所有时间字段使用 RFC3339 格式 (ISO 8601)，尽量传入东八区时间，如 `2025-11-15T18:00:00+08:00`
2. **认证**: 除了注册和登录接口，其他接口都需要携带 session cookie
3. **权限**: 大部分接口要求用户已登录
4. **分页**: 列表接口支持分页，`page` 从 1 开始
5. **文件访问**: 上传的文件通过 `/static/文件名` 访问，无需认证
6. **软删除**: 删除操作使用软删除
