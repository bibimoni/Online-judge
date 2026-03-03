package ratingserviceimpl

import (
	domain "contest/src/domain/entity"
	ratingservice "contest/src/service/rating"
	"context"
	"math"
	"sort"
	"time"
)

type EloMMRService struct {
	cfg domain.EloMMRConfig
}

func NewEloMMRService(cfg domain.EloMMRConfig) *EloMMRService {
	return &EloMMRService{cfg: cfg}
}

func NewDefaultEloMMRService() ratingservice.RatingService {
	return NewEloMMRService(domain.DefaultEloMMRConfig())
}

func (s *EloMMRService) ComputeRatings(
	ctx context.Context,
	contest *domain.Contest,
	standings []domain.StandingRow,
	oldStates map[string]domain.UserRatingState,
) ([]domain.RatingResult, map[string]domain.UserRatingState, error) {
	if len(standings) == 0 {
		return nil, nil, ratingservice.ErrNoStandings
	}

	n := len(standings)

	mus := make([]float64, n)
	sigmas := make([]float64, n)
	for i, row := range standings {
		st, ok := oldStates[row.Username]
		if !ok {
			st = domain.NewUserRatingState(row.Username, s.cfg)
		}
		st = s.applyDrift(st, contest.EndTime)
		mus[i] = st.Mu
		sigmas[i] = st.Sigma
	}

	perfs := s.inferPerformances(standings, mus, sigmas)

	newStates := make(map[string]domain.UserRatingState, n)
	results := make([]domain.RatingResult, n)
	for i, row := range standings {
		oldSt, ok := oldStates[row.Username]
		if !ok {
			oldSt = domain.NewUserRatingState(row.Username, s.cfg)
		}
		oldDisplay := oldSt.DisplayRating

		newMu, newSig := s.updateBelief(mus[i], sigmas[i], perfs[i])
		newDisplay := newMu - 2*newSig

		ns := domain.UserRatingState{
			Username:      row.Username,
			Mu:            newMu,
			Sigma:         newSig,
			DisplayRating: newDisplay,
			NumContests:   oldSt.NumContests + 1,
			LastContestAt: contest.EndTime,
		}
		newStates[row.Username] = ns

		results[i] = domain.RatingResult{
			ContestId:   contest.Id,
			Username:    row.Username,
			Rank:        row.Rank,
			OldRating:   math.Round(oldDisplay),
			NewRating:   math.Round(newDisplay),
			Delta:       math.Round(newDisplay - oldDisplay),
			Performance: math.Round(perfs[i]),
			CreatedAt:   time.Now(),
		}
	}

	return results, newStates, nil
}

func (s *EloMMRService) applyDrift(st domain.UserRatingState, contestTime time.Time) domain.UserRatingState {
	if s.cfg.DriftPerSec <= 0 || st.LastContestAt.IsZero() {
		return st
	}
	elapsed := contestTime.Sub(st.LastContestAt).Seconds()
	if elapsed <= 0 {
		return st
	}
	st.Sigma = math.Sqrt(st.Sigma*st.Sigma + s.cfg.DriftPerSec*elapsed)
	return st
}

func (s *EloMMRService) inferPerformances(
	standings []domain.StandingRow,
	mus, sigmas []float64,
) []float64 {
	n := len(standings)
	beta := s.cfg.Beta
	perfs := make([]float64, n)
	copy(perfs, mus)

	expectedWins := make([]float64, n)
	for i, row := range standings {
		expectedWins[i] = float64(n-row.Rank) + 0.5*(countTies(standings, row.Rank)-1)
	}


	for iter := 0; iter < 50; iter++ {
		maxDelta := 0.0
		for i := 0; i < n; i++ {

			var winsSum, gradSum float64
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				p := logistic(perfs[i], perfs[j], beta)
				winsSum += p
				gradSum += p * (1 - p) / beta
			}

			if gradSum < 1e-12 {
				continue
			}

			priorWeight := 1.0 / (sigmas[i]*sigmas[i] + beta*beta)
			priorPull := priorWeight * (mus[i] - perfs[i])

			delta := (expectedWins[i] - winsSum + priorPull) / (gradSum + priorWeight)
			perfs[i] += delta
			if math.Abs(delta) > maxDelta {
				maxDelta = math.Abs(delta)
			}
		}
		if maxDelta < 0.01 {
			break
		}
	}

	return perfs
}

func (s *EloMMRService) updateBelief(mu, sigma, perf float64) (newMu, newSig float64) {
	beta := s.cfg.Beta
	weight := math.Min(s.cfg.WeightLimit, 1.0) 
	priorPrec := 1.0 / (sigma * sigma)
	likelihoodPrec := weight / (beta * beta)

	postPrec := priorPrec + likelihoodPrec
	newMu = (mu*priorPrec + perf*likelihoodPrec) / postPrec
	newSig = math.Max(math.Sqrt(1.0/postPrec), s.cfg.SigLimit)

	return newMu, newSig
}

func logistic(ra, rb, beta float64) float64 {
	return 1.0 / (1.0 + math.Pow(10, (rb-ra)/(2*beta)))
}

func countTies(standings []domain.StandingRow, rank int) float64 {
	cnt := 0.0
	for _, s := range standings {
		if s.Rank == rank {
			cnt++
		}
	}
	return cnt
}

func StandingsFromScoreboard(
	snapshot *domain.ScoreboardSnapshot,
	contest *domain.Contest,
) []domain.StandingRow {
	ratedSet := make(map[string]bool)
	for _, c := range contest.Contestants {
		if c.ParticipantType == domain.RatedParticipant {
			ratedSet[c.Username] = true
		}
	}

	type entry struct {
		username string
		score    float64
		penalty  int
	}
	var entries []entry
	for _, row := range snapshot.Rows {
		if !ratedSet[row.Username] {
			continue
		}
		entries = append(entries, entry{
			username: row.Username,
			score:    row.Score,
			penalty:  row.Penalty,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].score != entries[j].score {
			return entries[i].score > entries[j].score
		}
		return entries[i].penalty < entries[j].penalty
	})
	standings := make([]domain.StandingRow, len(entries))
	for i := range entries {
		rank := i + 1
		if i > 0 && entries[i].score == entries[i-1].score && entries[i].penalty == entries[i-1].penalty {
			rank = standings[i-1].Rank
		}
		standings[i] = domain.StandingRow{
			Username: entries[i].username,
			Rank:     rank,
			Score:    entries[i].score,
		}
	}

	return standings
}
