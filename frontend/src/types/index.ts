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

export interface Group {
  ID: number;
  tenant_id: number;
  bot_id: number;
  chat_id: number;
  title: string;
  type: string;
  is_active: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

export interface FeedbackConfig {
  ID: number;
  group_id: number;
  post_to_group: boolean;
  forum_topic_id: number | null;
}

export interface Feedback {
  ID: number;
  tenant_id: number;
  group_id: number;
  message: string;
  admin_only: boolean;
  posted: boolean;
  CreatedAt: string;
}

export interface User {
  id: number;
  tenant_id: number;
  access_token: string;
  name: string;
}
