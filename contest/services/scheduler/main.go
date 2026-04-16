package scheduler

import (
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	redisClient *redis.Client                 = nil
	contestRepo contestrepo.ContestRepository = nil
	ctx                                       = context.Background()
)

// Types of task

var listFuncs map[string]func(string) = map[string]func(string){
	"CONTEST_START":    StartContest,
	"CONTEST_END":      EndContest,
	"CONTEST_FREEZE":   func(contestId string) {},
	"CONTEST_UNFREEZE": func(contestId string) {},
	"CONTEST_FINALIZE": func(contestId string) {},
}
var (
	CONTEST_START  string = "CONTEST_START"
	CONTEST_END    string = "CONTEST_END"
	CONTEST_FREEZE string = "CONTEST_FREEZE"
)

func StartContest(contestId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	contest, err := contestRepo.GetById(ctx, contestId)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to find contest %s", contestId)
		return
	}

	log.Info().Msgf("Contest %s (id %s) has started", contest.Name, contestId)

	contestRepo.UpdateOne(ctx, contestId, map[string]any{
		"scoreboard_visibility": domain.ScoreboardPublic,
		"status":                domain.Running,
	})
}

func EndContest(contestId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	contest, err := contestRepo.GetById(ctx, contestId)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to find contest %s", contestId)
		return
	}

	log.Info().Msgf("Contest %s (id %s) has ended", contest.Name, contestId)

	contestRepo.UpdateOne(ctx, contestId, map[string]any{
		"status": domain.Ended,
	})
}

var TaskCallback map[string](func(Task)) = map[string](func(Task)){
	CONTEST_START: nil,
}

type Task struct {
	Type string
	Date time.Time
	Data map[string]string
}

var popTaskLua = redis.NewScript(`
	local val = redis.call('zrangebyscore', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, 1)
	if #val > 0 then
		redis.call('zrem', KEYS[1], val[1])
		return val[1]
	end
	return nil
`)

func Scheduler(re *redis.Client, repo contestrepo.ContestRepository) {
	redisClient = re
	contestRepo = repo

	for {
		cur := time.Now()
		fmt.Printf("%v\n", cur)
		now := cur.Unix()

		// GetAllScheduledTasks()

		result, err := popTaskLua.Run(ctx, redisClient, []string{"CONTEST_SCHEDULING"}, now).Result()
		if err != nil {
			if err.Error() == "redis: nil" {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			log.Err(err)
			time.Sleep(1 * time.Second)
			continue
		}

		var t Task
		if err := json.Unmarshal([]byte(result.(string)), &t); err == nil {
			log.Info().Msgf("Its time to do task %s", t.Type)
			listFuncs[t.Type](t.Data["contestId"])
		} else {
			log.Err(err)
		}
	}
}

func AddTask(task Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		log.Err(err).Msg("Failed to marshal task for addition")
		return err
	}

	redisClient.ZAdd(ctx, "CONTEST_SCHEDULING", redis.Z{
		Score:  float64(task.Date.Unix()),
		Member: data,
	})

	log.Info().Msgf("[CONTEST SCHEDULER] ADDED TASK type %s data %v", task.Type, task.Data)

	return nil
}

func RemoveTask(task Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		log.Err(err).Msg("Failed to marshal task for removal")
		return err
	}

	cmd := redisClient.ZRem(ctx, "CONTEST_SCHEDULING", data)
	if cmd.Err() != nil {
		log.Err(cmd.Err()).Msgf("Failed to remove task of type %s from Redis", task.Type)
		return cmd.Err()
	}
	if res, err := cmd.Result(); err != nil {
		log.Err(cmd.Err())
		return err
	} else if res == 0 {
		err = errors.New("task doesn't exist")
		log.Err(err)
		// return err
		// currently task doesn't exist is okay
	}

	log.Info().Msgf("[CONTEST SCHEDULER] REMOVED TASK type %s data %v", task.Type, task.Data)

	return nil
}

func GetAllScheduledTasks() ([]Task, error) {
	results, err := redisClient.ZRange(ctx, "CONTEST_SCHEDULING", 0, -1).Result()
	if err != nil {
		log.Err(err).Msg("Failed to get scheduled tasks from Redis")
		return nil, err
	}

	var tasks []Task
	fmt.Println("--- Scheduled Tasks ---")
	for _, result := range results {
		var t Task
		if err := json.Unmarshal([]byte(result), &t); err != nil {
			log.Err(err).Msg("Failed to unmarshal task")
			continue
		}
		tasks = append(tasks, t)
		fmt.Printf("Type: %-15s | Date: %s | Data: %v\n", t.Type, t.Date.Format(time.RFC3339), t.Data)
	}
	if len(tasks) == 0 {
		fmt.Println("No scheduled tasks found.")
	}
	fmt.Println("-----------------------")

	return tasks, nil
}

func ChangeContestStartDate(contestId string, newDate time.Time, repo contestrepo.ContestRepository) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	contest, err := repo.GetById(ctx, contestId)
	if err != nil {
		return err
	}
	oldDate := contest.StartTime

	RemoveTask(Task{
		Type: CONTEST_START,
		Date: oldDate,
		Data: map[string]string{
			"contestId": contestId,
		},
	})

	AddTask(Task{
		Type: CONTEST_START,
		Date: newDate,
		Data: map[string]string{
			"contestId": contestId,
		},
	})

	return nil
}

func ChangeContestEndDate(contestId string, newDate time.Time, repo contestrepo.ContestRepository) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	contest, err := repo.GetById(ctx, contestId)
	if err != nil {
		return err
	}
	oldDate := contest.EndTime

	RemoveTask(Task{
		Type: CONTEST_END,
		Date: oldDate,
		Data: map[string]string{
			"contestId": contestId,
		},
	})

	AddTask(Task{
		Type: CONTEST_END,
		Date: newDate,
		Data: map[string]string{
			"contestId": contestId,
		},
	})

	return nil
}
