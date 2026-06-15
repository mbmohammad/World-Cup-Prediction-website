package scoring

// Scorer awards points for a prediction given the actual result.
type Scorer struct {
	pointsExact   int
	pointsGD      int
	pointsOutcome int
}

func New(exact, gd, outcome int) *Scorer {
	return &Scorer{pointsExact: exact, pointsGD: gd, pointsOutcome: outcome}
}

// Points computes the score for a single prediction.
//   - Exact match (predH==actH && predA==actA) -> pointsExact
//   - Correct goal difference (and correct sign/outcome but not exact) -> pointsGD
//   - Correct outcome only (W/D/L matches) -> pointsOutcome
//   - Otherwise 0
func (s *Scorer) Points(predH, predA, actH, actA int) int {
	if predH == actH && predA == actA {
		return s.pointsExact
	}
	predDiff := predH - predA
	actDiff := actH - actA
	predOutcome := sign(predDiff)
	actOutcome := sign(actDiff)
	if predOutcome != actOutcome {
		return 0
	}
	if predDiff == actDiff {
		return s.pointsGD
	}
	return s.pointsOutcome
}

func sign(x int) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}
