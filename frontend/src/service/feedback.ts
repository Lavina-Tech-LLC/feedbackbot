import { api } from '@/api';
import { queryKeys } from '@/api/queryKeys';
import { useQuery } from '@tanstack/react-query';
import type { Feedback } from '@/types';

interface FeedbackListResponse {
  data: {
    data: Feedback[];
    total: number;
    page: number;
    limit: number;
  };
}

export const useGetFeedbacks = (params: {
  groupId: string;
  adminOnly?: string;
  page?: number;
  limit?: number;
}) =>
  useQuery({
    queryKey: queryKeys.feedbacks.byGroup(params.groupId, params.adminOnly, params.page),
    queryFn: () => {
      const searchParams = new URLSearchParams({ group_id: params.groupId });
      if (params.adminOnly) searchParams.set('admin_only', params.adminOnly);
      if (params.page) searchParams.set('page', String(params.page));
      if (params.limit) searchParams.set('limit', String(params.limit));
      return api.get<never, FeedbackListResponse>(`/feedbacks?${searchParams.toString()}`);
    },
    enabled: !!params.groupId,
  });
