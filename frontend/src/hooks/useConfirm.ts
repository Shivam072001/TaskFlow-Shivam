import { useContext } from 'react';
import { ConfirmContext, type ConfirmFn } from '@core/context/ConfirmContext';

export function useConfirm(): ConfirmFn {
  const fn = useContext(ConfirmContext);
  if (!fn) throw new Error('useConfirm must be used within ConfirmDialogProvider');
  return fn;
}
