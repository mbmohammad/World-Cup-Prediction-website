import { useQuery } from '@tanstack/react-query';
import { api } from '../api/client';

type Entry = { rank: number; user_id: number; display_name: string; points: number; predictions: number };

export default function Leaderboard() {
  const q = useQuery({
    queryKey: ['leaderboard'],
    queryFn: async () => (await api.get<{ leaderboard: Entry[] }>('/leaderboard')).data.leaderboard,
  });

  if (q.isLoading) return <p>Loading…</p>;
  if (q.error) return <p className="text-red-600">Failed to load leaderboard.</p>;

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Leaderboard</h1>
      <div className="bg-white rounded shadow overflow-hidden">
        <table className="w-full">
          <thead className="bg-slate-100 text-left text-sm">
            <tr>
              <th className="px-4 py-2">Rank</th>
              <th className="px-4 py-2">Player</th>
              <th className="px-4 py-2 text-right">Predictions</th>
              <th className="px-4 py-2 text-right">Points</th>
            </tr>
          </thead>
          <tbody>
            {(q.data ?? []).map((e) => (
              <tr key={e.user_id} className="border-t">
                <td className="px-4 py-2 font-semibold">{e.rank}</td>
                <td className="px-4 py-2">{e.display_name}</td>
                <td className="px-4 py-2 text-right">{e.predictions}</td>
                <td className="px-4 py-2 text-right font-bold">{e.points}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
