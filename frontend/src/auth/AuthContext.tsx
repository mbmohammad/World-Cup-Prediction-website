import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { api, tokens } from '../api/client';

type User = { id: number; email: string; display_name: string };

type AuthCtx = {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, displayName: string) => Promise<void>;
  logout: () => void;
};

const Ctx = createContext<AuthCtx>(null as unknown as AuthCtx);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!tokens.access) { setLoading(false); return; }
    api.get<User>('/me')
      .then((r) => setUser(r.data))
      .catch(() => tokens.clear())
      .finally(() => setLoading(false));
  }, []);

  async function login(email: string, password: string) {
    const r = await api.post('/auth/login', { email, password });
    tokens.set(r.data.tokens.access_token, r.data.tokens.refresh_token);
    setUser(r.data.user);
  }

  async function register(email: string, password: string, display_name: string) {
    const r = await api.post('/auth/register', { email, password, display_name });
    tokens.set(r.data.tokens.access_token, r.data.tokens.refresh_token);
    setUser(r.data.user);
  }

  function logout() {
    tokens.clear();
    setUser(null);
  }

  return <Ctx.Provider value={{ user, loading, login, register, logout }}>{children}</Ctx.Provider>;
}

export const useAuth = () => useContext(Ctx);
