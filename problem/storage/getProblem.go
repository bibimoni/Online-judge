package storage

import (
	"context"
	"os"
	"problem/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetAllProblems() ([]uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ProblemList, err := db.Database(os.Getenv("PROBLEM_DATABASE_NAME")).Collection("Problems").Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	ProblemIdList := make([]uint64, 0)
	for ProblemList.Next(ctx) {
		var result models.Problem
		if err := ProblemList.Decode(&result); err != nil {
			return nil, err
		}
		ProblemIdList = append(ProblemIdList, result.ProblemId)
	}

	return ProblemIdList, nil
}
