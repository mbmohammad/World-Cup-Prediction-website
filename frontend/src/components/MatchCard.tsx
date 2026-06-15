import { useState } from 'react';

type Team = { id: number; name: string; code: string; flag_url: string };
type Match = {
  id: number;
  home_team: Team; away_team: Team;
  kickoff_utc: string; stage: string;
  home_score?: number; away_score?: number;
  status: 'scheduled' | 'live' | 'finished';
};
type Prediction = { pred_home: number; pred_away: number; points_awarded?: number };

export default function MatchCard({
  match,
  prediction,
  onSubmit,
}: {
  match: Match;
  prediction?: Prediction;
  onSubmit: (h: number, a: number) => void;
}) {
  const [h, setH] = useState(prediction?.pred_home ?? 0);
  const [a, setA] = useState(prediction?.pred_away ?? 0);

  const kickoff = new Date(match.kickoff_utc);
  const locked = kickoff.getTime() <= Date.now() || match.status !== 'scheduled';

  return (
    <div className="bg-white rounded shadow p-4">
      <div className="text-xs text-slate-500 mb-2 flex justify-between">
        <span>{match.stage}</span>
        <span>{kickoff.toLocaleString()}</span>
      </div>
      <div className="flex items-center justify-between">
        <span className="font-medium">{match.home_team.code} {match.home_team.name}</span>
        <span className="text-slate-400">vs</span>
        <span className="font-medium">{match.away_team.name} {match.away_team.code}</span>
      </div>

      {match.status === 'finished' && (
        <div className="mt-2 text-center font-bold text-lg">
          {match.home_score} - {match.away_score}
        </div>
      )}

      <div className="mt-3 flex items-center gap-2 justify-center">
        <input type="number" min={0} max={20} disabled={locked}
               className="w-16 border px-2 py-1 rounded text-center"
               value={h} onChange={(e) => setH(Number(e.target.value))} />
        <span>-</span>
        <input type="number" min={0} max={20} disabled={locked}
               className="w-16 border px-2 py-1 rounded text-center"
               value={a} onChange={(e) => setA(Number(e.target.value))} />
        <button disabled={locked}
                className="ml-2 bg-emerald-600 text-white px-3 py-1.5 rounded disabled:bg-slate-300"
                onClick={() => onSubmit(h, a)}>
          {prediction ? 'Update' : 'Predict'}
        </button>
      </div>

      {prediction?.points_awarded != null && (
        <div className="mt-2 text-sm text-emerald-700 text-center">
          +{prediction.points_awarded} pts
        </div>
      )}
      {locked && !prediction && (
        <div className="mt-2 text-sm text-slate-500 text-center">Predictions closed</div>
      )}
    </div>
  );
}
