import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

export default function Login() {
  const { login } = useAuth();
  const nav = useNavigate();
  const [email, setEmail] = useState('');
  const [pwd, setPwd] = useState('');
  const [err, setErr] = useState<string | null>(null);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await login(email, pwd);
      nav('/matches');
    } catch (ex: any) {
      setErr(ex.response?.data?.error ?? 'Login failed');
    }
  }

  return (
    <div className="max-w-sm mx-auto bg-white p-6 rounded shadow">
      <h1 className="text-xl font-semibold mb-4">Login</h1>
      <form onSubmit={onSubmit} className="space-y-3">
        <input className="w-full border px-3 py-2 rounded" placeholder="Email" type="email"
               value={email} onChange={(e) => setEmail(e.target.value)} required />
        <input className="w-full border px-3 py-2 rounded" placeholder="Password" type="password"
               value={pwd} onChange={(e) => setPwd(e.target.value)} required />
        {err && <p className="text-sm text-red-600">{err}</p>}
        <button className="w-full bg-emerald-600 text-white py-2 rounded">Sign in</button>
      </form>
      <p className="text-sm mt-3">No account? <Link to="/register" className="text-emerald-700">Register</Link></p>
    </div>
  );
}
