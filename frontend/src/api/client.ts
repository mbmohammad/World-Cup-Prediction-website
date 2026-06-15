import axios from 'axios';

const baseURL = import.meta.env.VITE_API_URL ?? '/api/v1';

export const api = axios.create({ baseURL });

const ACCESS_KEY = 'wc26.access';
const REFRESH_KEY = 'wc26.refresh';

export const tokens = {
  get access() { return localStorage.getItem(ACCESS_KEY); },
  get refresh() { return localStorage.getItem(REFRESH_KEY); },
  set(access: string | null, refresh?: string | null) {
    if (access === null) localStorage.removeItem(ACCESS_KEY);
    else localStorage.setItem(ACCESS_KEY, access);
    if (refresh !== undefined) {
      if (refresh === null) localStorage.removeItem(REFRESH_KEY);
      else localStorage.setItem(REFRESH_KEY, refresh);
    }
  },
  clear() { localStorage.removeItem(ACCESS_KEY); localStorage.removeItem(REFRESH_KEY); },
};

api.interceptors.request.use((config) => {
  const t = tokens.access;
  if (t) config.headers.Authorization = `Bearer ${t}`;
  return config;
});

let refreshing: Promise<string> | null = null;

api.interceptors.response.use(
  (r) => r,
  async (err) => {
    const original = err.config;
    if (err.response?.status === 401 && !original._retry && tokens.refresh) {
      original._retry = true;
      try {
        refreshing ??= axios
          .post(`${baseURL}/auth/refresh`, { refresh_token: tokens.refresh })
          .then((r) => {
            const newAccess = r.data.access_token as string;
            tokens.set(newAccess);
            return newAccess;
          })
          .finally(() => { refreshing = null; });
        const newAccess = await refreshing;
        original.headers.Authorization = `Bearer ${newAccess}`;
        return api(original);
      } catch {
        tokens.clear();
      }
    }
    return Promise.reject(err);
  },
);
