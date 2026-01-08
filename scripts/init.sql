-- Initialize Roles
INSERT INTO roles (name, `key`, data_scope, created_at, updated_at) VALUES
('Super Admin', 'super_admin', 'all', NOW(), NOW()),
('School Admin', 'admin', 'dept_and_sub', NOW(), NOW()),
('Counselor', 'counselor', 'dept', NOW(), NOW()),
('Teacher', 'teacher', 'dept', NOW(), NOW()),
('Student', 'student', 'self', NOW(), NOW());

-- Initialize Permissions (Optional for now as per requirement 4, but good practice)
INSERT INTO permissions (code, name, created_at, updated_at) VALUES
('dept:create', 'Create Department', NOW(), NOW()),
('class:create', 'Create Class', NOW(), NOW()),
('dept:join', 'Join Department', NOW(), NOW()),
('class:join', 'Join Class', NOW(), NOW());

