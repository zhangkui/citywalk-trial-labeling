import axios, { type AxiosRequestConfig } from 'axios';
import type { ApiResponse, SessionData } from './types';

const baseURL = import.meta.env.VITE_API_BASE_URL ?? '/api';
const sessionKey = 'citywalk.session';

export const api = axios.create({
  baseURL,
  headers: { 'Content-Type': 'application/json' },
});

export const loadSession = (): SessionData | null => {
  const raw = localStorage.getItem(sessionKey);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as SessionData;
  } catch {
    return null;
  }
};

export const saveSession = (session: SessionData | null) => {
  if (!session) {
    localStorage.removeItem(sessionKey);
    return;
  }
  localStorage.setItem(sessionKey, JSON.stringify(session));
};

export const authHeader = () => {
  const session = loadSession();
  return session?.accessToken ? `Bearer ${session.accessToken}` : '';
};

api.interceptors.request.use((config) => {
  const token = authHeader();
  if (token) {
    config.headers = config.headers ?? {};
    config.headers.Authorization = token;
  }
  return config;
});

let refreshing: Promise<string | null> | null = null;

async function refreshAccessToken() {
  if (refreshing) return refreshing;
  const session = loadSession();
  if (!session?.refreshToken) return null;
  refreshing = axios
    .post<ApiResponse<SessionData>>(`${baseURL}/v1/auth/refresh`, { refreshToken: session.refreshToken })
    .then((response) => {
      if (response.data.code !== 0) return null;
      const next = response.data.data;
      saveSession({ ...session, accessToken: next.accessToken, refreshToken: next.refreshToken, expiresIn: next.expiresIn });
      return next.accessToken;
    })
    .catch(() => null)
    .finally(() => {
      refreshing = null;
    });
  return refreshing;
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error?.response?.status === 401 && !error.config?._retry) {
      const token = await refreshAccessToken();
      if (token) {
        error.config._retry = true;
        error.config.headers = error.config.headers ?? {};
        error.config.headers.Authorization = `Bearer ${token}`;
        return api.request(error.config);
      }
    }
    return Promise.reject(error);
  },
);

export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  const response = await api.request<ApiResponse<T>>(config);
  if (response.data.code !== 0) {
    throw new Error(response.data.message || 'request failed');
  }
  return response.data.data;
}

export const get = <T>(url: string, params?: Record<string, unknown>) => request<T>({ url, method: 'GET', params });
export const post = <T>(url: string, data?: unknown) => request<T>({ url, method: 'POST', data });
export const put = <T>(url: string, data?: unknown) => request<T>({ url, method: 'PUT', data });
export const del = <T>(url: string, data?: unknown) => request<T>({ url, method: 'DELETE', data });
