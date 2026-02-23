export interface Tenant {
  id: number;
  name: string;
  slug: string;
  createdAt: string;
  updatedAt: string;
}

export interface Bot {
  id: number;
  tenant_id: number;
  bot_username: string;
  bot_name: string;
  verified: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Group {
  id: number;
  tenant_id: number;
  bot_id: number;
  chat_id: number;
  title: string;
  type: string;
  is_active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface FeedbackConfig {
  id: number;
  group_id: number;
  post_to_group: boolean;
  forum_topic_id: number | null;
}

export interface Feedback {
  id: number;
  tenant_id: number;
  group_id: number;
  message: string;
  admin_only: boolean;
  posted: boolean;
  createdAt: string;
}

export interface User {
  id: number;
  tenant_id: number;
  name: string;
  email: string;
  role: string;
}
