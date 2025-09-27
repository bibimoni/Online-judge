package repository

import (
	"contest/common"
	domain "contest/domain/entity"
	repository "contest/domain/repository/contest"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContestRepositoryImpl struct {
	collection *mongo.Collection
}

func NewContestRepositoryImpl(db *mongo.Database) *ContestRepositoryImpl {
	return &ContestRepositoryImpl{
		collection: db.Collection("Contest"),
	}
}

func NewContestRepository(db *mongo.Database) repository.ContestRepository {
	return NewContestRepositoryImpl(db)
}

func (cr *ContestRepositoryImpl) Create(ctx context.Context, author uint64) (string, error) {
	newContest := domain.Contest{
		Id:          uuid.NewString(),
		Name:        "",
		Description: "",

		Authors:     []uint64{},
		Curators:    []uint64{},
		Testers:     []uint64{},
		Contestants: []domain.Contestant{},

		ProblemLabels: []string{},
		Problems:      []uint64{},

		ScoreboardVisibility: domain.SCOREBOARD_HIDDEN,

		StartTime: time.Now(),
		EndTime:   time.Now(),
	}

	_, err := cr.collection.InsertOne(ctx, newContest)
	if err != nil {
		return "", err
	}

	log.Info().Msgf("New contest created, id : %s", newContest.Id)

	return newContest.Id, nil
}

func (cr *ContestRepositoryImpl) GetById(contestId string) (domain.Contest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var contest domain.Contest
	cr.collection.FindOne(ctx, bson.M{
		"id": contestId,
	}).Decode(&contest)

	return contest, nil
}

func (cr *ContestRepositoryImpl) AddPeople(contestId string, peopleType string, userId uint64) error {
	if !slices.Contains(common.CONTEST_PEOPLE, peopleType) {
		return fmt.Errorf("invalid peopleType")
	}

	contest, err := cr.GetById(contestId)
	if err != nil {
		return err
	}

	fmt.Printf("contest: %v\n", contest)

	// Check if already exists
	if peopleType == common.CONTEST_CONTESTANTS && contest.ContestantExist(userId) {
		return fmt.Errorf("contestant %d already in contest %s", userId, contestId)
	}
	if peopleType == common.CONTEST_AUTHORS && slices.Contains(contest.Authors, userId) {
		return fmt.Errorf("author %d already in contest %s", userId, contestId)
	}
	if peopleType == common.CONTEST_CURATORS && slices.Contains(contest.Curators, userId) {
		return fmt.Errorf("curator %d already in contest %s", userId, contestId)
	}
	if peopleType == common.CONTEST_TESTERS && slices.Contains(contest.Testers, userId) {
		return fmt.Errorf("tester %d already in contest %s", userId, contestId)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var data interface{}
	if peopleType == common.CONTEST_CONTESTANTS {
		data = domain.CreateContestant(userId)
	} else {
		data = userId
	}

	_, err = cr.collection.UpdateOne(
		ctx,
		bson.M{"id": contestId},
		bson.M{"$push": bson.M{peopleType: data}},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("added user %v to group %s of contest %s", data, peopleType, contestId)

	return nil
}

func (cr *ContestRepositoryImpl) RemovePeople(contestId string, peopleType string, userId uint64) error {
	if !slices.Contains(common.CONTEST_PEOPLE, peopleType) {
		return fmt.Errorf("invalid peopleType")
	}

	contest, err := cr.GetById(contestId)
	if err != nil {
		return err
	}

	fmt.Printf("contest: %v\n", contest)

	// Check if already exists
	if peopleType == common.CONTEST_CONTESTANTS && !contest.ContestantExist(userId) {
		return fmt.Errorf("contestant %d is not in contest %s", userId, contestId)
	}
	if peopleType == common.CONTEST_AUTHORS && !slices.Contains(contest.Authors, userId) {
		return fmt.Errorf("author %d is not in contest %s", userId, contestId)
	}
	if peopleType == common.CONTEST_CURATORS && !slices.Contains(contest.Curators, userId) {
		return fmt.Errorf("curator %d is not in contest %s", userId, contestId)
	}
	if peopleType == common.CONTEST_TESTERS && !slices.Contains(contest.Testers, userId) {
		return fmt.Errorf("tester %d is not in contest %s", userId, contestId)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var data interface{}
	if peopleType == common.CONTEST_CONTESTANTS {
		data = domain.CreateContestant(userId)
	} else {
		data = userId
	}

	_, err = cr.collection.UpdateOne(
		ctx,
		bson.M{"id": contestId},
		bson.M{"$pull": bson.M{peopleType: data}},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("removed user %v to group %s of contest %s", data, peopleType, contestId)

	return nil
}

func (cr *ContestRepositoryImpl) EditProblemset(contestId string, newProblemset []uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := cr.collection.UpdateOne(
		ctx,
		bson.M{"id": contestId},
		bson.M{"$set": bson.M{"problems": newProblemset}},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("new problemset for contest %s: %v", contestId, newProblemset)

	return nil
}
