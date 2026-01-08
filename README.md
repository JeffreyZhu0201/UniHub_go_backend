# UniHub - 校园用户架构管理系统

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Gin Framework](https://img.shields.io/badge/Gin-v1.10-00ADD8?style=flat&logo=go)
![GORM](https://img.shields.io/badge/GORM-v1.25-00ADD8?style=flat&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=flat&logo=mysql)
![Redis](https://img.shields.io/badge/Redis-7.0-DC382D?style=flat&logo=redis)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green.svg)

UniHub 是一个基于 **Golang** (Gin + GORM) 开发的现代化校园后端管理系统。它旨在解决复杂的校园组织架构管理问题，支持多级用户角色（辅导员、教师、学生），并集成了通知、请假审批、签到查寝以及开放平台等扩展服务。

系统采用标准的 Clean Architecture 分层架构，代码结构清晰，易于维护和扩展。

## ✨ 核心功能

### 1. 完善的用户与权限体系
- **多角色支持**：系统内置超级管理员、学校管理员、辅导员、教师、学生五种角色。
- **RBAC 权限控制**：基于角色的权限访问控制，无需笨重的 Casbin，轻量级实现接口鉴权。
- **数据权限范围 (Data Scope)**：
  - 辅导员只能查看自己创建的部门及部门下的学生。
  - 教师只能查看自己创建的班级及班级下的学生。
  - 学生只能访问自己的数据。
- **安全认证**：基于 JWT (JSON Web Token) 的身份验证，密码采用 Bcrypt 加密存储。

### 2. 灵活的组织架构管理
- **部门管理**：辅导员可创建部门，生成**8位随机邀请码**。
- **班级管理**：教师可创建班级，生成**8位随机邀请码**。
- **便捷加入**：学生通过邀请码一键加入部门或班级，系统自动建立关联。

### 3.微服务化扩展功能
- **通知中心**：辅导员/教师可向指定部门或班级发送通知，学生实时查看。
- **请假审批流**：学生发起请假，辅导员在线审批（通过/驳回）。审批通过后系统自动生成“销假签到”任务。
- **任务系统**：
  - **签到/查寝**：支持发布位置签到或拍照查寝任务。
  - **截止时间控制**：自动判断任务是否逾期。
  - **防重复提交**：学生每次任务仅能提交一次。

### 4. 开放平台 (Open Platform)
- **开发者注册**：外部开发者可注册并获取 `DevSecret`。
- **应用管理**：开发者可创建应用，获取 `AppID` 和 `AppSecret`。
- **API网关**：内置 API 速率限制 (Rate Limiting)，防止接口滥用。

## 🛠 技术栈

- **语言**: Golang 1.25+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.0
- **配置管理**: Viper
- **容器化**: Docker & Docker Compose

## 📂 目录结构

```
.
├── cmd/
│   └── server/          # 程序入口
├── configs/             # 配置文件 (yaml)
├── docs/                # 文档 (API, SQL等)
├── internal/
│   ├── config/          # 配置加载逻辑
│   ├── db/              # 数据库连接初始化
│   ├── handler/         # HTTP 处理器 (Controller)
│   ├── model/           # GORM 数据库模型
│   ├── router/          # 路由定义
│   └── service/         # 业务逻辑 (权限检查等)
├── pkg/
│   ├── jwtutil/         # JWT 工具包
│   └── middleware/      # 中间件 (Auth, RateLimit)
├── scripts/             # 初始化脚本 (SQL)
├── tests/               # 单元测试
├── Dockerfile           # Docker 构建文件
└── docker-compose.yml   # 容器编排配置
```

## 🚀 快速开始

### 前置要求
- Go 1.24+
- MySQL 8.0+
- Redis

### 本地运行

1. **克隆项目**
   ```bash
   git clone https://github.com/your/unihub.git
   cd unihub
   ```

2. **配置环境**
   复制并修改配置文件（可选，默认使用 `configs/config.yaml`）：
   ```bash
   # 确保 config.yaml 中的 MySQL 和 Redis 地址指向你的本地服务
   ```

3. **初始化数据库**
   在 MySQL 中执行 `scripts/init.sql` 以创建基础角色数据：
   ```bash
   mysql -u root -p unihub < scripts/init.sql
   ```
   *注意：系统启动时会自动迁移表结构 (Auto Migrate)。*

4. **启动服务**
   ```bash
   go run cmd/server/main.go
   ```

### 🐳 Docker 部署 (推荐)

一键启动所有服务 (App + MySQL + Redis)：

```bash
docker-compose up --build -d
```

- **API 服务**: `http://localhost:8080`
- **MySQL**: 端口 3306
- **Redis**: 端口 6379

> 首次启动时，MySQL 容器会自动执行 `scripts/init.sql` 初始化数据。

## 📖 API 文档

详细的 API 接口说明及 Curl 示例请查阅：
👉 [UniHub API Documentation](docs/api.md)

主要模块包括：
- `/api/v1/auth/*`: 认证
- `/api/v1/user/*`: 用户
- `/api/v1/departments/*`: 部门管理
- `/api/v1/classes/*`: 班级管理
- `/api/v1/notifications/*`: 通知
- `/api/v1/leaves/*`: 请假
- `/api/v1/tasks/*`: 任务
- `/api/v1/open/*`: 开放平台

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

