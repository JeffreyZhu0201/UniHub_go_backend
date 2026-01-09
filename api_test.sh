#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

# Check dependencies
if ! command -v jq &> /dev/null; then
    echo "Error: jq is not installed. Please install jq to run this script."
    exit 1
fi

echo "Cleaning up previous data (if possible)..."
RAND=$((RANDOM))
COUNSELOR_USER="counselor_$RAND"
TEACHER_USER="teacher_$RAND"
STUDENT_USER="student_$RAND"
PASSWORD="password123"

echo "---------------------------------------------------"
echo "1. Register Counselor ($COUNSELOR_USER)"
resp=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$COUNSELOR_USER\",
    \"password\": \"$PASSWORD\",
    \"role_key\": \"counselor\",
    \"staff_no\": \"C$RAND\"
  }")
echo "Response: $resp"
COUNSELOR_TOKEN=$(echo $resp | jq -r .token)

if [ "$COUNSELOR_TOKEN" == "null" ] || [ -z "$COUNSELOR_TOKEN" ]; then
    echo "Failed to register counselor"
    exit 1
fi
echo "Counselor Token acquired."

echo "---------------------------------------------------"
echo "2. Create Department"
resp=$(curl -s -X POST "$BASE_URL/departments" \
  -H "Authorization: Bearer $COUNSELOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Computer Science Dept"
  }')
echo "Response: $resp"
DEPT_INVITE_CODE=$(echo $resp | jq -r .invite_code)
echo "Department Invite Code: $DEPT_INVITE_CODE"

echo "---------------------------------------------------"
echo "3. Register Teacher ($TEACHER_USER)"
resp=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEACHER_USER\",
    \"password\": \"$PASSWORD\",
    \"role_key\": \"teacher\",
    \"staff_no\": \"T$RAND\"
  }")
echo "Response: $resp"
TEACHER_TOKEN=$(echo $resp | jq -r .token)
echo "Teacher Token acquired."

echo "---------------------------------------------------"
echo "4. Create Class"
resp=$(curl -s -X POST "$BASE_URL/classes" \
  -H "Authorization: Bearer $TEACHER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Software Engineering 101"
  }')
echo "Response: $resp"
CLASS_INVITE_CODE=$(echo $resp | jq -r .invite_code)
echo "Class Invite Code: $CLASS_INVITE_CODE"

echo "---------------------------------------------------"
echo "5. Register Student ($STUDENT_USER)"
resp=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$STUDENT_USER\",
    \"password\": \"$PASSWORD\",
    \"role_key\": \"student\",
    \"student_no\": \"S$RAND\"
  }")
echo "Response: $resp"
STUDENT_TOKEN=$(echo $resp | jq -r .token)
echo "Student Token acquired."

echo "---------------------------------------------------"
echo "6. Join Department"
resp=$(curl -s -X POST "$BASE_URL/departments/join" \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"invite_code\": \"$DEPT_INVITE_CODE\"
  }")
echo "Response: $resp"

echo "---------------------------------------------------"
echo "7. Join Class"
resp=$(curl -s -X POST "$BASE_URL/classes/join" \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"invite_code\": \"$CLASS_INVITE_CODE\"
  }")
echo "Response: $resp"

echo "---------------------------------------------------"
echo "8. Create Notification (Counselor)"
echo "Fetching Counselor Departments to get numeric ID..."
resp=$(curl -s -X GET "$BASE_URL/departments/mine" \
  -H "Authorization: Bearer $COUNSELOR_TOKEN")
# Pick the first one
REAL_DEPT_ID=$(echo $resp | jq -r '.[0].ID')
echo "Real Dept ID: $REAL_DEPT_ID"

resp=$(curl -s -X POST "$BASE_URL/notifications" \
  -H "Authorization: Bearer $COUNSELOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Welcome\",
    \"content\": \"Welcome to the department!\",
    \"target_type\": \"dept\",
    \"target_id\": $REAL_DEPT_ID
  }")
echo "Response: $resp"

echo "---------------------------------------------------"
echo "9. Check Notifications (Student)"
resp=$(curl -s -X GET "$BASE_URL/notifications/mine" \
  -H "Authorization: Bearer $STUDENT_TOKEN")
echo "Response: $resp"

echo "---------------------------------------------------"
echo "10. Apply Leave (Student)"
resp=$(curl -s -X POST "$BASE_URL/leaves" \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"type\": \"Sick\",
    \"start_time\": \"2023-11-20T09:00:00Z\",
    \"end_time\": \"2023-11-21T18:00:00Z\",
    \"reason\": \"Fever\"
  }")
echo "Response: $resp"
LEAVE_UUID=$(echo $resp | jq -r .id)
echo "Leave UUID: $LEAVE_UUID"

echo "---------------------------------------------------"
echo "11. Audit Leave (Counselor)"
resp=$(curl -s -X POST "$BASE_URL/leaves/$LEAVE_UUID/audit" \
  -H "Authorization: Bearer $COUNSELOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"status\": \"approved\"
  }")
echo "Response: $resp"

echo "---------------------------------------------------"
echo "12. Create Task (Teacher)"
echo "Fetching Teacher Classes..."
resp=$(curl -s -X GET "$BASE_URL/classes/mine" \
  -H "Authorization: Bearer $TEACHER_TOKEN")
REAL_CLASS_ID=$(echo $resp | jq -r '.[0].ID')
echo "Real Class ID: $REAL_CLASS_ID"

resp=$(curl -s -X POST "$BASE_URL/tasks" \
  -H "Authorization: Bearer $TEACHER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Class Sign In\",
    \"type\": \"sign_in\",
    \"target_type\": \"class\",
    \"target_id\": $REAL_CLASS_ID,
    \"deadline\": \"2026-12-31T23:59:59Z\",
    \"config\": {}
  }")
echo "Response: $resp"
TASK_UUID=$(echo $resp | jq -r .id)
echo "Task UUID: $TASK_UUID"

echo "---------------------------------------------------"
echo "13. Submit Task (Student)"
resp=$(curl -s -X POST "$BASE_URL/tasks/$TASK_UUID/submit" \
  -H "Authorization: Bearer $STUDENT_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"data\": {\"location\": \"Classroom 101\"}
  }")
echo "Response: $resp"

echo "---------------------------------------------------"
echo "Test Sequence Complete."

