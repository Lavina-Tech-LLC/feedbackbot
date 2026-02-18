import axios from 'axios';
import { api_constants } from './constants';
import { store } from '@/redux/store';
import { clearUser } from '@/redux/slices';
import Cookies from 'js-cookie';

declare module 'axios' {
  export interface AxiosRequestConfig {
    noAuth?: boolean;
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
    const session = Cookies.get('session');
    if (session) {
      const user = JSON.parse(session);
      config.headers.Authorization = user?.access_token;
    }
  }
  return config;
});

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      store.dispatch(clearUser());
    }
    return Promise.reject(error);
  },
);
