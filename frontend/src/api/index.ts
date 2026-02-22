import axios from 'axios';
import { api_constants } from './constants';
import { store } from '@/redux/store';
import { clearAuth, setToken, setRefreshToken } from '@/redux/slices';

declare module 'axios' {
  export interface AxiosRequestConfig {
    noAuth?: boolean;
    _retry?: boolean;
  }
}

export const api = axios.create({
  baseURL: api_constants.baseUrl,
});

api.interceptors.request.use((config) => {
  if (config.method === 'get' && !config.signal) {
    const controller = new AbortController();
    setTimeout(() => controller.abort(), 15000);
    config.signal = controller.signal;
  }

  if (config.noAuth === true) {
    delete config.headers.Authorization;
  } else {
    const token = store.getState().auth.token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

let isRefreshing = false;
let failedQueue: { resolve: (token: string) => void; reject: (err: unknown) => void }[] = [];

function processQueue(error: unknown, token: string | null) {
  failedQueue.forEach((p) => {
    if (token) p.resolve(token);
    else p.reject(error);
  });
  failedQueue = [];
}

api.interceptors.response.use(
  (response) => response.data,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.noAuth) {
      const refreshToken = store.getState().auth.refreshToken;
      if (!refreshToken) {
        store.dispatch(clearAuth());
        window.location.href = '/login';
        return Promise.reject(error);
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then((token) => {
          originalRequest.headers.Authorization = `Bearer ${token}`;
          originalRequest._retry = true;
          return api(originalRequest);
        });
      }

      isRefreshing = true;
      originalRequest._retry = true;

      try {
        const res = await axios.post(
          `${api_constants.baseUrl}/auth/refresh`,
          { refresh_token: refreshToken },
        );
        const { access_token, refresh_token } = res.data.data;
        store.dispatch(setToken(access_token));
        store.dispatch(setRefreshToken(refresh_token));
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        processQueue(null, access_token);
        return api(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        store.dispatch(clearAuth());
        window.location.href = '/login';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  },
);
