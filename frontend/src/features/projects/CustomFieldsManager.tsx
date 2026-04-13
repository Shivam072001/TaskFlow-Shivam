import { useState } from 'react';
import { Plus, Trash2, Loader2, Settings2 } from 'lucide-react';
import type { CustomFieldDefinition } from '@/types';
import { useCustomFieldDefinitions, useCreateCustomField, useDeleteCustomField } from '@hooks/useCustomFields';
import { useConfirm } from '@hooks/useConfirm';
import { Select } from '@components/ui/Select';
import { Badge } from '@components/ui/Badge';


interface Props {
  projectId: string;
  canManage: boolean;
}

const fieldTypeOptions = [
  { value: 'text', label: 'Text' },
  { value: 'number', label: 'Number' },
  { value: 'select', label: 'Select' },
];

function FieldRow({
  field,
  canManage,
  onDelete,
  isDeleting,
}: {
  field: CustomFieldDefinition;
  canManage: boolean;
  onDelete: (id: string) => void;
  isDeleting: boolean;
}) {
  const confirm = useConfirm();

  async function handleDelete() {
    const ok = await confirm({
      title: `Delete "${field.name}"`,
      description: 'Existing values for this field will be permanently lost.',
      confirmLabel: 'Delete field',
      variant: 'danger',
    });
    if (ok) onDelete(field.id);
  }

  return (
    <div className="flex items-center justify-between gap-4 rounded-md border border-border bg-card px-4 py-3">
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-foreground">{field.name}</span>
          <Badge variant="outline">{field.field_type}</Badge>
          {field.required && <Badge variant="destructive">Required</Badge>}
        </div>
        {field.field_type === 'select' && field.options.length > 0 && (
          <p className="mt-1 text-xs text-muted-foreground">
            Options: {field.options.join(', ')}
          </p>
        )}
      </div>
      {canManage && (
        <button
          type="button"
          onClick={handleDelete}
          disabled={isDeleting}
          className="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive disabled:opacity-50"
        >
          {isDeleting ? <Loader2 className="h-4 w-4 animate-spin" /> : <Trash2 className="h-4 w-4" />}
        </button>
      )}
    </div>
  );
}

export function CustomFieldsManager({ projectId, canManage }: Props) {
  const { data: fields, isLoading } = useCustomFieldDefinitions(projectId);
  const createField = useCreateCustomField(projectId);
  const deleteField = useDeleteCustomField(projectId);

  const [adding, setAdding] = useState(false);
  const [name, setName] = useState('');
  const [fieldType, setFieldType] = useState<'text' | 'number' | 'select'>('text');
  const [optionsRaw, setOptionsRaw] = useState('');
  const [required, setRequired] = useState(false);

  function resetForm() {
    setName('');
    setFieldType('text');
    setOptionsRaw('');
    setRequired(false);
    setAdding(false);
  }

  function handleCreate() {
    const options =
      fieldType === 'select'
        ? optionsRaw
            .split(',')
            .map((s) => s.trim())
            .filter(Boolean)
        : undefined;

    createField.mutate({ name, field_type: fieldType, options, required }, { onSuccess: resetForm });
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <Loader2 className="h-5 w-5 animate-spin" />
      </div>
    );
  }

  const list = fields ?? [];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h3 className="flex items-center gap-2 text-sm font-semibold text-foreground">
          <Settings2 className="h-4 w-4" />
          Custom Fields
          {list.length > 0 && (
            <span className="rounded-full bg-muted px-2 py-0.5 text-xs font-normal text-muted-foreground">
              {list.length}
            </span>
          )}
        </h3>
        {canManage && !adding && (
          <button
            type="button"
            onClick={() => setAdding(true)}
            className="inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground hover:bg-primary/90"
          >
            <Plus className="h-3.5 w-3.5" />
            Add Field
          </button>
        )}
      </div>

      {adding && (
        <div className="space-y-3 rounded-lg border border-border bg-card p-4">
          <div>
            <label className="mb-1 block text-xs font-medium text-foreground">Field Name</label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. Story Points"
              className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
            />
          </div>

          <div>
            <label className="mb-1 block text-xs font-medium text-foreground">Field Type</label>
            <Select
              value={fieldType}
              onChange={(v) => setFieldType(v as 'text' | 'number' | 'select')}
              options={fieldTypeOptions}
            />
          </div>

          {fieldType === 'select' && (
            <div>
              <label className="mb-1 block text-xs font-medium text-foreground">Options (comma-separated)</label>
              <input
                type="text"
                value={optionsRaw}
                onChange={(e) => setOptionsRaw(e.target.value)}
                placeholder="e.g. S, M, L, XL"
                className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring"
              />
            </div>
          )}

          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={required}
              onChange={(e) => setRequired(e.target.checked)}
              className="h-4 w-4 rounded border-border text-primary focus:ring-ring"
            />
            <span className="text-sm text-foreground">Required</span>
          </label>

          <div className="flex gap-2 pt-1">
            <button
              type="button"
              onClick={handleCreate}
              disabled={!name.trim() || createField.isPending || (fieldType === 'select' && !optionsRaw.trim())}
              className="inline-flex items-center gap-1.5 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50"
            >
              {createField.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Plus className="h-4 w-4" />
              )}
              Create
            </button>
            <button
              type="button"
              onClick={resetForm}
              className="rounded-md px-4 py-2 text-sm font-medium text-muted-foreground hover:bg-muted"
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {list.length === 0 && !adding ? (
        <p className="py-6 text-center text-sm text-muted-foreground">
          No custom fields defined.{canManage ? ' Add one to get started.' : ''}
        </p>
      ) : (
        <div className="space-y-2">
          {list.map((field) => (
            <FieldRow
              key={field.id}
              field={field}
              canManage={canManage}
              onDelete={(id) => deleteField.mutate(id)}
              isDeleting={deleteField.isPending}
            />
          ))}
        </div>
      )}
    </div>
  );
}
