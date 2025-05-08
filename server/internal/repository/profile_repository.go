package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eduhub/server/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4" // For pgx.ErrNoRows
)

const profileTable = "profiles"

type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile *models.Profile) error
	GetProfileByUserID(ctx context.Context, userID string) (*models.Profile, error)
	GetProfileByID(ctx context.Context, profileID int) (*models.Profile, error)
	UpdateProfile(ctx context.Context, profile *models.Profile) error
	// Optional: DeleteProfile(ctx context.Context, profileID int) error
}

type profileRepository struct {
	DB *DB
}

func NewProfileRepository(db *DB) ProfileRepository {
	return &profileRepository{DB: db}
}

func (r *profileRepository) CreateProfile(ctx context.Context, profile *models.Profile) error {
	now := time.Now()
	if profile.JoinedAt.IsZero() {
		profile.JoinedAt = now
	}
	profile.LastActive = now
	profile.CreatedAt = now
	profile.UpdatedAt = now

	if profile.Preferences == nil {
		profile.Preferences = make(models.JSONMap)
	}
	if profile.SocialLinks == nil {
		profile.SocialLinks = make(models.JSONMap)
	}

	query := r.DB.SQ.Insert(profileTable).
		Columns(
			"user_id", "college_id", "bio", "profile_image", "phone_number",
			"address", "date_of_birth", "joined_at", "last_active",
			"preferences", "social_links", "created_at", "updated_at",
		).
		Values(
			profile.UserID, profile.CollegeID, profile.Bio, profile.ProfileImage, profile.PhoneNumber,
			profile.Address, profile.DateOfBirth, profile.JoinedAt, profile.LastActive,
			profile.Preferences, profile.SocialLinks, profile.CreatedAt, profile.UpdatedAt,
		).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateProfile: failed to build query: %w", err)
	}

	err = r.DB.Pool.QueryRow(ctx, sql, args...).Scan(&profile.ID)
	if err != nil {
		return fmt.Errorf("CreateProfile: failed to execute query or scan ID: %w", err)
	}
	return nil
}

func (r *profileRepository) GetProfileByUserID(ctx context.Context, userID string) (*models.Profile, error) {
	profile := &models.Profile{}
	queryFields := []string{"id", "user_id", "college_id", "bio", "profile_image", "phone_number", "address", "date_of_birth", "joined_at", "last_active", "preferences", "social_links", "created_at", "updated_at"}
	query := r.DB.SQ.Select(queryFields...).
		From(profileTable).
		Where(squirrel.Eq{"user_id": userID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetProfileByUserID: failed to build query: %w", err)
	}

	err = pgxscan.Get(ctx, r.DB.Pool, profile, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetProfileByUserID: profile for user ID %s not found", userID)
		}
		return nil, fmt.Errorf("GetProfileByUserID: failed to execute query or scan: %w", err)
	}
	return profile, nil
}

func (r *profileRepository) GetProfileByID(ctx context.Context, profileID int) (*models.Profile, error) {
	profile := &models.Profile{}
	queryFields := []string{"id", "user_id", "college_id", "bio", "profile_image", "phone_number", "address", "date_of_birth", "joined_at", "last_active", "preferences", "social_links", "created_at", "updated_at"}
	query := r.DB.SQ.Select(queryFields...).
		From(profileTable).
		Where(squirrel.Eq{"id": profileID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetProfileByID: failed to build query: %w", err)
	}

	err = pgxscan.Get(ctx, r.DB.Pool, profile, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetProfileByID: profile with ID %d not found", profileID)
		}
		return nil, fmt.Errorf("GetProfileByID: failed to execute query or scan: %w", err)
	}
	return profile, nil
}

func (r *profileRepository) UpdateProfile(ctx context.Context, profile *models.Profile) error {
	now := time.Now()
	profile.LastActive = now
	profile.UpdatedAt = now

	if profile.Preferences == nil {
		profile.Preferences = make(models.JSONMap)
	}
	if profile.SocialLinks == nil {
		profile.SocialLinks = make(models.JSONMap)
	}

	query := r.DB.SQ.Update(profileTable).
		Set("college_id", profile.CollegeID).
		Set("bio", profile.Bio).
		Set("profile_image", profile.ProfileImage).
		Set("phone_number", profile.PhoneNumber).
		Set("address", profile.Address).
		Set("date_of_birth", profile.DateOfBirth).
		Set("last_active", profile.LastActive).
		Set("preferences", profile.Preferences).
		Set("social_links", profile.SocialLinks).
		Set("updated_at", profile.UpdatedAt).
		Where(squirrel.Eq{"id": profile.ID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateProfile: failed to build query: %w", err)
	}

	commandTag, err := r.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateProfile: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateProfile: no profile found with ID %d, or no changes made", profile.ID)
	}
	return nil
}
