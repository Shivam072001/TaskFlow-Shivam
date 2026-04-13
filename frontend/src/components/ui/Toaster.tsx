import { Toaster as SonnerToaster } from 'sonner';

export function Toaster() {
  return (
    <SonnerToaster
      position="bottom-right"
      gap={8}
      toastOptions={{
        unstyled: true,
        classNames: {
          toast:
            'group flex w-[360px] items-start gap-3 rounded-xl border bg-card px-4 py-3.5 shadow-lg backdrop-blur-sm',
          title: 'text-sm font-semibold text-foreground',
          description: 'mt-0.5 text-xs text-muted-foreground leading-relaxed',
          actionButton:
            'ml-auto shrink-0 rounded-md bg-primary px-3 py-1.5 text-xs font-medium text-primary-foreground hover:bg-primary/90',
          cancelButton:
            'ml-auto shrink-0 rounded-md px-3 py-1.5 text-xs font-medium text-muted-foreground hover:bg-muted',
          success: 'border-emerald-500/30 bg-emerald-50/90 dark:bg-emerald-950/50 dark:border-emerald-500/20',
          error: 'border-destructive/30 bg-red-50/90 dark:bg-red-950/50 dark:border-destructive/20',
          warning: 'border-amber-500/30 bg-amber-50/90 dark:bg-amber-950/50 dark:border-amber-500/20',
          info: 'border-blue-500/30 bg-blue-50/90 dark:bg-blue-950/50 dark:border-blue-500/20',
          closeButton:
            'absolute -right-1.5 -top-1.5 rounded-full border border-border bg-card p-0.5 text-muted-foreground opacity-0 shadow-sm transition-opacity group-hover:opacity-100 hover:text-foreground',
        },
      }}
    />
  );
}
