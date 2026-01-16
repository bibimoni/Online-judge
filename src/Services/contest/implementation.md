# Implementation Plan

Below is the implementation plan for adding contest functionality to the online judge system by following the SOLID principles

## Rating systems

We are going to use AtCoder-style rating system for our competitive programming platform.

## Domain model

We already have `Contest` and `Contestant`. Extend as follows:
### Contest
Add: 
- `ScoringType` (enum: `ICPC`, `IOI`, `CUSTOM`)
- `PenaltyRules` (for ICPC scoring config: penalty minutes, freeze config, etc.)
- `IOI Rules` (sum of subtasks, max-sumbission-per-problem, etc.)
- `RatingPolicy` (enum: ATCODER, CODEFORCES)
- `Rated` (bool)
- `FreezeTime` (optional; even if not know, design for it)
- `Status` (enum: DRAFT, SCHEDULED, ONGOING, ENDED, RATED)
- `FinalizedAt` (time.Time)
- `RejudgeWindowEnd` (time.Time, the time to rejudge, after this the scoreboard will be finalized)

### ContestProblem (new)
- `ProblemID` (reference to Problem)
- `Label` (or `shortname`) (string) 
- `MaxPoints` (float64; for IOI scoring, ICPC can default to 1)

### ContestSubmission (new)

To support rejudge/skip reliably we store the contest view of the submission.

- `SubmissionID` (reference to Submission)
- `ContestID` (reference to Contest)
- `Username` (string)
- `ProblemID` (reference to ContestProblem)
- `SubmittedAt` (time.Time)
- `Verdict` (enum: ACCEPTED, WRONG_ANSWER, TIME_LIMIT_EXCEEDED, etc.)
- `Points` (float64; for IOI scoring)
- `JudgedAt` (time.Time)
- `EvalStatus` (enum: PENDING, JUDGING, FINISHED)
- `Ignored` (bool)
- `Updated At` (time.Time)

### ScoreboardSnapshot (new)
- `ContestID` (reference to Contest)
- `Kind` (enum: LIVE, FINAL)
- `Payload` (rows + per-problem results)
- `CreatedAt` (time.Time)

### Rating Result (new)
- `ContestID` (reference to Contest)
- `Username` (string)
- `OldRating` (int)
- `NewRating` (int)
- `Delta` (int)
- `Rank`
- `Performance` (optional)

- `SubmissionIDs` (string[])
- `ProblemID` (reference to Problem)
- `RequestedBy` (string)
<!-- - `Status` (enum: PENDING, JUDGING, COMPLETED, FAILED) -->

## Repository interfaces

### Contest Repository
- `Create(ctx, contest Contest) (id, error)`
- `GetByID(ctx, id) (Contest, error)`
- `UpdateMeta(ctx, id, patch) error`
- `UpdateProblems(ctx, id, problems []ContestProblem) error`
- `RegisterContestant(ctx, contestID, username) error`
- `SetRated(ctx, contestID, rated bool, policy string) error`
- `SetScoringType(ctx, contestID, scoringType string, rules any) error`

### ContestSubmission Repository
- `UpsertFromJudgeEvent(ctx, contestID, submission ...) error`
- `ListByContest(ctx, contestID, filters...) ([]ContestSubmission, error)`
- `SetIgnored(ctx, contestID, submissionID, ignored bool) error`
<!-- - `BulkSetIgnoredByContestant(ctx, contestID, username, ignored bool) error` -->
- `SetIgnoredByUser(ctx, contestID, username, ignored bool) error`
- `UpdateVerdictPointsForRejudge(ctx, contestID, submissionID, verdict, points) error`

### Scoreboard Repository
- `SaveSnapshot(ctx, contestID, kind string, payload ScoreboardSnapshot) error`
- `GetSnapshot(ctx, contestID, kind string) (ScoreboardSnapshot, error)`

### ContestRatingResult Repository
- `GetRating(ctx, contestID, username) (rating, error)`
- `BatchGetRatings(ctx, contestID, username []string) (...)`
- `SaveRatingChanges(ctx, results []RatingResult) error`

## Services
### Scoring Service
- Name: `ScoreboardCalculator`
- Interface: 
  - `BuildScoreboard(contest Contest, subs []ContestSubmission) (ScoreboardSnapshot, error)`
- Implementations:
    - `ICPCScoreboardCalculator`
    - `IOIScoreboardCalculator`
    ...

### Rating Service
- Name: `RatingPolicy`
- Interface:
  - `Compute(contest Contest, standings ScoreboardSnapshot, oldRatings map[user]rating) ([]RatingResult, error)`
- Implementations:
    - `AtCoderRatingPolicy` (the naming here is just for reading comfort; can be just `RatingPolicy`)
    ...

### Contest Service
- Name: `Contest`
<...> implement create/edit/register/permissions

### ContestEvent Service
- Name: `ContestEvent`
<...> listen to judge events, update contest submissions, trigger scoreboard recalculation and rating calculation

### Rejudge Service
- Name: `ContestRejudge`
<...> handle rejudge requests, update contest submissions, trigger scoreboard recalculation.
### Moderation Service
- Name: `ContestModeration`
<...> handle ignore/unignore requests for contest submissions (and more...)

### Ingestion Service
- Name: `ContestIngestionService`
- Role: 
  - Validate contest exists
  - persist/upsert `ContestSubmission`
  - recompute scoreboard

## Usecases
Below is the current list of use cases i can think of right now. Each usecase should calls services (not repositories directly) to perform its job. (note: If this rule makes things harder to implement, we can ignore it, but please try to follow it as much as possible).

- `create_contest`
- `get_contest`
- `list_contests`
- `update_contest_meta`
- `set_contest_problems`
- `register_contest`
- `get_scoreboard` (call calculator)
- `finalize_scoreboard`
- `ingest_submission (internal)` (from judge events, who said usecases can't be used internally? lmao)
- `skip_submission`
- `unskip_submission`
- `request_rejudge`
- `rejudge_contest`
- `compute_ratings`
- `finalize_ratings`

## Endpoint list (changed a bit)
### Contest 
- `POST /api/v1/contest` (create)
- `GET /api/v1/contest/:id` 
- `GET /api/v1/contest` (list; status=upcoming|running|ended, paging)
- `PATCH /api/v1/contest/:id` (meta: name, desc, start/end, visibility)
- `PUT /api/v1/contest/:id/problems`
- `PUT /api/v1/contest/:id/scoring` (ICPC/IOI + rules)
- `PUT /api/v1/contest/:id/rating` (enable rated + policy)

### Contestant
- `POST /api/v1/contest/:id/register`
- `POST /api/v1/contest/:id/unregister`
- `POST /api/v1/contest/:id/submit`

### Scoreboard
- `GET /api/v1/contest/:id/scoreboard` (live)
- `POST /api/v1/contest/:id/scoreboard/finalize` (admin/author; produces snapshot)
- `GET /api/v1/contest/:id/scoreboard/final`
- `WS /api/v1/contest/:id/scoreboard/ws` (live updates)

### Moderation
- `POST /api/v1/contest/:id/submission/:sid/skip`
- `POST /api/v1/contest/:id/submission/:sid/unskip` (Optional)
- `POST /api/v1/contest/:id/user/:username/skip` (bulk ignore)

### Rejudge
Requests judge to rejudge
- `POST /api/v1/contest/:id/rejudge/submission/:sid`
- `POST /api/v1/contest/:id/rejudge/problem/:pid`
- `POST /api/v1/contest/:id/rejudge` (whole contest)

### Internal integration
- `POST /api/v1/internal/contest/submission-events` (ingest submission events from judge)

## Implementation steps
1. Fix schema
2. Contest
3. Submission ingestion
4. Scoreboard calculation
5. Skip solution
6. Rejudge flows
7. Finalize scoreboard snapshot
8. Rating finalize

## Important notes
### What is submission ingestion?
1. The client submit in the contest via `/submit` endpoint of the `contest` service. Then the `contest` service forwards the submission to the `submission-judge` service. Receives the `submission_id`. Create `ContestSubmission` object. It will stores mapping in its DB. 

2. Then when the judging is done, the `submission-judge` service calls the `ingest_submission` usecase in the `contest` service. 

3. Contest service then lookup `ContestSubmission` by `submisison_id`, if exists update `ContestSubmission`/scoreboard, otherwsie ignore (not a contest submission).
3.1 It also check if the submission belongs to a rejudge job, and update accordingly.

### What is scoreboard recalculation?
When a new submission is ingested, or a submission is skipped/unskipped, or a rejudge is done, the scoreboard needs to be recalculated.

Note that this is different than rating calculation. Scoreboard recalculation means building the current scoreboard view (live or final) based on the current `ContestSubmission` entries. 

### Live Scoreboard
Mechanism to update live scoreboard: 
1. Contest service computes/update live scoreboard snapshosts (in Redis) every 20 seconds (configurable).
2. After each snapshot update, contest service publishes an `update_event` on Redis PubSub.
- Channel: `contest:scoreboard:{contestID}`
- Payload: `{contest_id, version, generated_at}` then the client fetches via HTTP. (Why? Full snapshot can be big, delta update is complex for client to handle)

### Rating 
The actual rating will be stored in the user(or `auth`) service. Not in the `contest` service. The `contest` service will store the rating changes per contest.

### Rejudge 
Rejudge is an async operation that re-runs judging for an existing set of submissions and updates:
- `ContestSubmission`
- Live scoreboard (recompute)
- Final scoreboard + rating

Because it's async, the rejudge request will create a `RejudgeJob` and return a rejudge_job_id

## Implementation note
### Auth route
- All user information will be stored in the `request context` by the auth middlware 
### `ingest_submission` usecase
1. validate contest exists
2. store/Upsert `ContestSubmission` 
3. Enqueue scoreboard recalculation job (async) 
4. return `ok`

### Background worker: `ScoreboardRecalculationWorker`
1. Consumes `recompute` jobs (Redis)
2. Recalculate scoreboard
3. store snapshot in Redis (and Mongo for final)
4. Publish `update_event` on Redis PubSub

### In `submission-judge` service
After judging a submission (status = `FINISHED`), call `ingest_submission` usecase in `contest` service via HTTP API. (or maybe even `JUDGING`)

Also implement a rejudge API call.
`POST /api/v1/internal/rejudge/submissions` with batch of submission IDs to rejudge.

### Rejudge
1. User calls rejudge endpoint
2. Create `RejudgeJob`
3. Call `submission-judge` rejudge API with submission IDs
4. `submission-judge` rejudges (by adding the submission to its queue). Update eval, submission status in its database.
5. After each submission is rejudged, `submission-judge` calls `ingest_submission` usecase in `contest` service
6. `ingest_submission` updates `ContestSubmission`, or in this case also `RejudgeJob` status if all submissions are done.
