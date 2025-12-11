package repository

import (
	"contest/src/common"
	domain "contest/src/domain/entity"
	repository "contest/src/domain/repository/contest"
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

		Authors:     []uint64{author},
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
	err := cr.collection.FindOne(ctx, bson.M{
		"id": contestId,
	}).Decode(&contest)
	if err != nil {
		return domain.Contest{}, err
	}

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
		// For pulling, we might need to match by UserID if it's an object
		// But $pull with object should work if it matches exactly.
		// However, CreateContestant creates a new object with timestamp, so it might not match exactly if we just recreate it.
		// We should pull by UserID for contestants.
		// But existing code used data = CreateContestant(userId) which implies exact match or maybe structure match?
		// Actually, for $pull with array of objects, we can specify a query.
		// Let's check how it was done.
		// Old code: data = domain.CreateContestant(userId)
		// This suggests it was trying to remove exact object. But timestamp would differ.
		// This looks like a bug in original code or I misunderstand.
		// I will fix it to pull by UserID for contestants.
		data = bson.M{"user_id": userId}
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

	log.Info().Msgf("removed user %v from group %s of contest %s", data, peopleType, contestId)

	return nil
}
