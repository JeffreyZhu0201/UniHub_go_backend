# Drop database unihub;
# create database unihub;

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
('dept:list','List Departments', NOW(), NOW()),
('ding:create','Create Ding', NOW(), NOW()),
('leave:approve','Approval leave', NOW(), NOW()),
('class:join', 'Join Class', NOW(), NOW());

INSERT INTO role_permissions (role_id, permission_id) VALUES
((SELECT id FROM roles WHERE `key` = 'super_admin'), (SELECT id FROM permissions WHERE code = 'dept:create')),
((SELECT id FROM roles WHERE `key` = 'super_admin'), (SELECT id FROM permissions WHERE code = 'class:create')),
((SELECT id FROM roles WHERE `key` = 'admin'), (SELECT id FROM permissions WHERE code = 'dept:create')),
((SELECT id FROM roles WHERE `key` = 'admin'), (SELECT id FROM permissions WHERE code = 'class:create')),

((SELECT id FROM roles WHERE `key` = 'counselor'), (SELECT id FROM permissions WHERE code = 'dept:create')),
((SELECT id FROM roles WHERE `key` = 'counselor'), (SELECT id FROM permissions WHERE code = 'dept:list')),
((SELECT id FROM roles WHERE `key` = 'counselor'), (SELECT id FROM permissions WHERE code = 'class:create')),
((SELECT id FROM roles WHERE `key` = 'counselor'), (SELECT id FROM permissions WHERE code = 'leave:approve')),
((SELECT id FROM roles WHERE `key` = 'counselor'), (SELECT id FROM permissions WHERE code = 'ding:create')),

((SELECT id FROM roles WHERE `key` = 'teacher'), (SELECT id FROM permissions WHERE code = 'class:create')),
((SELECT id FROM roles WHERE `key` = 'teacher'), (SELECT id FROM permissions WHERE code = 'ding:create')),
((SELECT id FROM roles WHERE `key` = 'student'), (SELECT id FROM permissions WHERE code = 'class:join')),
((SELECT id FROM roles WHERE `key` = 'student'), (SELECT id FROM permissions WHERE code = 'dept:join'));