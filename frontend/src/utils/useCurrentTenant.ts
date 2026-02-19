import { useSelector } from 'react-redux';
import type { RootState } from '@/redux/store';

export function useCurrentTenant(): string {
  const user = useSelector((state: RootState) => state.auth.user);
  return user?.tenant_id ? String(user.tenant_id) : '';
}
