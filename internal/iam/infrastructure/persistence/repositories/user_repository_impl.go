package repositories

import (
	"context"
	"errors"
	"finanzas-backend/internal/iam/domain/model/entities"
	"finanzas-backend/internal/iam/domain/model/valueobjects"
	domain_repos "finanzas-backend/internal/iam/domain/repositories"
	"finanzas-backend/internal/iam/infrastructure/persistence/models"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain_repos.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *entities.User) error {
	model := models.FromEntity(user)

	if user.ID().IsZero() {
		// Create - Generate new UUID
		model.ID = valueobjects.GenerateUserID().Value()
		if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
			return err
		}
		userID, _ := valueobjects.NewUserID(model.ID)
		user.SetID(userID)
	} else {
		// Update
		if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id valueobjects.UserID) (*entities.User, error) {
	return r.FindByIDValue(ctx, id.String())
}

func (r *userRepositoryImpl) FindByIDValue(ctx context.Context, id string) (*entities.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity()
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model models.UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToEntity()
}

func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.UserModel{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
