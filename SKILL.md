# Horse Stable Management — Schema-Driven Architecture Skill

**How to extend the system without writing code: Add schemas, not components.**

---

## Overview

This system is **schema-driven**: objects, fields, and dashboards are defined in the database, not code. Once the generic form builder, table renderer, and dashboard widget renderer are built (Phase 1), adding new activity types or reports is configuration work, not development work.

**Configuration (no code needed):**
- Add a new activity type (Farrier visit, Vet checkup, Supplement log)
- Add a new report widget (top 10 horses by training hours, equipment maintenance log)
- Customize field validation or options

**Development (requires code):**
- Add a new field type (phone, location, etc.)
- Add a new widget type (Gantt, heatmap)
- Change form layout or business logic

---

## How to: Add a New Activity Type (No Code)

### Step 1: Create Object Definition

```sql
INSERT INTO object_definitions (name, label, plural_label, icon, description) VALUES
  ('farrier_visit', 'Farrier Visit', 'Farrier Visits', '🐴', 'Farrier hoof trimming log');
```

### Step 2: Define Fields

```sql
INSERT INTO field_definitions (object_id, name, label, type, required, position, metadata) VALUES
  ((SELECT id FROM object_definitions WHERE name='farrier_visit'), 'date', 'Date', 'date', true, 1, '{}'),
  ((SELECT id FROM object_definitions WHERE name='farrier_visit'), 'horse_id', 'Horse', 'relation', true, 2, '{"relation_to_object": "horses"}'),
  ((SELECT id FROM object_definitions WHERE name='farrier_visit'), 'farrier_name', 'Farrier Name', 'text', true, 3, '{}'),
  ((SELECT id FROM object_definitions WHERE name='farrier_visit'), 'services', 'Services', 'multiselect', true, 4, '{"options": ["trim", "shoe", "balance"]}'),
  ((SELECT id FROM object_definitions WHERE name='farrier_visit'), 'notes', 'Notes', 'text', false, 5, '{}'),
  ((SELECT id FROM object_definitions WHERE name='farrier_visit'), 'photos', 'Photos', 'file', false, 6, '{}');
```

**Available field types:** text, number, date, datetime, select, multiselect, relation, file, boolean

### Step 3: Done

Frontend automatically renders form + table. No code change needed.

---

## How to: Add a Dashboard Widget (No Code)

```sql
INSERT INTO dashboard_definitions (object_id, name, label, widgets) VALUES
  (
    (SELECT id FROM object_definitions WHERE name='farrier_visit'),
    'farrier_summary',
    'Farrier Summary',
    '[
      {
        "id": "total_visits",
        "type": "stat",
        "title": "Total Visits",
        "query": "SELECT COUNT(*) as count FROM records WHERE object_id = (SELECT id FROM object_definitions WHERE name=''farrier_visit'') AND deleted_at IS NULL",
        "metric": "count"
      }
    ]'
  );
```

**Widget types:** stat, chart (bar/pie/line), table, timeline

---

## Best Practices

1. **Field naming**: Use snake_case (`horse_id`, not `hid`)
2. **Relations**: Use relation field type with `{"relation_to_object": "table"}`
3. **Validation**: Positive constraints (`maxLength: 255`)
4. **Dashboard queries**: Keep simple (SUM, COUNT, GROUP BY)

---

**Phase 1**: 6 activity types + 4 dashboards seeded in SQL  
**Phase 2**: Admin UI to create/edit schemas without SQL

**Details**: See SPEC.md, CONSTRAINTS.md, CLAUDE.md, tasks/plan.md  
**Support**: admin@eri.co.id  
**Author**: nurafandito (nurafandito@gmail.com)  
**Organization**: eri-co-id
