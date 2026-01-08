# UniHub API 文档

## 基本信息
- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: 所有受保护接口需要在 Header 中携带 Token
  - `Authorization: Bearer <your_token>`

---

## 1. 认证模块 (Auth)

### 用户注册
**POST** `/auth/register`

支持学生、辅导员、教师注册。注册后自动返回 Token。

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |
| role_key | string | 是 | 角色标识: `student`, `counselor`, `teacher` |
| org_id | uint | 否 | 组织ID (可选) |
| staff_no | string | 否 | 工号 (教师/辅导员必填) |
| student_no | string | 否 | 学号 (学生必填) |

**Curl示例:**
```bash
# 学生注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "student001",
    "password": "password123",
    "role_key": "student",
    "student_no": "2024001"
  }'

# 辅导员注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "counselor001",
    "password": "password123",
    "role_key": "counselor",
    "staff_no": "T001"
  }'
```

### 用户登录
**POST** `/auth/login`

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "student001",
    "password": "password123"
  }'
```

---

## 2. 用户模块 (User)

### 获取个人资料
**GET** `/user/profile`

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <TOKEN>"
```

### 获取学生列表
**GET** `/students`

- **辅导员**: 返回其管理的所有部门下的学生。
- **教师**: 返回其管理的所有班级下的学生。

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/students \
  -H "Authorization: Bearer <TOKEN>"
```

---

## 3. 组织架构模块 (Org)

### 创建部门 (辅导员)
**POST** `/departments`

系统自动生成8位邀请码。

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| name | string | 是 | 部门名称 |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Authorization: Bearer <COUNSELOR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "计算机2024级1系"}'
```

### 查看我的部门 (辅导员)
**GET** `/departments/mine`

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/departments/mine \
  -H "Authorization: Bearer <COUNSELOR_TOKEN>"
```

### 创建班级 (教师)
**POST** `/classes`

系统自动生成8位邀请码。

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| name | string | 是 | 班级名称 |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Authorization: Bearer <TEACHER_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "高等数学A班"}'
```

### 查看我的班级 (教师)
**GET** `/classes/mine`

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/classes/mine \
  -H "Authorization: Bearer <TEACHER_TOKEN>"
```

### 加入部门 (学生)
**POST** `/departments/join`

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| invite_code | string | 是 | 8位邀请码 |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/departments/join \
  -H "Authorization: Bearer <STUDENT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"invite_code": "ABCDEFGH"}'
```

### 加入班级 (学生)
**POST** `/classes/join`

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| invite_code | string | 是 | 8位邀请码 |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/classes/join \
  -H "Authorization: Bearer <STUDENT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"invite_code": "XYZ12345"}'
```

---

## 4. 通知服务 (Notification)

### 发布通知 (辅导员/教师)
**POST** `/notifications`

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| title | string | 是 | 标题 |
| content | string | 是 | 内容 |
| target_type | string | 是 | 目标类型: `dept` (部门) 或 `class` (班级) |
| target_id | uint | 是 | 部门ID 或 班级ID |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "紧急会议",
    "content": "请所有同学下午2点到会议室集合",
    "target_type": "dept",
    "target_id": 1
  }'
```

### 查看我的通知 (学生)
**GET** `/notifications/mine`

返回学生所在部门和班级的所有通知。

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/notifications/mine \
  -H "Authorization: Bearer <STUDENT_TOKEN>"
```

---

## 5. 请假服务 (Leave)

### 提交请假申请 (学生)
**POST** `/leaves`

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| type | string | 是 | 请假类型 (如: 病假, 事假) |
| start_time | string | 是 | 开始时间 (RFC3339) |
| end_time | string | 是 | 结束时间 (RFC3339) |
| reason | string | 是 | 请假理由 |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/leaves \
  -H "Authorization: Bearer <STUDENT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "病假",
    "start_time": "2024-01-20T08:00:00Z",
    "end_time": "2024-01-22T18:00:00Z",
    "reason": "发烧"
  }'
```

### 查看我的请假记录 (学生)
**GET** `/leaves/mine`

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/leaves/mine \
  -H "Authorization: Bearer <STUDENT_TOKEN>"
```

### 查看待审批请假 (辅导员)
**GET** `/leaves/pending`

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/leaves/pending \
  -H "Authorization: Bearer <COUNSELOR_TOKEN>"
```

### 审批请假 (辅导员)
**POST** `/leaves/:uuid/audit`

审批通过后，系统会自动为学生生成一个“销假签到”任务。

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| status | string | 是 | 审批状态: `approved` (通过) 或 `rejected` (驳回) |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/leaves/some-uuid-123/audit \
  -H "Authorization: Bearer <COUNSELOR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"status": "approved"}'
```

---

## 6. 任务服务 (Task)

### 发布任务 (辅导员/教师)
**POST** `/tasks`

支持发布签到、查寝等任务。

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| title | string | 是 | 任务标题 |
| type | string | 是 | 任务类型: `sign_in` (签到) 或 `dorm_check` (查寝) |
| target_type | string | 是 | 目标类型: `dept` 或 `class` |
| target_id | uint | 是 | 目标ID |
| deadline | string | 是 | 截止时间 (RFC3339) |
| config | object | 否 | 任务配置 (如位置坐标等) |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "晚点名",
    "type": "dorm_check",
    "target_type": "dept",
    "target_id": 1,
    "deadline": "2024-01-20T23:00:00Z",
    "config": {"lat": 30.123, "lng": 120.456, "radius": 100}
  }'
```

### 查看我的任务 (学生)
**GET** `/tasks/mine`

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/tasks/mine \
  -H "Authorization: Bearer <STUDENT_TOKEN>"
```

### 提交任务 (学生)
**POST** `/tasks/:uuid/submit`

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| data | object | 是 | 提交内容 (如 {"location": "...", "photo": "..."}) |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks/some-uuid-123/submit \
  -H "Authorization: Bearer <STUDENT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "data": {"photo_url": "http://oss.example.com/photo.jpg", "location": "宿舍楼A栋"}
  }'
```

---

## 7. 开放平台 (Open Platform)

### 注册开发者
**POST** `/api/v1/open/register`

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| name | string | 是 | 开发者名称/公司名 |
| email | string | 是 | 联系邮箱 |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/open/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "第三方科技公司",
    "email": "dev@example.com"
  }'
```

### 创建应用
**POST** `/api/v1/open/apps`

需要在 Header 中携带注册时返回的 Developer Secret。

**Header:**
- `X-Dev-Secret`: `<YOUR_DEV_SECRET>`

**请求参数:**
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| name | string | 是 | 应用名称 |

**Curl示例:**
```bash
curl -X POST http://localhost:8080/api/v1/open/apps \
  -H "X-Dev-Secret: <YOUR_DEV_SECRET>" \
  -H "Content-Type: application/json" \
  -d '{"name": "智慧校园考勤助手"}'
```

### 调用开放接口 (带限流)
**GET** `/api/v1/start/v1/user/:id/public_profile`

需要在 Header 中携带 App ID 和 App Secret。

**Header:**
- `X-App-ID`: `<APP_ID>`
- `X-App-Secret`: `<APP_SECRET>`

**Curl示例:**
```bash
curl http://localhost:8080/api/v1/start/v1/user/1/public_profile \
  -H "X-App-ID: <APP_ID>" \
  -H "X-App-Secret: <APP_SECRET>"
```

