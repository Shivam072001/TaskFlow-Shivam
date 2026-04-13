import type React from 'react';
import { useState, useRef, useEffect, useCallback, type ReactNode } from 'react';
import { ChevronDown, Check } from 'lucide-react';
import { cn } from '@utils/cn';

export interface SelectOption {
  value: string;
  label: string;
  icon?: ReactNode;
}

interface SelectProps {
  value: string;
  onChange: (value: string) => void;
  options: SelectOption[];
  placeholder?: string;
  disabled?: boolean;
  className?: string;
}

export function Select({ value, onChange, options, placeholder = 'Select…', disabled, className }: SelectProps) {
  const [open, setOpen] = useState(false);
  const [focusIndex, setFocusIndex] = useState(-1);
  const containerRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLUListElement>(null);

  const selected = options.find((o) => o.value === value);

  const close = useCallback(() => {
    setOpen(false);
    setFocusIndex(-1);
  }, []);

  useEffect(() => {
    function handleOutside(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) close();
    }
    document.addEventListener('mousedown', handleOutside);
    return () => document.removeEventListener('mousedown', handleOutside);
  }, [close]);

  useEffect(() => {
    if (open && focusIndex >= 0) {
      const items = listRef.current?.querySelectorAll('[role="option"]');
      (items?.[focusIndex] as HTMLElement)?.scrollIntoView({ block: 'nearest' });
    }
  }, [focusIndex, open]);

  function handleKeyDown(e: React.KeyboardEvent<HTMLButtonElement>) {
    if (disabled) return;

    switch (e.key) {
      case 'Enter':
      case ' ':
        e.preventDefault();
        if (!open) {
          setOpen(true);
          setFocusIndex(options.findIndex((o) => o.value === value));
        } else if (focusIndex >= 0) {
          onChange(options[focusIndex].value);
          close();
        }
        break;
      case 'ArrowDown':
        e.preventDefault();
        if (!open) {
          setOpen(true);
          setFocusIndex(0);
        } else {
          setFocusIndex((i) => Math.min(i + 1, options.length - 1));
        }
        break;
      case 'ArrowUp':
        e.preventDefault();
        setFocusIndex((i) => Math.max(i - 1, 0));
        break;
      case 'Escape':
        close();
        break;
    }
  }

  return (
    <div ref={containerRef} className={cn('relative', className)}>
      <button
        type="button"
        role="combobox"
        aria-expanded={open}
        aria-haspopup="listbox"
        disabled={disabled}
        onClick={() => !disabled && setOpen((o) => !o)}
        onKeyDown={handleKeyDown}
        className={cn(
          'flex w-full items-center justify-between gap-2 rounded-md border border-border bg-card px-3 py-2 text-sm',
          'transition-colors hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
          disabled && 'cursor-not-allowed opacity-50',
        )}
      >
        <span className={cn('flex items-center gap-2 truncate', !selected && 'text-muted-foreground')}>
          {selected?.icon}
          {selected?.label ?? placeholder}
        </span>
        <ChevronDown className={cn('h-4 w-4 shrink-0 text-muted-foreground transition-transform', open && 'rotate-180')} />
      </button>

      <div
        className={cn(
          'absolute z-50 mt-1 w-full rounded-md border border-border bg-card shadow-lg',
          'transition-all duration-150',
          open ? 'pointer-events-auto translate-y-0 opacity-100' : 'pointer-events-none -translate-y-1 opacity-0',
        )}
      >
        <ul
          ref={listRef}
          role="listbox"
          aria-activedescendant={focusIndex >= 0 ? `select-opt-${focusIndex}` : undefined}
          className="max-h-60 overflow-auto p-1"
        >
          {options.map((opt, i) => {
            const isSelected = opt.value === value;
            const isFocused = i === focusIndex;
            return (
              <li
                key={opt.value}
                id={`select-opt-${i}`}
                role="option"
                aria-selected={isSelected}
                className={cn(
                  'flex cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm transition-colors',
                  isFocused && 'bg-muted',
                  isSelected && 'font-medium',
                )}
                onMouseEnter={() => setFocusIndex(i)}
                onClick={() => {
                  onChange(opt.value);
                  close();
                }}
              >
                {opt.icon}
                <span className="flex-1 truncate">{opt.label}</span>
                {isSelected && <Check className="h-4 w-4 text-primary" />}
              </li>
            );
          })}
        </ul>
      </div>
    </div>
  );
}
