package repository

import (
	"context"
	"fmt"
	"log"
	"recommendation/internal/entity"
	"recommendation/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserBehaviorRepository interface {
	GetUserInteractionsSince(userID model.UserID, period time.Duration) (*model.UserInteraction, error)
	AddUserInteractions(userInteraction model.UserInteraction) error
}

type userBehaviorRepository struct {
	collection *mongo.Collection
}

func NewUserBehaviorRepository(c *mongo.Collection) UserBehaviorRepository {

	return &userBehaviorRepository{
		collection: c,
	}
}

func (r *userBehaviorRepository) GetUserInteractionsSince(userID model.UserID, period time.Duration) (*model.UserInteraction, error) {

	var (
		query   = buildQuery(userID, period)
		results []entity.UserInteraction
	)

	cursor, err := r.collection.Find(context.TODO(), query)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into the results slice
	//TODO use request context
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Print(err)
		return nil, err
	}

	return buildResponse(results), nil

}

func buildResponse(userInteractions []entity.UserInteraction) *model.UserInteraction {
	if len(userInteractions) == 0 {
		return &model.UserInteraction{}
	}

	var response = model.UserInteraction{
		UserID:       model.UserID(userInteractions[0].UserID),
		Interactions: make([]model.Interaction, 0, len(userInteractions)),
	}

	for _, interaction := range userInteractions {
		i := model.Interaction{
			Product:              model.ProductSKU(interaction.ProductSKU),
			Action:               model.Action(interaction.Action),
			InteractionTimestamp: interaction.InteractionTime,
			InteractionDuration:  interaction.InteractionDuration,
		}
		response.Interactions = append(response.Interactions, i)
	}
	return &response
}

func (r *userBehaviorRepository) AddUserInteractions(userInteraction model.UserInteraction) error {

	// Create a new UserInteraction instance
	interactions := entity.NewUserInteractionsBuilder().WithModel(userInteraction).Build()

	insertResult, err := r.collection.InsertMany(context.TODO(), *interactions.AsInterfaceSlice())
	if err != nil {
		log.Print(err)
		return err
	}

	// Print the ID of the inserted document
	fmt.Printf("Inserted document with ID: %v\n", insertResult.InsertedIDs...)
	return nil
}

// TODO refactor this to a more berbose builder
func buildQuery(id model.UserID, period time.Duration) bson.D {
	sincePeriodAgo := time.Now().Add(-period)
	query := bson.D{
		{Key: "user_id", Value: id},
		{Key: "interaction_timestamp", Value: bson.D{
			{Key: "$gte", Value: sincePeriodAgo},
		}},
	}
	return query
}
