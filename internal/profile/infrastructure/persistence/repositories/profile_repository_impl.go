package repositories

import (
	"context"
	"errors"

	domain_repos "finanzas-backend/internal/profile/domain/repositories"
	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
	"finanzas-backend/internal/profile/infrastructure/persistence/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type profileRepositoryImpl struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) domain_repos.ProfileRepository {
	return &profileRepositoryImpl{db: db}
}

func (r *profileRepositoryImpl) Save(ctx context.Context, profile *entities.Profile) error {
	// Generate new UUID for profile
	profileID, err := valueobjects.NewProfileID(uuid.New())
	if err != nil {
		return err
	}
	profile.SetID(profileID)

	model := models.FromEntity(profile)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *profileRepositoryImpl) Update(ctx context.Context, profile *entities.Profile) error {
	model := models.FromEntity(profile)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *profileRepositoryImpl) FindByID(ctx context.Context, id valueobjects.ProfileID) (*entities.Profile, error) {
	var model models.ProfileModel
	err := r.db.WithContext(ctx).Where("id = ?", id.Value()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity()
}

func (r *profileRepositoryImpl) FindByUserID(ctx context.Context, userID valueobjects.UserID) (*entities.Profile, error) {
	var model models.ProfileModel
	err := r.db.WithContext(ctx).Where("user_id = ?", userID.Value()).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity()
}

func (r *profileRepositoryImpl) FindByDNI(ctx context.Context, dni string) (*entities.Profile, error) {
	var model models.ProfileModel
	err := r.db.WithContext(ctx).Where("dni = ?", dni).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity()
}

func (r *profileRepositoryImpl) ExistsByUserID(ctx context.Context, userID valueobjects.UserID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ProfileModel{}).Where("user_id = ?", userID.Value()).Count(&count).Error
	return count > 0, err
}

func (r *profileRepositoryImpl) ExistsByDNI(ctx context.Context, dni string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ProfileModel{}).Where("dni = ?", dni).Count(&count).Error
	return count > 0, err
}
