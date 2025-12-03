package repositories

import (
	"context"
	"errors"

	"finanzas-backend/internal/profile/domain/model/entities"
	"finanzas-backend/internal/profile/domain/model/valueobjects"
	domain_repos "finanzas-backend/internal/profile/domain/repositories"
	"finanzas-backend/internal/profile/infrastructure/persistence/models"
	"finanzas-backend/internal/shared/infrastructure/security"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type profileRepositoryImpl struct {
	db                *gorm.DB
	encryptionService *security.EncryptionService
}

func NewProfileRepository(db *gorm.DB, encryptionService *security.EncryptionService) domain_repos.ProfileRepository {
	return &profileRepositoryImpl{
		db:                db,
		encryptionService: encryptionService,
	}
}

func (r *profileRepositoryImpl) Save(ctx context.Context, profile *entities.Profile) error {
	// Generate new UUID for profile
	profileID, err := valueobjects.NewProfileID(uuid.New())
	if err != nil {
		return err
	}
	profile.SetID(profileID)

	// Convert to model
	model := models.FromEntity(profile)

	// Encrypt DNI in the model
	encryptedDNI, err := r.encryptionService.Encrypt(profile.DNI().Value())
	if err != nil {
		return err
	}
	model.DNIEncrypted = encryptedDNI

	return r.db.WithContext(ctx).Create(model).Error
}

func (r *profileRepositoryImpl) Update(ctx context.Context, profile *entities.Profile) error {
	// Convert to model
	model := models.FromEntity(profile)

	// Encrypt DNI in the model
	encryptedDNI, err := r.encryptionService.Encrypt(profile.DNI().Value())
	if err != nil {
		return err
	}
	model.DNIEncrypted = encryptedDNI

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

	profile, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	// Decrypt DNI
	if err := r.decryptProfileDNI(profile); err != nil {
		return nil, err
	}

	return profile, nil
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

	profile, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	// Decrypt DNI
	if err := r.decryptProfileDNI(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepositoryImpl) FindByDNI(ctx context.Context, dni string) (*entities.Profile, error) {
	// Encrypt DNI for search
	encryptedDNI, err := r.encryptionService.Encrypt(dni)
	if err != nil {
		return nil, err
	}

	var model models.ProfileModel
	err = r.db.WithContext(ctx).Where("dni_encrypted = ?", encryptedDNI).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	profile, err := model.ToEntity()
	if err != nil {
		return nil, err
	}

	// Decrypt DNI
	if err := r.decryptProfileDNI(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

func (r *profileRepositoryImpl) ExistsByUserID(ctx context.Context, userID valueobjects.UserID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ProfileModel{}).Where("user_id = ?", userID.Value()).Count(&count).Error
	return count > 0, err
}

func (r *profileRepositoryImpl) ExistsByDNI(ctx context.Context, dni string) (bool, error) {
	// Encrypt DNI for search
	encryptedDNI, err := r.encryptionService.Encrypt(dni)
	if err != nil {
		return false, err
	}

	var count int64
	err = r.db.WithContext(ctx).Model(&models.ProfileModel{}).Where("dni_encrypted = ?", encryptedDNI).Count(&count).Error
	return count > 0, err
}

// decryptProfileDNI decrypts the DNI field of a profile
func (r *profileRepositoryImpl) decryptProfileDNI(profile *entities.Profile) error {
	encryptedValue := profile.DNI().EncryptedValue()
	if encryptedValue == "" {
		return nil // No encrypted value to decrypt
	}

	decryptedDNI, err := r.encryptionService.Decrypt(encryptedValue)
	if err != nil {
		return err
	}

	// Create new DNI with decrypted value
	dni, err := valueobjects.NewDNI(decryptedDNI)
	if err != nil {
		return err
	}

	// Update the profile entity with decrypted DNI
	profile.SetDNI(dni)
	return nil
}
