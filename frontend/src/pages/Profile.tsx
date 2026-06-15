import { useAuth } from '../auth/AuthContext';

export default function Profile() {
  const { user } = useAuth();
  if (!user) return null;
  return (
    <div className="bg-white rounded shadow p-6 max-w-md">
      <h1 className="text-xl font-semibold mb-4">Your profile</h1>
      <dl className="space-y-2 text-sm">
        <div><dt className="text-slate-500 inline">Display name: </dt><dd className="inline font-medium">{user.display_name}</dd></div>
        <div><dt className="text-slate-500 inline">Email: </dt><dd className="inline font-medium">{user.email}</dd></div>
      </dl>
    </div>
  );
}
