package model

import (
	"fmt"
	"time"
)

var validActions = []Action{"view", "click", "add_to_cart"}

type (
	ProductSKU string
	UserID     string
	Action     string

	// UserInteraction represents a user's interactions with products.
	UserInteraction struct {
		UserID       UserID        `json:"user_id"`      // Required
		Interactions []Interaction `json:"interactions"` // Required
	}

	// Interaction represents a single user interaction with a product.
	Interaction struct {
		Product              ProductSKU `json:"product_sku"`                    // Required
		Action               Action     `json:"action"`                         // Required, one of "view", "click", "add_to_cart"
		InteractionTimestamp time.Time  `json:"interaction_timestamp"`          // Required
		InteractionDuration  *int       `json:"interaction_duration,omitempty"` // Optional
	}

	// ProductsRecommendation represents product recommendations for a user.
	ProductsRecommendation struct {
		Products []ProductSKU `json:"products"` // List of product SKUs recommended to the user
	}
)

func (u UserID) Validate() APIError {
	if u == "" {
		return NewBadRequestApiError("Invalid User ID")
	}

	return nil
}

func (a Action) Validate() bool {
	for _, validAction := range validActions {
		if a == validAction {
			return true
		}
	}

	return false
}

func (ui UserInteraction) Validate() APIError {
	if err := ui.UserID.Validate(); err != nil {
		return err
	}

	if len(ui.Interactions) == 0 {
		return NewBadRequestApiError("Missing user interactions")
	}

	for _, interaction := range ui.Interactions {
		if err := interaction.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (i Interaction) Validate() APIError {
	if i.Product == "" {
		return NewBadRequestApiError("Invalid Product SKU")
	}

	if !i.Action.Validate() {
		return NewBadRequestApiError(fmt.Sprintf("Invalid Action: %s", i.Action))
	}

	return nil
}
