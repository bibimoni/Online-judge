package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"problem/models"
	"problem/utils/polygon"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
CheckAlreadyAdded(ProblemId)
*/
func CheckAlreadyAdded(ProblemId uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.Database(os.Getenv("PROBLEM_DATABASE_NAME")).Collection("Problems").FindOne(
		ctx,
		bson.M{"problem-id": int64(ProblemId)},
	).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

/*
AddProblem(ProblemId, PackageId)

Workflow:
- Download the packages using DownloadPackage()

FUTURE:
- Automatically get the latest package of the problem
*/
func AddProblem(ProblemId uint64) error {
	if found, err := CheckAlreadyAdded(ProblemId); err != nil {
		return fmt.Errorf("error while checking the problem repository: %s", err.Error())
	} else if found {
		return fmt.Errorf("problem already existed in the repository")
	}

	var PackageId uint64
	PackageId, err := polygon.GetLastestPackage(ProblemId)
	if err != nil {
		return fmt.Errorf("error getting latest package: %s", err.Error())
	}

	if err := polygon.DownloadPackage(ProblemId, PackageId); err != nil {
		return fmt.Errorf("error downloading package: %s", err.Error())
	}

	file, err := os.Open(fmt.Sprintf("%s/%d/problem.json", os.Getenv("PROBLEM_STORAGE_DIR"), ProblemId))
	if err != nil {
		return fmt.Errorf("error opening problem.xml: %s", err.Error())
	}
	defer file.Close()

	var problem models.Problem
	if err := json.NewDecoder(file).Decode(&problem); err != nil {
		return fmt.Errorf("error decoding problem.xml: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := db.Database("problem-db").Collection("Problems").InsertOne(ctx, problem); err != nil {
		return fmt.Errorf("error saving problem to database: %s", err.Error())
	}

	return nil
}
