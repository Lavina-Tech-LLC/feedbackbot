import { api } from '@/api';
import { queryKeys } from '@/api/queryKeys';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { Bot } from '@/types';

export const useGetBots = () =>
  useQuery({
    queryKey: queryKeys.bots.all(),
    queryFn: () => api.get<never, { data: Bot[] }>('/bots'),
  });

export const useGetBot = (id: string) =>
  useQuery({
    queryKey: queryKeys.bots.byId(id),
    queryFn: () => api.get<never, { data: Bot }>(`/bots/${id}`),
    enabled: !!id,
  });

export const useCreateBot = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { tenant_id: number; token: string }) =>
      api.post<never, { data: Bot }>('/bots', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.bots.all() });
    },
  });
};

export const useDeleteBot = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => api.delete(`/bots/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.bots.all() });
    },
  });
};
