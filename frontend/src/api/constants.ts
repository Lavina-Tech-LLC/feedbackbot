const HOSTNAME_MAP: Record<string, string> = {
  'feedbackbot.lavina.tech': 'https://feedbackbot-api.lavina.tech',
  'stage-feedbackbot.lavina.tech': 'https://stage-feedbackbot-api.lavina.tech',
  localhost: '/api',
};

export const getApiBaseUrl = (): string => {
  return HOSTNAME_MAP[window.location.hostname] ?? '/api';
};

export const api_constants = {
  baseUrl: getApiBaseUrl(),
};
