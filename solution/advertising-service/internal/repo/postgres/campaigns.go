package postgres

import (
	"advertising/advertising-service/internal/dto"
	"advertising/advertising-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CampaignsRepo struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
}

func NewCampaignsRepo(db *sqlx.DB) *CampaignsRepo {
	return &CampaignsRepo{
		db: db,
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (cr *CampaignsRepo) CreateCampaign(ctx context.Context, advertiserId uuid.UUID, data dto.CampaignData) (uuid.UUID, error) {
	op := "CampaignsRepo.CreateCampaign"

	columns := []string{
		"advertiser_id", "impressions_limit", "clicks_limit",
		"cost_per_impression", "cost_per_click",
		"ad_title", "ad_text", "start_date", "end_date",
	}
	values := []any{
		advertiserId, data.ImpressionsLimit, data.ClicksLimit,
		data.CostPerImpression, data.CostPerClick,
		data.AdTitle, data.AdText, data.StartDate, data.EndDate,
	}

	if data.Gender != nil {
		columns = append(columns, "gender")
		values = append(values, *data.Gender)
	}
	if data.AgeFrom != nil {
		columns = append(columns, "age_from")
		values = append(values, *data.AgeFrom)
	}
	if data.AgeTo != nil {
		columns = append(columns, "age_to")
		values = append(values, *data.AgeTo)
	}
	if data.Location != nil {
		columns = append(columns, "location")
		values = append(values, *data.Location)
	}

	query, args, err := cr.sq.
		Insert("campaigns").
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var id uuid.UUID
	if err := cr.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23503":
				return uuid.UUID{}, models.ErrAdvertiserNotFound
			}
		}
		return uuid.UUID{}, fmt.Errorf("%s: db.QueryRowContext: %w", op, err)
	}

	return id, nil
}

func (cr *CampaignsRepo) GetCampaignById(ctx context.Context, campaignId uuid.UUID) (models.Campaign, error) {
	op := "CampaignsRepo.GetCampaignById"

	query, args, err := cr.sq.
		Select(
			"id", "advertiser_id", "impressions_limit", "clicks_limit",
			"cost_per_impression", "cost_per_click",
			"ad_title", "ad_text", "ad_image_url",
			"start_date", "end_date",
			"gender", "age_from", "age_to", "location",
		).From("campaigns").
		Where(sq.Eq{"id": campaignId}).
		ToSql()
	if err != nil {
		return models.Campaign{}, fmt.Errorf("%s: build query: %w", op, err)
	}

	var campaign models.Campaign
	if err := cr.db.GetContext(ctx, &campaign, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Campaign{}, models.ErrCampaignNotFound
		}
		return models.Campaign{}, fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	return campaign, nil
}

func (cr *CampaignsRepo) ListCampaignsForAdvertiser(ctx context.Context, advertiserId uuid.UUID, params dto.PaginationParams) ([]models.Campaign, error) {
	op := "CampaignsRepo.ListCampaignsForAdvertiser"

	qb := cr.sq.
		Select(
			"id", "advertiser_id", "impressions_limit", "clicks_limit",
			"cost_per_impression", "cost_per_click",
			"ad_title", "ad_text", "ad_image_url",
			"start_date", "end_date",
			"gender", "age_from", "age_to", "location",
		).From("campaigns").
		Where(sq.Eq{"advertiser_id": advertiserId}).
		Limit(uint64(params.Size)).
		Offset(uint64(params.Page-1) * uint64(params.Size)).
		OrderBy("created_at DESC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: build query: %w", op, err)
	}

	campaigns := []models.Campaign{}
	if err := cr.db.SelectContext(ctx, &campaigns, query, args...); err != nil {
		return nil, fmt.Errorf("%s: db.SelectContext: %w", op, err)
	}

	return campaigns, nil
}

func (cr *CampaignsRepo) UpdateCampaign(ctx context.Context, campaignId uuid.UUID, data dto.CampaignData) error {
	op := "CampaignsRepo.UpdateCampaign"

	query, args, err := cr.sq.
		Update("campaigns").
		Set("impressions_limit", data.ImpressionsLimit).
		Set("clicks_limit", data.ClicksLimit).
		Set("cost_per_impression", data.CostPerImpression).
		Set("cost_per_click", data.CostPerClick).
		Set("ad_title", data.AdTitle).
		Set("ad_text", data.AdText).
		Set("start_date", data.StartDate).
		Set("end_date", data.EndDate).
		Set("gender", data.Gender).
		Set("age_from", data.AgeFrom).
		Set("age_to", data.AgeTo).
		Set("location", data.Location).
		Where(sq.Eq{"id": campaignId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	res, err := cr.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: db.GetContext: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: res.RowsAffected: %w", op, err)
	}

	if rowsAffected == 0 {
		return models.ErrCampaignNotFound
	}

	return nil
}

func (cr *CampaignsRepo) SetCampaignAdImageUrl(ctx context.Context, campaignId uuid.UUID, adImageUrl *string) error {
	op := "CampaignsRepo.SetCampaignAdImageUrl"

	query, args, err := cr.sq.
		Update("campaigns").
		Set("ad_image_url", adImageUrl).
		Where(sq.Eq{"id": campaignId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	res, err := cr.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: res.RowsAffected: %w", op, err)
	}

	if rowsAffected == 0 {
		return models.ErrCampaignNotFound
	}

	return nil
}

func (cr *CampaignsRepo) DeleteCampaign(ctx context.Context, campaignId uuid.UUID) error {
	op := "CampaignsRepo.DeleteCampaigns"

	query, args, err := cr.sq.
		Delete("campaigns").
		Where(sq.Eq{"id": campaignId}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: build query: %w", op, err)
	}

	res, err := cr.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: db.ExecContext: %w", op, err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: res.RowsAffected: %w", op, err)
	}

	if affectedRows == 0 {
		return models.ErrCampaignNotFound
	}

	return nil
}
