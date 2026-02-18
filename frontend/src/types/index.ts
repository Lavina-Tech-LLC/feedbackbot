export interface Tenant {
  ID: number;
  name: string;
  slug: string;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface Bot {
  ID: number;
  tenant_id: number;
  bot_username: string;
  bot_name: string;
  verified: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface User {
  id: number;
  access_token: string;
  name: string;
}
