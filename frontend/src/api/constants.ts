const HOSTNAME_MAP: Record<string, string> = {
  'feedbackbot.lavina.tech': 'https://api-feedbackbot.lavina.tech',
  'stage-feedbackbot.lavina.tech': 'https://api-stage-feedbackbot.lavina.tech',
  localhost: '/api',
};

export const getApiBaseUrl = (): string => {
  return HOSTNAME_MAP[window.location.hostname] ?? '/api';
};

export const api_constants = {
  baseUrl: getApiBaseUrl(),
};
