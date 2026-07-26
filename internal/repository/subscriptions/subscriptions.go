package subscriptions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/major75/online-subscriptions/database"
	"github.com/major75/online-subscriptions/internal/models"
	"github.com/major75/online-subscriptions/pkg/logger"
)

type UserSubscription interface {
	Create(ctx context.Context, subscription models.CreateUserSubscriptionRequest) (models.CreateUserSubscriptionResponse, error)
	Update(ctx context.Context, id uint32, subscription models.CreateUserSubscriptionRequest) (models.UserSubscription, error)
	Patch(ctx context.Context, id uint32, subscription models.PatchUserSubscriptionRequest) (models.UserSubscription, error)
	GetSubscription(ctx context.Context, id uint32) (models.UserSubscription, error)
	DeleteSubscription(ctx context.Context, id uint32) (int64, error)
	GetSubscriptions(ctx context.Context, id string) ([]models.UserSubscription, error)
	GetSubscriptionsTotal(ctx context.Context, filter models.SubscriptionsTotalFilter) (models.SubscriptionsTotal, error)
}

type userSubscription struct {
	db     *database.DB
	logger logger.Logger
}

func NewUserSubscriptionRepository(db *database.DB, logger logger.Logger) UserSubscription {
	return &userSubscription{
		db:     db,
		logger: logger,
	}
}

func (r *userSubscription) Create(ctx context.Context, subscription models.CreateUserSubscriptionRequest) (models.CreateUserSubscriptionResponse, error) {
	startDate := subscription.StartDate.Truncate(24 * time.Hour)

	var stopDate *time.Time
	if subscription.StopDate != nil {
		temp := subscription.StopDate.Truncate(24 * time.Hour)
		stopDate = &temp
	}

	var response models.CreateUserSubscriptionResponse
	query := `
		insert into user_subscriptions (user_id, service_name, price, start_date, stop_date)
		values ($1, $2, $3, $4, $5)
		returning id, created_at, updated_at
	`
	err := r.db.Pool(ctx).QueryRow(
		ctx, query,
		subscription.UserID, subscription.ServiceName, subscription.Price, startDate, stopDate,
	).Scan(&response.ID, &response.CreatedAt, &response.UpdatedAt)

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *userSubscription) Update(ctx context.Context, id uint32, subscription models.CreateUserSubscriptionRequest) (models.UserSubscription, error) {
	startDate := subscription.StartDate.Truncate(24 * time.Hour)

	var stopDate *time.Time
	if subscription.StopDate != nil {
		temp := subscription.StopDate.Truncate(24 * time.Hour)
		stopDate = &temp
	}

	var response models.UserSubscription
	query := `
		update user_subscriptions
		set user_id = $1,
		    service_name = $2,
		    price = $3,
		    start_date = $4,
		    stop_date = $5
		where id = $6
		returning id, user_id, service_name, price, start_date, stop_date, created_at, updated_at
	`
	err := r.db.Pool(ctx).QueryRow(
		ctx, query,
		subscription.UserID, subscription.ServiceName, subscription.Price, startDate, stopDate, id,
	).Scan(
		&response.ID,
		&response.UserID,
		&response.ServiceName,
		&response.Price,
		&response.StartDate,
		&response.StopDate,
		&response.CreatedAt,
		&response.UpdatedAt)

	if err != nil {
		r.logger.Warn("failed to update user subscription", "error", err)
		return response, err
	}

	return response, nil
}

func (r *userSubscription) Patch(ctx context.Context, id uint32, subscription models.PatchUserSubscriptionRequest) (models.UserSubscription, error) {
	startDate := subscription.StartDate.Truncate(24 * time.Hour)

	var stopDate *time.Time
	if subscription.StopDate != nil {
		temp := subscription.StopDate.Truncate(24 * time.Hour)
		stopDate = &temp
	}

	var response models.UserSubscription
	query := `
		update user_subscriptions
		set price = $1,
		    start_date = $2,
		    stop_date = $3
		where id = $4
		returning id, user_id, service_name, price, start_date, stop_date, created_at, updated_at
	`
	err := r.db.Pool(ctx).QueryRow(
		ctx, query, subscription.Price, startDate, stopDate, id,
	).Scan(
		&response.ID,
		&response.UserID,
		&response.ServiceName,
		&response.Price,
		&response.StartDate,
		&response.StopDate,
		&response.CreatedAt,
		&response.UpdatedAt)

	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *userSubscription) GetSubscription(ctx context.Context, id uint32) (models.UserSubscription, error) {
	var response models.UserSubscription
	query := `
		select id, user_id, service_name, price, start_date, stop_date, created_at, updated_at
		from user_subscriptions
		where id = $1
	`
	err := r.db.Pool(ctx).QueryRow(ctx, query, id).Scan(
		&response.ID,
		&response.UserID,
		&response.ServiceName,
		&response.Price,
		&response.StartDate,
		&response.StopDate,
		&response.CreatedAt,
		&response.UpdatedAt,
	)

	if err != nil {
		return response, err
	}

	return response, err
}

func (r *userSubscription) DeleteSubscription(ctx context.Context, id uint32) (int64, error) {
	query := `delete from user_subscriptions where id = $1`
	cmdTag, err := r.db.Pool(ctx).Exec(ctx, query, id)
	if err != nil {
		return 0, err
	}

	return cmdTag.RowsAffected(), nil
}

func (r *userSubscription) GetSubscriptions(ctx context.Context, id string) ([]models.UserSubscription, error) {
	query := `
		select id, user_id, service_name, price, start_date, stop_date, created_at, updated_at
		from user_subscriptions
		where user_id = $1
	`
	rows, err := r.db.Pool(ctx).Query(ctx, query, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []models.UserSubscription
	for rows.Next() {
		s := models.UserSubscription{}
		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.ServiceName,
			&s.Price,
			&s.StartDate,
			&s.StopDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return ret, err
		}
		ret = append(ret, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *userSubscription) GetSubscriptionsTotal(ctx context.Context, filter models.SubscriptionsTotalFilter) (models.SubscriptionsTotal, error) {
	// Compose where clause
	var args []interface{}
	argIdx := 1
	whereClauses := []string{}

	args = append(args, filter.DateFrom)
	argIdx++

	args = append(args, filter.DateTo)
	argIdx++

	if filter.UserID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("t.user_id = $%d", argIdx))
		args = append(args, *filter.UserID)
		argIdx++
	}

	if filter.ServiceName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("t.service_name = $%d", argIdx))
		args = append(args, *filter.ServiceName)
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "where " + strings.Join(whereClauses, " and ")
	}

	query := fmt.Sprintf(`
select coalesce(sum(t.price * months_between(from_date, till_date)),0) 
from (
	select 
	    user_id,
	    service_name,
	    price,
	    greatest(start_date, $1) as from_date,
	    least(coalesce(stop_date, $2), now()::date) as till_date 
	from user_subscriptions
	where greatest(start_date, $1) between $1 and $2
) t %s`, whereSQL)

	var total models.SubscriptionsTotal
	err := r.db.Pool(ctx).QueryRow(ctx, query, args...).Scan(&total.Total)
	if err != nil {
		return total, err
	}

	return total, nil
}
