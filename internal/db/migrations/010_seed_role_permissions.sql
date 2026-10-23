-- Seed default role_permissions (SPEC.md §5 "Seed: Default role permissions")
INSERT INTO role_permissions (role, resource, action) VALUES
    ('super_admin', 'users', 'create'), ('super_admin', 'users', 'read'), ('super_admin', 'users', 'update'), ('super_admin', 'users', 'delete'),
    ('super_admin', 'objects', 'create'), ('super_admin', 'objects', 'read'), ('super_admin', 'objects', 'update'), ('super_admin', 'objects', 'delete'),
    ('super_admin', 'schemas', 'create'), ('super_admin', 'schemas', 'read'), ('super_admin', 'schemas', 'update'), ('super_admin', 'schemas', 'delete'),
    ('super_admin', 'horses', 'create'), ('super_admin', 'horses', 'read'), ('super_admin', 'horses', 'update'), ('super_admin', 'horses', 'delete'),
    ('super_admin', 'equipment', 'create'), ('super_admin', 'equipment', 'read'), ('super_admin', 'equipment', 'update'), ('super_admin', 'equipment', 'delete'),
    ('super_admin', 'lesson_codes', 'create'), ('super_admin', 'lesson_codes', 'read'), ('super_admin', 'lesson_codes', 'update'), ('super_admin', 'lesson_codes', 'delete'),
    ('super_admin', 'dashboards', 'create'), ('super_admin', 'dashboards', 'read'), ('super_admin', 'dashboards', 'update'), ('super_admin', 'dashboards', 'delete'),
    ('super_admin', 'reports', 'read'), ('super_admin', 'reports', 'export'),

    ('admin', 'objects', 'create'), ('admin', 'objects', 'read'), ('admin', 'objects', 'update'), ('admin', 'objects', 'delete'),
    ('admin', 'horses', 'create'), ('admin', 'horses', 'read'), ('admin', 'horses', 'update'), ('admin', 'horses', 'delete'),
    ('admin', 'equipment', 'create'), ('admin', 'equipment', 'read'), ('admin', 'equipment', 'update'), ('admin', 'equipment', 'delete'),
    ('admin', 'lesson_codes', 'create'), ('admin', 'lesson_codes', 'read'), ('admin', 'lesson_codes', 'update'), ('admin', 'lesson_codes', 'delete'),
    ('admin', 'dashboards', 'read'),
    ('admin', 'reports', 'read'), ('admin', 'reports', 'export'),

    ('trainer', 'objects', 'create'), ('trainer', 'objects', 'read'), ('trainer', 'horses', 'read'), ('trainer', 'equipment', 'read'), ('trainer', 'lesson_codes', 'read'),
    ('groom', 'objects', 'create'), ('groom', 'objects', 'read'), ('groom', 'horses', 'read'), ('groom', 'equipment', 'read'), ('groom', 'lesson_codes', 'read'),
    ('arit', 'objects', 'create'), ('arit', 'objects', 'read'), ('arit', 'horses', 'read'),
    ('cleaner', 'objects', 'create'), ('cleaner', 'objects', 'read'), ('cleaner', 'horses', 'read');
