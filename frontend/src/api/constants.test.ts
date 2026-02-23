import { describe, it, expect, vi } from 'vitest';
import { getApiBaseUrl } from './constants';

describe('getApiBaseUrl', () => {
  it('should return production URL for feedbackbot.lavina.tech', () => {
    vi.stubGlobal('location', { hostname: 'feedbackbot.lavina.tech' });
    expect(getApiBaseUrl()).toBe('https://feedbackbot-api.lavina.tech');
    vi.unstubAllGlobals();
  });

  it('should return staging URL for stage-feedbackbot.lavina.tech', () => {
    vi.stubGlobal('location', { hostname: 'stage-feedbackbot.lavina.tech' });
    expect(getApiBaseUrl()).toBe('https://stage-feedbackbot-api.lavina.tech');
    vi.unstubAllGlobals();
  });

  it('should return /api for localhost', () => {
    vi.stubGlobal('location', { hostname: 'localhost' });
    expect(getApiBaseUrl()).toBe('/api');
    vi.unstubAllGlobals();
  });

  it('should return /api for unknown hostnames', () => {
    vi.stubGlobal('location', { hostname: 'unknown.example.com' });
    expect(getApiBaseUrl()).toBe('/api');
    vi.unstubAllGlobals();
  });
});
