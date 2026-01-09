# Postman API Test Collection

This document contains a list of API endpoints for the UniHub system, along with `curl` commands to test them.

**Base URL**: `http://localhost:8080`

## 1. Authentication

### 1.1 Register (Student)
Registers a new student account.

**Endpoint**: `POST /api/v1/auth/register`

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "Student One",
    "email": "student1@test.com",
    "password": "password123",
    "role_key": "student",
    "student_no": "2023001"
  }'
```

### 1.2 Register (Counselor)
Registers a new counselor account.

**Endpoint**: `POST /api/v1/auth/register`

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "Counselor One",
    "email": "counselor1@test.com",
    "password": "password123",
    "role_key": "counselor",
    "staff_no": "T2023001"
  }'
```

### 1.3 Register (Teacher)
Registers a new teacher account.

**Endpoint**: `POST /api/v1/auth/register`

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "Teacher One",
    "email": "teacher1@test.com",
    "password": "password123",
    "role_key": "teacher",
    "staff_no": "T2023002"
  }'
```

### 1.4 Login
Logs in a user and returns a JWT token.

**Endpoint**: `POST /api/v1/auth/login`

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student1@test.com",
    "password": "password123"
  }'
```

> **Note**: Save the `token` from the response for subsequent requests. Replace `YOUR_TOKEN` in the commands below.

---

## 2. Organization Management (Counselor & Teacher)

### 2.1 Create Department (Counselor)
Creates a new department.

**Endpoint**: `POST /api/v1/departments`

```bash
curl -X POST http://localhost:8080/api/v1/departments \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Computer Science Dept"
  }'
```

> **Response**: Note the `invite_code`.

### 2.2 List My Departments (Counselor)
Lists departments managed by the logged-in counselor.

**Endpoint**: `GET /api/v1/departments/mine`

```bash
curl -X GET http://localhost:8080/api/v1/departments/mine \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 2.3 Create Class (Teacher)
Creates a new class.

**Endpoint**: `POST /api/v1/classes`

```bash
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Math Class 101"
  }'
```

> **Response**: Note the `invite_code`.

### 2.4 List My Classes (Teacher)
Lists classes managed by the logged-in teacher.

**Endpoint**: `GET /api/v1/classes/mine`

```bash
curl -X GET http://localhost:8080/api/v1/classes/mine \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 3. Student Actions

### 3.1 Join Department
Student joins a department using an invite code.

**Endpoint**: `POST /api/v1/departments/join`

```bash
curl -X POST http://localhost:8080/api/v1/departments/join \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "invite_code": "DEPT_INVITE_CODE"
  }'
```

### 3.2 Join Class
Student joins a class using an invite code.

**Endpoint**: `POST /api/v1/classes/join`

```bash
curl -X POST http://localhost:8080/api/v1/classes/join \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "invite_code": "CLASS_INVITE_CODE"
  }'
```

---

## 4. Leave Management

### 4.1 Apply for Leave (Student)
Submits a leave request.

**Endpoint**: `POST /api/v1/leaves`

```bash
curl -X POST http://localhost:8080/api/v1/leaves \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "Sick Leave",
    "start_time": "2023-11-01T08:00:00Z",
    "end_time": "2023-11-03T18:00:00Z",
    "reason": "Flu"
  }'
```

### 4.2 List My Leaves (Student)
Lists student's own leave requests.

**Endpoint**: `GET /api/v1/leaves/mine`

```bash
curl -X GET http://localhost:8080/api/v1/leaves/mine \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4.3 List Pending Leaves (Counselor)
Lists pending leave requests from students in managed departments.

**Endpoint**: `GET /api/v1/leaves/pending`

```bash
curl -X GET http://localhost:8080/api/v1/leaves/pending \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 4.4 Audit Leave (Counselor)
Approves or rejects a leave request.

**Endpoint**: `POST /api/v1/leaves/:uuid/audit` 
*(Note: standard path convention uses parameter, but request body contains `leave_id`)*

```bash
curl -X POST http://localhost:8080/api/v1/leaves/audit/audit \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "leave_id": 1,
    "status": "approved"
  }'
```

---

## 5. Notifications

### 5.1 Create Notification (Counselor/Teacher)
Sends a notification to a department or class.

**Endpoint**: `POST /api/v1/notifications`

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Exam Schedule",
    "content": "Midterm exams start next week.",
    "target_type": "dept",
    "target_id": 1
  }'
```

### 5.2 List My Notifications (Student)
Lists notifications for the student's department and classes.

**Endpoint**: `GET /api/v1/notifications/mine`

```bash
curl -X GET http://localhost:8080/api/v1/notifications/mine \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 6. Ding Tasks (Sign-In/Dorm Check)

### 6.1 Create Ding Task (Counselor/Teacher)
Creates a location-based check-in task.

**Endpoint**: `POST /api/v1/createdings`

```bash
curl -X POST http://localhost:8080/api/v1/createdings \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Morning Check-in",
    "start_time": "2023-11-01T08:00:00Z",
    "end_time": "2023-11-01T09:00:00Z",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "radius": 500,
    "dept_id": 1
  }'
```

### 6.2 List My Dings (Student)
Lists pending and completed ding tasks.

**Endpoint**: `GET /api/v1/mydings`

```bash
curl -X GET http://localhost:8080/api/v1/mydings \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 6.3 List Created Dings (Counselor/Teacher)
Lists tasks created by the current user.

**Endpoint**: `GET /api/v1/mycreateddings`

```bash
curl -X GET http://localhost:8080/api/v1/mycreateddings \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 7. Open Platform

### 7.1 Register Developer
Registers a new developer account.

**Endpoint**: `POST /api/v1/open/register`

```bash
curl -X POST http://localhost:8080/api/v1/open/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dev One",
    "email": "dev1@test.com"
  }'
```

> **Response**: Note `dev_secret`.

### 7.2 Create App
Creates a new app for the developer.

**Endpoint**: `POST /api/v1/open/apps`

```bash
curl -X POST http://localhost:8080/api/v1/open/apps \
  -H "Content-Type: application/json" \
  -H "X-Dev-Secret: YOUR_DEV_SECRET" \
  -d '{
    "name": "My Cool App"
  }'
```

> **Response**: Note `app_id` and `app_secret`.

### 7.3 Access Public Profile (Third-party App)
Simulates a third-party app accessing public data.

**Endpoint**: `GET /api/v1/start/v1/user/:id/public_profile`

```bash
curl -X GET http://localhost:8080/api/v1/start/v1/user/1/public_profile \
  -H "X-App-ID: YOUR_APP_ID" \
  -H "X-App-Secret: YOUR_APP_SECRET"
```

---

## 8. User Info

### 8.1 Get Profile
Gets current user's profile.

**Endpoint**: `GET /api/v1/user/profile`

```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 8.2 List Students (Counselor/Teacher)
Lists students in their managed departments or classes.

**Endpoint**: `GET /api/v1/students`

```bash
curl -X GET http://localhost:8080/api/v1/students \
  -H "Authorization: Bearer YOUR_TOKEN"
```

