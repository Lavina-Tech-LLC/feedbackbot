import { api } from '@/api';
import { queryKeys } from '@/api/queryKeys';
import { useMutation, useQuery } from '@tanstack/react-query';
import type { AuthConfig, User } from '@/types';

export const useAuthConfig = () =>
  useQuery({
    queryKey: queryKeys.auth.config(),
    queryFn: () => api.get<never, { data: AuthConfig }>('/auth/config', { noAuth: true }),
  });

export const useExchangeToken = () =>
  useMutation({
    mutationFn: (data: { code: string; redirect_uri: string }) =>
      api.post<never, { data: User }>('/auth/token', data, { noAuth: true }),
  });

export const useRegister = () =>
  useMutation({
    mutationFn: (data: { name: string; email: string; password: string }) =>
      api.post<never, { data: User & { access_token: string } }>('/auth/register', data, { noAuth: true }),
  });

export const useLogin = () =>
  useMutation({
    mutationFn: (data: { email: string; password: string }) =>
      api.post<never, { data: { access_token: string } }>('/auth/login', data, { noAuth: true }),
  });

export const useForgotPassword = () =>
  useMutation({
    mutationFn: () =>
      api.get<never, { data: { url: string } }>('/auth/forgot-password', { noAuth: true }),
  });
