# 后端实际 JSON 响应格式（基于 CPN 代码）

根据 `CPN` 目录下的后端代码整理，供前端对接时参考。

---

## 1. 统一响应包装

所有接口均通过 `utils.SuccessResponse` / `utils.ErrorResponse` 包装：

```json
// 成功
{
  "code": 200,
  "data": { ... },
  "message": "success"
}

// 失败
{
  "code": 500,
  "data": null,
  "message": "错误信息"
}
```

**注意**：成功时 `code` 为 `200`（不是 0），失败时 `code` 为 `500`。前端判断应使用 `code === 200`。

---

## 2. 各接口实际返回结构

### 2.1 登录 `POST /login`

**成功 (200)**：
```json
{
  "code": 200,
  "data": {
    "id": "用户ID字符串",
    "username": "用户名",
    "roles": "user|admin",
    "source": "Competition_Platform"
  },
  "message": "success"
}
```

- 登录成功后后端会设置 Cookie：`LCP-Cookie`，前端需使用 `credentials: 'include'`。
- 后续请求会自动携带该 Cookie，无需手动传 token。

---

### 2.2 注册 `POST /register`

**成功 (200)**：
```json
{
  "code": 200,
  "data": "注册成功",
  "message": "success"
}
```

---

### 2.3 获取登录信息 `OPTIONS /getInfo`

**成功 (200)**：
```json
{
  "code": 200,
  "data": {
    "id": "用户ID",
    "username": "用户名",
    "roles": "user|admin",
    "source": "Competition_Platform"
  },
  "message": "success"
}
```

---

### 2.4 比赛列表 `GET /competition?pageNum=1&pageSize=10`

**成功 (200)**：
```json
{
  "code": 200,
  "data": {
    "total": 总数,
    "competitionList": [
      {
        "id": "赛题ID",
        "author": { "id", "username", "roles", "source" },
        "title": "比赛题目",
        "abstract": "比赛简介（文档链接）",
        "sort_order": "升序|降序",
        "time_limit": 20,
        "status": "准备中|进行中|已结束",
        "create_time": "2006-01-02 15:04:05",
        "latest_time": "2006-01-02 15:04:05"
      }
    ]
  },
  "message": "success"
}
```

**字段对应关系**（后端用下划线）：
| 前端建议字段 | 后端返回字段 |
|-------------|--------------|
| id | id |
| title | title |
| abstract | abstract |
| sortOrder | sort_order |
| timeLimit | time_limit |
| status | status |
| createTime | create_time |
| latestTime | latest_time |
| author | author |

---

### 2.5 单个比赛详情 `GET /competition/{competitionID}`

**成功 (200)**：同 `competitionList` 中单个元素结构，无列表包装。

---

### 2.6 提交记录列表 `GET /record?pageNum=1&pageSize=10`

- 管理员：返回所有记录
- 普通用户：仅返回自己的记录

**成功 (200)**：
```json
{
  "code": 200,
  "data": {
    "total": 总数,
    "recordList": [
      {
        "id": "记录ID",
        "user": { "id", "username", "roles", "source" },
        "competition": { 完整比赛对象 },
        "status": "上传完成|...",
        "run_time": -1,
        "score": -1,
        "errors": "",
        "create_time": "2006-01-02 15:04:05",
        "latest_time": "2006-01-02 15:04:05",
        "finish_time": "2006-01-02 15:04:05"
      }
    ]
  },
  "message": "success"
}
```

**字段对应关系**：
| 前端建议字段 | 后端返回字段 |
|-------------|--------------|
| id | id |
| userID | user.id |
| username | user.username |
| competitionID | competition.id |
| competitionTitle | competition.title |
| score | score |
| runTime | run_time |
| status | status |
| errors | errors |
| createTime | create_time |

---

### 2.7 比赛排名 `GET /competition/{competitionID}/record`

**成功 (200)**：
```json
{
  "code": 200,
  "data": {
    "total": 总数,
    "recordList": [
      {
        "id": "记录ID",
        "user": { ... },
        "competition": { ... },
        "status": "...",
        "run_time": 123,
        "score": 80,
        "errors": "",
        "create_time": "...",
        "latest_time": "...",
        "finish_time": "...",
        "rank": 1
      }
    ]
  },
  "message": "success"
}
```

注意：`recordList` 中每项多一个 `rank` 字段表示排名。

---

### 2.8 单个提交详情 `GET /competition/{competitionID}/record/{recordID}`

**成功 (200)**：同 `recordList` 中单个元素，无列表包装。

---

### 2.9 用户列表 `GET /user?pageNum=1&pageSize=10`（管理员）

**成功 (200)**：
```json
{
  "code": 200,
  "data": {
    "total": 总数,
    "userList": [
      {
        "id": "用户ID",
        "username": "用户名",
        "school": "学校",
        "email": "邮箱",
        "roles": "user|admin",
        "source": "Competition_Platform",
        "is_delete": false,
        "create_time": "2006-01-02 15:04:05",
        "latest_time": "2006-01-02 15:04:05"
      }
    ]
  },
  "message": "success"
}
```

---

### 2.10 当前用户信息 `GET /user/{userID}`

无论 path 传什么，都返回当前登录用户：

**成功 (200)**：
```json
{
  "code": 200,
  "data": {
    "id": "用户ID",
    "username": "用户名",
    "roles": "user|admin",
    "source": "Competition_Platform"
  },
  "message": "success"
}
```

---

## 3. 创建/更新接口的请求参数

### 3.1 创建比赛 `POST /competition`

**formData 必填**：
- `title`：比赛题目
- `abstract`：比赛简介（文档链接）
- `dockerImage`：赛题 docker 镜像名

**formData 可选**：
- `sortOrder`：`升序` | `降序`，默认 `降序`
- `timeLimit`：字符串，如 `"20"`，默认 `"20"`

### 3.2 更新比赛 `PUT /competition/{competitionID}`

所有字段可选：`title`, `abstract`, `sortOrder`, `timeLimit`, `dockerImage`, `status`（`准备中` | `进行中` | `已结束`）

### 3.3 提交代码 `POST /competition/{competitionID}/record`

**formData 必填**：
- `submission`：文件（必须是 zip）

---

## 4. 前端需要适配的改动汇总

1. **判断成功**：用 `code === 200`（不要用 `code === 0`）
2. **取数据**：`data` 即为业务数据，无需再解一层
3. **比赛列表**：用 `data.competitionList`，字段用 `abstract`、`sort_order`、`time_limit`、`create_time` 等（或前端统一转驼峰）
4. **记录列表**：用 `data.recordList`，字段用 `run_time`、`errors`、`create_time` 等
5. **用户列表**：用 `data.userList`
6. **认证方式**：Cookie `LCP-Cookie`，请求时 `credentials: 'include'`

---

## 5. 后端潜在问题（可选与后端沟通）

- `record/app.go` 中 `recordCreate` 部分错误返回直接用了 `context.JSON(500, "字符串")`，没有用 `utils.ErrorResponse`，前端拿到的可能是裸字符串，建议后端统一改成 `ErrorResponse`
- 用户 model 中 `roles` 在 bson 里是 `roles`，但数据库字段为 `roles`，与 swagger 的 `role` 单数略有差异，实际返回为 `roles`
