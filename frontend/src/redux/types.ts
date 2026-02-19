import type { User } from '@/types';

export interface AuthState {
  user?: User;
  token?: string;
}
