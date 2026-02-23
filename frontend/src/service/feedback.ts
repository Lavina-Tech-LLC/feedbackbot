import { api } from '@/api';
import { api_constants } from '@/api/constants';
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

export interface FeedbackParams {
  groupId: string;
  adminOnly?: string;
  page?: number;
  limit?: number;
  dateFrom?: string;
  dateTo?: string;
  search?: string;
}

function buildSearchParams(params: FeedbackParams): URLSearchParams {
  const searchParams = new URLSearchParams({ group_id: params.groupId });
  if (params.adminOnly) searchParams.set('admin_only', params.adminOnly);
  if (params.page) searchParams.set('page', String(params.page));
  if (params.limit) searchParams.set('limit', String(params.limit));
  if (params.dateFrom) searchParams.set('date_from', params.dateFrom);
  if (params.dateTo) searchParams.set('date_to', params.dateTo);
  if (params.search) searchParams.set('search', params.search);
  return searchParams;
}

export const useGetFeedbacks = (params: FeedbackParams) =>
  useQuery({
    queryKey: queryKeys.feedbacks.byGroup(params.groupId, params.adminOnly, params.page, params.dateFrom, params.dateTo, params.search),
    queryFn: () => {
      const searchParams = buildSearchParams(params);
      return api.get<never, FeedbackListResponse>(`/feedbacks?${searchParams.toString()}`);
    },
    enabled: !!params.groupId,
  });

export function getExportCsvUrl(params: FeedbackParams): string {
  const searchParams = buildSearchParams(params);
  return `${api_constants.baseUrl}/feedbacks/export?${searchParams.toString()}`;
}
