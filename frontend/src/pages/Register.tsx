import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

export default function Register() {
  const { register } = useAuth();
  const nav = useNavigate();
  const [email, setEmail] = useState('');
  const [pwd, setPwd] = useState('');
  const [name, setName] = useState('');
  const [err, setErr] = useState<string | null>(null);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await register(email, pwd, name);
      nav('/matches');
    } catch (ex: any) {
      setErr(ex.response?.data?.error ?? 'Registration failed');
    }
  }

  return (
    <div className="max-w-sm mx-auto bg-white p-6 rounded shadow">
      <h1 className="text-xl font-semibold mb-4">Sign up</h1>
      <form onSubmit={onSubmit} className="space-y-3">
        <input className="w-full border px-3 py-2 rounded" placeholder="Display name"
               value={name} onChange={(e) => setName(e.target.value)} required minLength={2} />
        <input className="w-full border px-3 py-2 rounded" placeholder="Email" type="email"
               value={email} onChange={(e) => setEmail(e.target.value)} required />
        <input className="w-full border px-3 py-2 rounded" placeholder="Password (min 8 chars)" type="password"
               value={pwd} onChange={(e) => setPwd(e.target.value)} required minLength={8} />
        {err && <p className="text-sm text-red-600">{err}</p>}
        <button className="w-full bg-emerald-600 text-white py-2 rounded">Create account</button>
      </form>
    </div>
  );
}
