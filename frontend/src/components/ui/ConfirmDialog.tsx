import { useCallback, useRef, useState, type ReactNode } from 'react';
import { AlertTriangle, Trash2, X } from 'lucide-react';
import { cn } from '@utils/cn';
import { ConfirmContext, type ConfirmFn, type ConfirmOptions } from '@core/context/ConfirmContext';

export function ConfirmDialogProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<(ConfirmOptions & { open: boolean }) | null>(null);
  const resolveRef = useRef<((value: boolean) => void) | null>(null);

  const confirm: ConfirmFn = useCallback((options) => {
    return new Promise<boolean>((resolve) => {
      resolveRef.current = resolve;
      setState({ ...options, open: true });
    });
  }, []);

  const handleClose = useCallback((result: boolean) => {
    resolveRef.current?.(result);
    resolveRef.current = null;
    setState(null);
  }, []);

  const variant = state?.variant ?? 'default';

  const iconBg = {
    danger: 'bg-destructive/10 text-destructive',
    warning: 'bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400',
    default: 'bg-primary/10 text-primary',
  }[variant];

  const confirmBtnClass = {
    danger: 'bg-destructive text-white hover:bg-destructive/90',
    warning: 'bg-amber-600 text-white hover:bg-amber-700 dark:bg-amber-600 dark:hover:bg-amber-500',
    default: 'bg-primary text-primary-foreground hover:bg-primary/90',
  }[variant];

  const Icon = variant === 'danger' ? Trash2 : AlertTriangle;

  return (
    <ConfirmContext.Provider value={confirm}>
      {children}

      {state?.open && (
        <div
          className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-[2px] animate-in fade-in duration-150"
          onClick={() => handleClose(false)}
        >
          <div
            className="w-full max-w-sm rounded-xl border border-border bg-card p-6 shadow-xl animate-in zoom-in-95 duration-150"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-start gap-4">
              <div className={cn('flex h-10 w-10 shrink-0 items-center justify-center rounded-full', iconBg)}>
                <Icon className="h-5 w-5" />
              </div>
              <div className="flex-1">
                <div className="flex items-start justify-between">
                  <h3 className="text-base font-semibold text-foreground">{state.title}</h3>
                  <button
                    onClick={() => handleClose(false)}
                    className="rounded-md p-0.5 text-muted-foreground hover:text-foreground"
                  >
                    <X className="h-4 w-4" />
                  </button>
                </div>
                {state.description && (
                  <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">{state.description}</p>
                )}
              </div>
            </div>

            <div className="mt-6 flex justify-end gap-2">
              <button
                type="button"
                onClick={() => handleClose(false)}
                className="rounded-lg px-4 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted"
              >
                {state.cancelLabel ?? 'Cancel'}
              </button>
              <button
                type="button"
                onClick={() => handleClose(true)}
                className={cn('rounded-lg px-4 py-2 text-sm font-medium transition-colors', confirmBtnClass)}
              >
                {state.confirmLabel ?? 'Confirm'}
              </button>
            </div>
          </div>
        </div>
      )}
    </ConfirmContext.Provider>
  );
}
