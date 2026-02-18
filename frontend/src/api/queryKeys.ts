export const queryKeys = {
  tenants: {
    all: () => ['tenants'] as const,
    byId: (id: string) => ['tenants', id] as const,
  },
  bots: {
    all: () => ['bots'] as const,
    byId: (id: string) => ['bots', id] as const,
  },
};
