import { api } from '@/api';
import { useMutation } from '@tanstack/react-query';
import type { User } from '@/types';

interface AuthResponse {
  data: {
    access_token: string;
    refresh_token: string;
    user: User;
  };
}

export const useRegister = () =>
  useMutation({
    mutationFn: (data: { name: string; email: string; password: string }) =>
      api.post<never, AuthResponse>('/auth/register', data, { noAuth: true }),
  });

export const useLogin = () =>
  useMutation({
    mutationFn: (data: { email: string; password: string }) =>
      api.post<never, AuthResponse>('/auth/login', data, { noAuth: true }),
  });

export const useRefreshToken = () =>
  useMutation({
    mutationFn: (data: { refresh_token: string }) =>
      api.post<never, AuthResponse>('/auth/refresh', data, { noAuth: true }),
  });

export const useMe = () =>
  useMutation({
    mutationFn: () => api.get<never, { data: User }>('/auth/me'),
  });
