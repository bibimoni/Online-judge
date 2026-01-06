package contestrepoimpl

import (
	domain "contest/src/domain/entity"
	contestrepo "contest/src/domain/repository/contest"
	"context"
	"errors"
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

func NewContestRepository(db *mongo.Database) contestrepo.ContestRepository {
	return NewContestRepositoryImpl(db)
}

func (cr *ContestRepositoryImpl) CanCreateContest(ctx context.Context, role string) bool {
	// Only problem setter and admin can create contest
	if role == "admin" || role == "problem_setter" {
		return true
	}
	return false
}

func (cr *ContestRepositoryImpl) Create(ctx context.Context, creator string, contestname string) (string, error) {
	newContestRule := domain.ContestRule{
		ScoringType:                     domain.ICPC, // defaults to ICPC
		PenaltyMinutes:                  0,
		FreezeStartTime:                 time.Now(),
		FreezeTime:                      0,
		MaxAllowedSubmissionsPerProblem: 0,
	}
	newContest := domain.Contest{
		Name:        contestname,
		Description: "",

		Authors:     []string{},
		Admins:      []string{creator},
		Testers:     []string{},
		Contestants: []domain.Contestant{},

		Problems: []domain.ContestProblem{},

		ScoreboardVisibility: domain.ScoreboardHidden,

		StartTime: time.Now(),
		EndTime:   time.Now(),

		ContestRule: newContestRule,

		Status:           domain.Draft,
		FinalizeAt:       time.Now(),
		RejudgeWindowEnd: time.Now(),
	}

	result, err := cr.collection.InsertOne(ctx, newContest)
	if err != nil {
		return "", err
	}

	contestId := result.InsertedID.(bson.ObjectID).Hex()
	log.Info().Msgf("New contest created, id : %s", contestId)

	return contestId, nil
}

func (cr *ContestRepositoryImpl) GetById(ctx context.Context, contestId string) (*domain.Contest, error) {
	cId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return nil, err
	}

	var contest domain.Contest
	err = cr.collection.FindOne(ctx, bson.M{
		"_id": cId,
	}).Decode(&contest)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, contestrepo.ErrNoContestFound
		}
		return nil, err
	}

	return &contest, nil
}

func (cr *ContestRepositoryImpl) AddPeople(ctx context.Context, contestId string, peopleType contestrepo.PeopleType, username string) error {
	if !slices.Contains(contestrepo.ContestPeopple, peopleType) {
		return fmt.Errorf("invalid peopleType")
	}

	contest, err := cr.GetById(ctx, contestId)
	if err != nil {
		return err
	}

	cId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return err
	}

	// Check if already exists
	if peopleType == contestrepo.Contestant && contest.ContestantExist(username) {
		return fmt.Errorf("contestant %s already in contest %s", username, contestId)
	}
	if peopleType == contestrepo.Author && slices.Contains(contest.Authors, username) {
		return fmt.Errorf("author %s already in contest %s", username, contestId)
	}
	if peopleType == contestrepo.Admin && slices.Contains(contest.Admins, username) {
		return fmt.Errorf("curator %s already in contest %s", username, contestId)
	}
	if peopleType == contestrepo.Tester && slices.Contains(contest.Testers, username) {
		return fmt.Errorf("tester %s already in contest %s", username, contestId)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var data any
	if peopleType == contestrepo.Contestant {
		data = domain.CreateContestant(username)
	} else {
		data = username
	}

	_, err = cr.collection.UpdateOne(
		ctx,
		bson.M{"_id": cId},
		bson.M{"$push": bson.M{string(peopleType): data}},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("added user %v to group %s of contest %s", data, peopleType, contestId)

	return nil
}

func (cr *ContestRepositoryImpl) RemovePeople(ctx context.Context, contestId string, peopleType contestrepo.PeopleType, username string) error {
	if !slices.Contains(contestrepo.ContestPeopple, peopleType) {
		return fmt.Errorf("invalid peopleType")
	}

	contest, err := cr.GetById(ctx, contestId)
	if err != nil {
		return err
	}

	cId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return err
	}

	// Check if already exists
	if peopleType == contestrepo.Contestant && !contest.ContestantExist(username) {
		return fmt.Errorf("contestant %s is not in contest %s", username, contestId)
	}
	if peopleType == contestrepo.Author && !slices.Contains(contest.Authors, username) {
		return fmt.Errorf("author %s is not in contest %s", username, contestId)
	}
	if peopleType == contestrepo.Admin && !slices.Contains(contest.Admins, username) {
		return fmt.Errorf("curator %s is not in contest %s", username, contestId)
	}
	if peopleType == contestrepo.Tester && !slices.Contains(contest.Testers, username) {
		return fmt.Errorf("tester %s is not in contest %s", username, contestId)
	}

	var data any
	if peopleType == contestrepo.Contestant {
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
		data = bson.M{"username": username}
	} else {
		data = username
	}

	_, err = cr.collection.UpdateOne(
		ctx,
		bson.M{"_id": cId},
		bson.M{"$pull": bson.M{string(peopleType): data}},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("removed user %v from group %s of contest %s", data, peopleType, contestId)

	return nil
}

func (cr *ContestRepositoryImpl) ReplaceOne(ctx context.Context, contestId string, updatedContest *domain.Contest) error {
	cId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return err
	}

	_, err = cr.collection.ReplaceOne(ctx, bson.M{"_id": cId}, updatedContest)
	if err != nil {
		return err
	}

	log.Info().Msgf("replaced contest %s", contestId)
	return nil
}

func (cr *ContestRepositoryImpl) UpdateOne(ctx context.Context, contestId string, updateData map[string]any) error {
	cId, err := bson.ObjectIDFromHex(contestId)
	if err != nil {
		return err
	}

	_, err = cr.collection.UpdateOne(ctx, bson.M{"_id": cId}, bson.M{"$set": updateData})
	if err != nil {
		return err
	}

	log.Info().Msgf("updated contest %s with data %v", contestId, updateData)
	return nil
}
