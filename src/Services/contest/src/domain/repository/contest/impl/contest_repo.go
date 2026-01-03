package contestrepoimpl

import (
	"contest/src/common"
	domain "contest/src/domain/entity"
	repository "contest/src/domain/repository/contest"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (cr *ContestRepositoryImpl) Create(ctx context.Context, author string) (string, error) {
	newContest := domain.Contest{
		// Id:          uuid.NewString(),
		Name:        "",
		Description: "",

		Authors:     []string{author},
		Curators:    []string{},
		Testers:     []string{},
		Contestants: []domain.Contestant{},

		ProblemLabels: []string{},
		Problems:      []uint64{},

		ScoreboardVisibility: domain.ScoreboardHidden,

		StartTime: time.Now(),
		EndTime:   time.Now(),
	}

	result, err := cr.collection.InsertOne(ctx, newContest)
	if err != nil {
		return "", err
	}

	contestId := result.InsertedID.(bson.ObjectID).Hex()
	log.Info().Msgf("New contest created, id : %s", contestId)

	return contestId, nil
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

func (cr *ContestRepositoryImpl) AddPeople(contestId string, peopleType string, username string) error {
	if !slices.Contains(common.CONTEST_PEOPLE, peopleType) {
		return fmt.Errorf("invalid peopleType")
	}

	contest, err := cr.GetById(contestId)
	if err != nil {
		return err
	}

	// Check if already exists
	if peopleType == common.CONTEST_CONTESTANTS && contest.ContestantExist(username) {
		return fmt.Errorf("contestant %d already in contest %s", username, contestId)
	}
	if peopleType == common.CONTEST_AUTHORS && slices.Contains(contest.Authors, username) {
		return fmt.Errorf("author %d already in contest %s", username, contestId)
	}
	if peopleType == common.CONTEST_CURATORS && slices.Contains(contest.Curators, username) {
		return fmt.Errorf("curator %d already in contest %s", username, contestId)
	}
	if peopleType == common.CONTEST_TESTERS && slices.Contains(contest.Testers, username) {
		return fmt.Errorf("tester %d already in contest %s", username, contestId)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var data interface{}
	if peopleType == common.CONTEST_CONTESTANTS {
		data = domain.CreateContestant(username)
	} else {
		data = username
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

func (cr *ContestRepositoryImpl) RemovePeople(contestId string, peopleType string, username string) error {
	if !slices.Contains(common.CONTEST_PEOPLE, peopleType) {
		return fmt.Errorf("invalid peopleType")
	}

	contest, err := cr.GetById(contestId)
	if err != nil {
		return err
	}

	// Check if already exists
	if peopleType == common.CONTEST_CONTESTANTS && !contest.ContestantExist(username) {
		return fmt.Errorf("contestant %d is not in contest %s", username, contestId)
	}
	if peopleType == common.CONTEST_AUTHORS && !slices.Contains(contest.Authors, username) {
		return fmt.Errorf("author %d is not in contest %s", username, contestId)
	}
	if peopleType == common.CONTEST_CURATORS && !slices.Contains(contest.Curators, username) {
		return fmt.Errorf("curator %d is not in contest %s", username, contestId)
	}
	if peopleType == common.CONTEST_TESTERS && !slices.Contains(contest.Testers, username) {
		return fmt.Errorf("tester %d is not in contest %s", username, contestId)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var data interface{}
	if peopleType == common.CONTEST_CONTESTANTS {
		// For pulling, we might need to match by UserID if it's an object
		// But $pull with object should work if it matches exactly.
		// However, CreateContestant creates a new object with timestamp, so it might not match exactly if we just recreate it.
		// We should pull by UserID for contestants.
		// But existing code used data = CreateContestant(username) which implies exact match or maybe structure match?
		// Actually, for $pull with array of objects, we can specify a query.
		// Let's check how it was done.
		// Old code: data = domain.CreateContestant(username)
		// This suggests it was trying to remove exact object. But timestamp would differ.
		// This looks like a bug in original code or I misunderstand.
		// I will fix it to pull by UserID for contestants.
		data = bson.M{"user_id": username}
	} else {
		data = username
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
