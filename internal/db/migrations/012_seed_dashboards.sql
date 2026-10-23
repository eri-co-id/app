-- Seed one dashboard per Phase 1 activity type (SPEC.md §7 dashboard example,
-- reimplemented as structured widget descriptors — see internal/dashboard —
-- rather than a stored raw SQL string, to avoid executing admin-authored SQL
-- directly against the database (CWE-89).

INSERT INTO dashboard_definitions (object_id, name, label, widgets)
SELECT id, 'riding_summary', 'Riding Summary', '[
    {"id":"total_rides","type":"stat","title":"Total Rides","metric":"count"},
    {"id":"fees_by_lesson","type":"chart","title":"Fees by Lesson Code","metric":"sum","field":"fee","group_by":"code_lesson","chart_type":"pie"},
    {"id":"avg_duration","type":"stat","title":"Average Duration (hours)","metric":"avg","field":"duration_hours"}
]'::jsonb
FROM object_definitions WHERE name = 'riding';

INSERT INTO dashboard_definitions (object_id, name, label, widgets)
SELECT id, 'training_summary', 'Training Summary', '[
    {"id":"total_trainings","type":"stat","title":"Total Trainings","metric":"count"},
    {"id":"fees_by_lesson","type":"chart","title":"Fees by Lesson Code","metric":"sum","field":"fee","group_by":"code_lesson","chart_type":"pie"}
]'::jsonb
FROM object_definitions WHERE name = 'training';

INSERT INTO dashboard_definitions (object_id, name, label, widgets)
SELECT id, 'horse_activity_summary', 'Horse Activity Summary', '[
    {"id":"total_activities","type":"stat","title":"Total Activities","metric":"count"},
    {"id":"activities_by_type","type":"chart","title":"Activities by Type","metric":"count","group_by":"activity_type","chart_type":"bar"}
]'::jsonb
FROM object_definitions WHERE name = 'horse_activity';

INSERT INTO dashboard_definitions (object_id, name, label, widgets)
SELECT id, 'consumable_balance', 'Consumable Balance', '[
    {"id":"total_in","type":"stat","title":"Total Stock In","metric":"sum","field":"quantity_in"},
    {"id":"total_out","type":"stat","title":"Total Stock Out","metric":"sum","field":"quantity_out"},
    {"id":"by_item_type","type":"chart","title":"Logs by Item Type","metric":"count","group_by":"item_type","chart_type":"bar"}
]'::jsonb
FROM object_definitions WHERE name = 'consumable_inventory';

INSERT INTO dashboard_definitions (object_id, name, label, widgets)
SELECT id, 'equipment_assignment', 'Equipment Assignment', '[
    {"id":"total_logs","type":"stat","title":"Total Equipment Logs","metric":"count"},
    {"id":"by_action","type":"chart","title":"Logs by Action","metric":"count","group_by":"action","chart_type":"bar"}
]'::jsonb
FROM object_definitions WHERE name = 'equipment_log';

INSERT INTO dashboard_definitions (object_id, name, label, widgets)
SELECT id, 'cleaning_summary', 'Cleaning Summary', '[
    {"id":"total_cleanings","type":"stat","title":"Total Cleanings","metric":"count"},
    {"id":"by_area","type":"chart","title":"Cleanings by Area","metric":"count","group_by":"area","chart_type":"bar"}
]'::jsonb
FROM object_definitions WHERE name = 'cleaning';

INSERT INTO dashboard_definitions (object_id, name, label, widgets)
SELECT id, 'grass_delivery_summary', 'Grass Delivery Summary', '[
    {"id":"total_deliveries","type":"stat","title":"Total Deliveries","metric":"count"},
    {"id":"total_quantity","type":"stat","title":"Total Quantity Delivered","metric":"sum","field":"quantity"}
]'::jsonb
FROM object_definitions WHERE name = 'grass_delivery';
