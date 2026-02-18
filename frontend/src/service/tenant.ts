import { api } from '@/api';
import { queryKeys } from '@/api/queryKeys';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { Tenant } from '@/types';

export const useGetTenant = (id: string) =>
  useQuery({
    queryKey: queryKeys.tenants.byId(id),
    queryFn: () => api.get<never, { data: Tenant }>(`/tenants/${id}`),
    enabled: !!id,
  });

export const useCreateTenant = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { name: string; slug: string }) =>
      api.post<never, { data: Tenant }>('/tenants', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.tenants.all() });
    },
  });
};
