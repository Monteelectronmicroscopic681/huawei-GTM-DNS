interface ConfirmOptions {
  title?: string;
  description?: string;
  confirmText?: string;
  cancelText?: string;
  variant?: 'destructive' | 'default';
}

let resolveRef: ((value: boolean) => void) | null = null;

export const confirmState = $state({
  open: false,
  title: '',
  description: '',
  confirmText: '',
  cancelText: '',
  variant: 'destructive' as 'destructive' | 'default',
});

export function confirm(options: ConfirmOptions): Promise<boolean> {
  return new Promise((resolve) => {
    confirmState.title = options.title ?? '';
    confirmState.description = options.description ?? '';
    confirmState.confirmText = options.confirmText ?? '';
    confirmState.cancelText = options.cancelText ?? '';
    confirmState.variant = options.variant ?? 'destructive';
    resolveRef = resolve;
    confirmState.open = true;
  });
}

// Called by ConfirmDialog on Action click or dialog close.
// Action onclick fires BEFORE onOpenChange, so resolveRef is nulled
// before the close handler runs — no double-resolve.
export function resolveConfirm(value: boolean) {
  resolveRef?.(value);
  resolveRef = null;
}
