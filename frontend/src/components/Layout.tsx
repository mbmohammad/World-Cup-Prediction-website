import { Link, Outlet } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';

export default function Layout() {
  const { user, logout } = useAuth();
  return (
    <div className="min-h-screen flex flex-col">
      <header className="bg-white border-b">
        <nav className="max-w-5xl mx-auto px-4 py-3 flex items-center gap-6">
          <Link to="/" className="font-bold text-lg">WC2026 Predict</Link>
          <Link to="/matches" className="hover:underline">Matches</Link>
          <Link to="/leaderboard" className="hover:underline">Leaderboard</Link>
          <div className="ml-auto flex items-center gap-3">
            {user ? (
              <>
                <Link to="/profile" className="text-sm">{user.display_name}</Link>
                <button onClick={logout} className="text-sm text-slate-600 hover:text-slate-900">Logout</button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-sm">Login</Link>
                <Link to="/register" className="text-sm bg-emerald-600 text-white px-3 py-1.5 rounded">Sign up</Link>
              </>
            )}
          </div>
        </nav>
      </header>
      <main className="flex-1 max-w-5xl w-full mx-auto px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}
