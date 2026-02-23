import { describe, it, expect, beforeEach, vi } from 'vitest';
import authReducer, {
  setUser,
  clearUser,
  setToken,
  clearToken,
  setRefreshToken,
  clearAuth,
} from './authSlice';
import type { AuthState } from '@/redux/types';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
  length: 0,
  key: vi.fn(),
};
Object.defineProperty(global, 'localStorage', { value: localStorageMock });

describe('authSlice', () => {
  const initialState: AuthState = {
    user: undefined,
    token: undefined,
    refreshToken: undefined,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should return initial state', () => {
    expect(authReducer(undefined, { type: 'unknown' })).toEqual(initialState);
  });

  it('should set user', () => {
    const user = { id: 1, tenant_id: 1, name: 'Test', email: 'test@test.com', role: 'user' };
    const state = authReducer(initialState, setUser(user));
    expect(state.user).toEqual(user);
  });

  it('should clear user', () => {
    const stateWithUser: AuthState = {
      ...initialState,
      user: { id: 1, tenant_id: 1, name: 'Test', email: 'test@test.com', role: 'user' },
    };
    const state = authReducer(stateWithUser, clearUser());
    expect(state.user).toBeUndefined();
  });

  it('should set token and save to localStorage', () => {
    const state = authReducer(initialState, setToken('my-token'));
    expect(state.token).toBe('my-token');
    expect(localStorageMock.setItem).toHaveBeenCalledWith('access_token', 'my-token');
  });

  it('should clear token and remove from localStorage', () => {
    const stateWithToken: AuthState = { ...initialState, token: 'old-token' };
    const state = authReducer(stateWithToken, clearToken());
    expect(state.token).toBeUndefined();
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('access_token');
  });

  it('should set refresh token and save to localStorage', () => {
    const state = authReducer(initialState, setRefreshToken('refresh-token'));
    expect(state.refreshToken).toBe('refresh-token');
    expect(localStorageMock.setItem).toHaveBeenCalledWith('refresh_token', 'refresh-token');
  });

  it('should clear all auth state', () => {
    const fullState: AuthState = {
      user: { id: 1, tenant_id: 1, name: 'Test', email: 'test@test.com', role: 'user' },
      token: 'token',
      refreshToken: 'refresh',
    };
    const state = authReducer(fullState, clearAuth());
    expect(state.user).toBeUndefined();
    expect(state.token).toBeUndefined();
    expect(state.refreshToken).toBeUndefined();
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('access_token');
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('refresh_token');
  });
});
