-- Seed object_definitions + field_definitions for the 7 Phase 1 activity types
-- (SPEC.md §3.1). Schema-driven: generic form builder renders these, no custom
-- component per activity type.
--
-- Deviations from the literal spec text (per project owner):
--   - fee_amount -> fee
--   - separate date + time fields merged into one `datetime` field

INSERT INTO object_definitions (name, label, plural_label, icon, description) VALUES
    ('riding', 'Riding Log', 'Riding Logs', 'horse', 'Trainer logs a riding session with a horse'),
    ('training', 'Horse Training Log', 'Horse Training Logs', 'dumbbell', 'Instructor logs a horse training session'),
    ('horse_activity', 'Horse Activity Log', 'Horse Activity Logs', 'activity', 'Groom logs daily horse care (feed, bathe, walk, groom, etc.)'),
    ('consumable_inventory', 'Consumable Inventory Log', 'Consumable Inventory Logs', 'package', 'Stock in/out for consumable items (serbuk, rumput, vital, etc.)'),
    ('equipment_log', 'Equipment Assignment Log', 'Equipment Assignment Logs', 'tool', 'Equipment assigned, returned, moved, or inspected'),
    ('cleaning', 'Facility Cleaning Log', 'Facility Cleaning Logs', 'spray-can', 'Cleaner logs facility cleaning tasks'),
    ('grass_delivery', 'Grass Delivery Log', 'Grass Delivery Logs', 'truck', 'Grass cutter (Arit) logs grass delivery');

-- riding (date only — no separate time field in spec, so nothing to merge here)
INSERT INTO field_definitions (object_id, name, label, type, required, position, metadata)
SELECT id, 'date', 'Date', 'date', true, 1, '{}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'rider', 'Rider', 'relation', true, 2, '{"relation_to_object":"users"}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'horse', 'Horse', 'relation', true, 3, '{"relation_to_object":"horses"}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'code_lesson', 'Lesson Code', 'relation', true, 4, '{"relation_to_object":"lesson_codes"}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'instructor', 'Instructor', 'relation', false, 5, '{"relation_to_object":"users"}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'fee', 'Fee', 'number', false, 6, '{}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'duration_hours', 'Duration (hours)', 'number', false, 7, '{}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'notes', 'Notes', 'text', false, 8, '{}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'equipment_used', 'Equipment Used', 'multiselect', false, 9, '{"options":["saddle","bridle","halter","blanket"]}'::jsonb FROM object_definitions WHERE name = 'riding'
UNION ALL SELECT id, 'photos', 'Photos', 'file', true, 10, '{"min":1,"max":5,"maxSizeMB":10}'::jsonb FROM object_definitions WHERE name = 'riding';

-- training (date only)
INSERT INTO field_definitions (object_id, name, label, type, required, position, metadata)
SELECT id, 'date', 'Date', 'date', true, 1, '{}'::jsonb FROM object_definitions WHERE name = 'training'
UNION ALL SELECT id, 'horse', 'Horse', 'relation', true, 2, '{"relation_to_object":"horses"}'::jsonb FROM object_definitions WHERE name = 'training'
UNION ALL SELECT id, 'code_lesson', 'Lesson Code', 'relation', true, 3, '{"relation_to_object":"lesson_codes"}'::jsonb FROM object_definitions WHERE name = 'training'
UNION ALL SELECT id, 'instructor', 'Instructor', 'relation', true, 4, '{"relation_to_object":"users"}'::jsonb FROM object_definitions WHERE name = 'training'
UNION ALL SELECT id, 'duration_hours', 'Duration (hours)', 'number', false, 5, '{}'::jsonb FROM object_definitions WHERE name = 'training'
UNION ALL SELECT id, 'fee', 'Fee', 'number', false, 6, '{}'::jsonb FROM object_definitions WHERE name = 'training'
UNION ALL SELECT id, 'notes', 'Notes', 'text', false, 7, '{}'::jsonb FROM object_definitions WHERE name = 'training'
UNION ALL SELECT id, 'equipment_used', 'Equipment Used', 'multiselect', false, 8, '{"options":["saddle","bridle","halter","blanket"]}'::jsonb FROM object_definitions WHERE name = 'training'
UNION ALL SELECT id, 'photos', 'Photos', 'file', true, 9, '{"min":1,"max":5,"maxSizeMB":10}'::jsonb FROM object_definitions WHERE name = 'training';

-- horse_activity (date+time merged into datetime)
INSERT INTO field_definitions (object_id, name, label, type, required, position, metadata)
SELECT id, 'datetime', 'Date & Time', 'datetime', true, 1, '{}'::jsonb FROM object_definitions WHERE name = 'horse_activity'
UNION ALL SELECT id, 'horse', 'Horse', 'relation', true, 2, '{"relation_to_object":"horses"}'::jsonb FROM object_definitions WHERE name = 'horse_activity'
UNION ALL SELECT id, 'activity_type', 'Activity Type', 'select', true, 3, '{"options":["feed","bathe","walk","groom","vet_check","trim","stretch"]}'::jsonb FROM object_definitions WHERE name = 'horse_activity'
UNION ALL SELECT id, 'duration_minutes', 'Duration (minutes)', 'number', false, 4, '{}'::jsonb FROM object_definitions WHERE name = 'horse_activity'
UNION ALL SELECT id, 'notes', 'Notes', 'text', false, 5, '{}'::jsonb FROM object_definitions WHERE name = 'horse_activity'
UNION ALL SELECT id, 'photos', 'Photos', 'file', true, 6, '{"min":1,"max":5,"maxSizeMB":10}'::jsonb FROM object_definitions WHERE name = 'horse_activity';

-- consumable_inventory (date+time merged into datetime)
INSERT INTO field_definitions (object_id, name, label, type, required, position, metadata)
SELECT id, 'datetime', 'Date & Time', 'datetime', true, 1, '{}'::jsonb FROM object_definitions WHERE name = 'consumable_inventory'
UNION ALL SELECT id, 'item_type', 'Item Type', 'select', true, 2, '{"options":["serbuk","rumput","vital","equijewel","hay","chaff"]}'::jsonb FROM object_definitions WHERE name = 'consumable_inventory'
UNION ALL SELECT id, 'quantity_in', 'Quantity In', 'number', false, 3, '{}'::jsonb FROM object_definitions WHERE name = 'consumable_inventory'
UNION ALL SELECT id, 'quantity_out', 'Quantity Out', 'number', false, 4, '{}'::jsonb FROM object_definitions WHERE name = 'consumable_inventory'
UNION ALL SELECT id, 'unit', 'Unit', 'select', true, 5, '{"options":["kg","liter"]}'::jsonb FROM object_definitions WHERE name = 'consumable_inventory'
UNION ALL SELECT id, 'current_stock', 'Current Stock', 'number', false, 6, '{"readonly":true,"calculated":true}'::jsonb FROM object_definitions WHERE name = 'consumable_inventory'
UNION ALL SELECT id, 'horse_allocated', 'Horse Allocated', 'relation', false, 7, '{"relation_to_object":"horses"}'::jsonb FROM object_definitions WHERE name = 'consumable_inventory'
UNION ALL SELECT id, 'notes', 'Notes', 'text', false, 8, '{}'::jsonb FROM object_definitions WHERE name = 'consumable_inventory';

-- equipment_log (date+time merged into datetime)
INSERT INTO field_definitions (object_id, name, label, type, required, position, metadata)
SELECT id, 'datetime', 'Date & Time', 'datetime', true, 1, '{}'::jsonb FROM object_definitions WHERE name = 'equipment_log'
UNION ALL SELECT id, 'equipment', 'Equipment', 'relation', true, 2, '{"relation_to_object":"equipment"}'::jsonb FROM object_definitions WHERE name = 'equipment_log'
UNION ALL SELECT id, 'action', 'Action', 'select', true, 3, '{"options":["assign","return","move","inspect"]}'::jsonb FROM object_definitions WHERE name = 'equipment_log'
UNION ALL SELECT id, 'from_horse', 'From Horse', 'relation', false, 4, '{"relation_to_object":"horses"}'::jsonb FROM object_definitions WHERE name = 'equipment_log'
UNION ALL SELECT id, 'to_horse', 'To Horse', 'relation', false, 5, '{"relation_to_object":"horses"}'::jsonb FROM object_definitions WHERE name = 'equipment_log'
UNION ALL SELECT id, 'condition', 'Condition', 'select', true, 6, '{"options":["good","worn","damaged"]}'::jsonb FROM object_definitions WHERE name = 'equipment_log'
UNION ALL SELECT id, 'notes', 'Notes', 'text', false, 7, '{}'::jsonb FROM object_definitions WHERE name = 'equipment_log'
UNION ALL SELECT id, 'photos', 'Photos', 'file', false, 8, '{"min":1,"max":5,"maxSizeMB":10}'::jsonb FROM object_definitions WHERE name = 'equipment_log';

-- cleaning (date+time merged into datetime)
INSERT INTO field_definitions (object_id, name, label, type, required, position, metadata)
SELECT id, 'datetime', 'Date & Time', 'datetime', true, 1, '{}'::jsonb FROM object_definitions WHERE name = 'cleaning'
UNION ALL SELECT id, 'cleaned_by', 'Cleaned By', 'relation', true, 2, '{"relation_to_object":"users"}'::jsonb FROM object_definitions WHERE name = 'cleaning'
UNION ALL SELECT id, 'area', 'Area', 'select', true, 3, '{"options":["stable","arena","tack_room","pasture","wash_area"]}'::jsonb FROM object_definitions WHERE name = 'cleaning'
UNION ALL SELECT id, 'status', 'Status', 'select', true, 4, '{"options":["completed","in_progress"],"default":"completed"}'::jsonb FROM object_definitions WHERE name = 'cleaning'
UNION ALL SELECT id, 'notes', 'Notes', 'text', false, 5, '{}'::jsonb FROM object_definitions WHERE name = 'cleaning'
UNION ALL SELECT id, 'photos', 'Photos', 'file', false, 6, '{"min":1,"max":5,"maxSizeMB":10}'::jsonb FROM object_definitions WHERE name = 'cleaning';

-- grass_delivery (date+time merged into datetime)
-- NOTE: spec names this field rumput_type (select: species/cut-type, required) without
-- an enumerated option list anywhere in SPEC.md — renamed to grass_type for English
-- naming consistency; options left empty pending confirmation.
INSERT INTO field_definitions (object_id, name, label, type, required, position, metadata)
SELECT id, 'datetime', 'Date & Time', 'datetime', true, 1, '{}'::jsonb FROM object_definitions WHERE name = 'grass_delivery'
UNION ALL SELECT id, 'delivered_by', 'Delivered By', 'relation', true, 2, '{"relation_to_object":"users"}'::jsonb FROM object_definitions WHERE name = 'grass_delivery'
UNION ALL SELECT id, 'grass_type', 'Grass Type', 'select', true, 3, '{"options":[]}'::jsonb FROM object_definitions WHERE name = 'grass_delivery'
UNION ALL SELECT id, 'quantity', 'Quantity', 'number', true, 4, '{}'::jsonb FROM object_definitions WHERE name = 'grass_delivery'
UNION ALL SELECT id, 'unit', 'Unit', 'select', true, 5, '{"options":["kg"]}'::jsonb FROM object_definitions WHERE name = 'grass_delivery'
UNION ALL SELECT id, 'horse_allocated', 'Horse Allocated', 'relation', false, 6, '{"relation_to_object":"horses"}'::jsonb FROM object_definitions WHERE name = 'grass_delivery'
UNION ALL SELECT id, 'location', 'Location', 'select', false, 7, '{"options":["hay_barn","feeding_area"]}'::jsonb FROM object_definitions WHERE name = 'grass_delivery'
UNION ALL SELECT id, 'photos', 'Photos', 'file', true, 8, '{"min":1,"max":5,"maxSizeMB":10}'::jsonb FROM object_definitions WHERE name = 'grass_delivery';
