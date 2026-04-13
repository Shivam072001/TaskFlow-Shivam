CREATE TABLE custom_field_definitions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    field_type TEXT NOT NULL CHECK (field_type IN ('text', 'number', 'select')),
    options    JSONB DEFAULT '[]',
    required   BOOLEAN NOT NULL DEFAULT FALSE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_custom_fields_project ON custom_field_definitions (project_id);

CREATE TABLE custom_field_values (
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id  UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    field_id UUID NOT NULL REFERENCES custom_field_definitions(id) ON DELETE CASCADE,
    value    TEXT NOT NULL DEFAULT '',
    UNIQUE (task_id, field_id)
);

CREATE INDEX idx_custom_field_values_task ON custom_field_values (task_id);
