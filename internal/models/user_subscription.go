package models

import (
	"github.com/google/uuid"
	"time"

	"github.com/major75/online-subscriptions/pkg/types"
)

type CreateUserSubscriptionRequest struct {
	UserID      uuid.UUID         `json:"user_id" validate:"required" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName string            `json:"service_name" validate:"min=1,max=255"`
	Price       int64             `json:"price" validate:"gte=1" example:"1"`
	StartDate   *types.MMYYYYDate `json:"start_date" validate:"required" swaggertype:"string" example:"07-2026"`
	StopDate    *types.MMYYYYDate `json:"stop_date,omitempty" validate:"omitempty" swaggertype:"string" example:"07-2026"`
}

type CreateUserSubscriptionResponse struct {
	ID        uint32    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserSubscription struct {
	ID          uint32     `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	ServiceName string     `json:"service_name"`
	Price       int64      `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	StopDate    *time.Time `json:"stop_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PatchUserSubscriptionRequest struct {
	Price     int64             `json:"price" validate:"gte=1" example:"1"`
	StartDate *types.MMYYYYDate `json:"start_date" validate:"required" swaggertype:"string" example:"07-2026"`
	StopDate  *types.MMYYYYDate `json:"stop_date,omitempty" validate:"omitempty" swaggertype:"string" example:"07-2026"`
}

type SubscriptionsTotalFilter struct {
	DateFrom    time.Time
	DateTo      time.Time
	UserID      *uuid.UUID
	ServiceName *string
}

type SubscriptionsTotal struct {
	Total uint64 `json:"total"`
}
