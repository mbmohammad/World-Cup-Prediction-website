import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../api/client';
import MatchCard from '../components/MatchCard';

type Team = { id: number; name: string; code: string; flag_url: string };
type Match = {
  id: number;
  home_team: Team; away_team: Team;
  kickoff_utc: string; stage: string; group?: string;
  home_score?: number; away_score?: number;
  status: 'scheduled' | 'live' | 'finished';
};
type Prediction = { match_id: number; pred_home: number; pred_away: number; points_awarded?: number };

export default function Matches() {
  const qc = useQueryClient();

  const matchesQ = useQuery({
    queryKey: ['matches'],
    queryFn: async () => (await api.get<{ matches: Match[] }>('/matches')).data.matches,
  });
  const predsQ = useQuery({
    queryKey: ['my-predictions'],
    queryFn: async () => (await api.get<{ predictions: Prediction[] }>('/predictions')).data.predictions,
  });

  const submit = useMutation({
    mutationFn: async (p: { match_id: number; pred_home: number; pred_away: number }) =>
      (await api.put(`/predictions/${p.match_id}`, { pred_home: p.pred_home, pred_away: p.pred_away })).data,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['my-predictions'] }),
  });

  if (matchesQ.isLoading) return <p>Loading matches…</p>;
  if (matchesQ.error) return <p className="text-red-600">Failed to load matches.</p>;

  const predByMatch = new Map((predsQ.data ?? []).map((p) => [p.match_id, p]));

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Fixtures</h1>
      <div className="grid gap-3 md:grid-cols-2">
        {(matchesQ.data ?? []).map((m) => (
          <MatchCard
            key={m.id}
            match={m}
            prediction={predByMatch.get(m.id)}
            onSubmit={(h, a) => submit.mutate({ match_id: m.id, pred_home: h, pred_away: a })}
          />
        ))}
      </div>
    </div>
  );
}
