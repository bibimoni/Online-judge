# Contest service

## Current choice and future consideration
- Run all tests in contest (no pretests)
- No rating system
- No virtual participation
- Authors can edit contest's settings
- Scoreboard only works for running contests, ended contests just need a json of the scoreboard

## API Schema

```go
type Contest struct {
	Id          string `bson:"id"`
	Name        string `bson:"name"`
	Description string `bson:"description"`

	Authors     []string                                     `bson:"authors"`
	Curators    []string                                     `bson:"curators"`
	Testers     []string                                     `bson:"testers"`
	Contestants *sortedset.SortedSet[uint64, uint64, uint64] `bson:"contestants"`

	ProblemLabels []string `bson:"problem_labels"`
	Problems      []uint64 `bson:"problems"`

	ScoreboardVisibility string `bson:"scoreboard_visibility"`

	StartTime time.Time `bson:"start_time"`
	EndTime   time.Time `bson:"end_time"`
}

type Contestant struct {
	UserID    uint64    `bson:"user_id"`
	Points    uint64    `bson:"points"`
	RealStart time.Time `bson:"real_start"`

	Submissions []uint64 `bson:"submissions"`
}

```

### Scoreboard Visibility
- SCOREBOARD_HIDDEN: Scoreboard is not visible by anyone
- SCOREBOARD_PUBLIC: Scoreboard is public to everyone
- SCOREBOARD_CONTESTANT_ONLY: Only contestant could see this scoreboard
