package entity

import (
	"recommendation/internal/model"
	"time"
)

type (
	UserInteractionsBuilder struct {
		userInteractions UserInteractions
	}

	UserInteractions []UserInteraction
	UserInteraction  struct {
		UserID              string    `bson:"user_id"`
		ProductSKU          string    `bson:"product_sku"`
		Action              string    `bson:"action"`
		InteractionTime     time.Time `bson:"interaction_timestamp"`
		InteractionDuration *int      `bson:"interaction_duration,omitempty"` // Optional field
	}
)

func NewUserInteractionsBuilder() *UserInteractionsBuilder {
	return &UserInteractionsBuilder{
		userInteractions: make(UserInteractions, 0),
	}
}

func (ui *UserInteractionsBuilder) WithModel(u model.UserInteraction) *UserInteractionsBuilder {

	for _, v := range u.Interactions {
		interaction := UserInteraction{
			UserID:              string(u.UserID),
			ProductSKU:          string(v.Product),
			Action:              string(v.Action),
			InteractionTime:     v.InteractionTimestamp,
			InteractionDuration: v.InteractionDuration,
		}
		ui.userInteractions = append(ui.userInteractions, interaction)

	}

	return ui
}

func (ui *UserInteractionsBuilder) Build() *UserInteractions {
	return &ui.userInteractions
}

func (ui *UserInteractions) AsInterfaceSlice() *[]interface{} {
	interfaces := make([]interface{}, len(*ui))
	for i, v := range *ui {
		interfaces[i] = v
	}
	return &interfaces
}
