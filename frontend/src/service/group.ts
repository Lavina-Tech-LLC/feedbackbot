import { api } from '@/api';
import { queryKeys } from '@/api/queryKeys';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { Group, FeedbackConfig } from '@/types';

export const useGetGroups = (tenantId: string) =>
  useQuery({
    queryKey: queryKeys.groups.byTenant(tenantId),
    queryFn: () => api.get<never, { data: Group[] }>(`/groups?tenant_id=${tenantId}`),
    enabled: !!tenantId,
  });

export const useUpdateGroup = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: { is_active?: boolean } }) =>
      api.patch<never, { data: Group }>(`/groups/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] });
    },
  });
};

export const useUpdateGroupConfig = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      groupId,
      data,
    }: {
      groupId: number;
      data: { post_to_group?: boolean; forum_topic_id?: number };
    }) => api.patch<never, { data: FeedbackConfig }>(`/groups/${groupId}/config`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['groups'] });
    },
  });
};
