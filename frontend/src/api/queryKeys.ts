export const queryKeys = {
  tenants: {
    all: () => ['tenants'] as const,
    byId: (id: string) => ['tenants', id] as const,
  },
  bots: {
    all: () => ['bots'] as const,
    byId: (id: string) => ['bots', id] as const,
  },
  groups: {
    all: () => ['groups'] as const,
    byTenant: (tenantId: string) => ['groups', 'tenant', tenantId] as const,
    byId: (id: string) => ['groups', id] as const,
  },
  feedbacks: {
    byGroup: (groupId: string, adminOnly?: string, page?: number, dateFrom?: string, dateTo?: string, search?: string) =>
      ['feedbacks', groupId, adminOnly, page, dateFrom, dateTo, search] as const,
  },
};
