package core

import (
	"context"
	"errors"

	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const appSettingsCollection = "app_settings"

// AppSettingsRepository - MongoDB-backed repository for AppSettings.
type AppSettingsRepository struct {
	db *mongo.Database
}

// NewAppSettingsRepository - returns new AppSettingsRepository instance.
func NewAppSettingsRepository() *AppSettingsRepository {
	return &AppSettingsRepository{
		db: persist.GetDatabase(),
	}
}

// Insert - inserts AppSettings into DB.
func (ar *AppSettingsRepository) Insert(settings *models.AppSettings) error {
	collection := ar.db.Collection(appSettingsCollection)

	currentSettings, _ := ar.GetAppSettings()
	if currentSettings != nil {
		return errors.New("AppSettings already exist")
	}

	settings.CreatedAt = models.NullableTimeNow()
	settings.UpdatedAt = models.NullableTimeNow()

	result, err := collection.InsertOne(context.TODO(), settings)
	if err != nil {
		return err
	}

	settings.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// GetAppSettings - returns AppSettings.
func (ar *AppSettingsRepository) GetAppSettings() (*models.AppSettings, error) {
	collection := ar.db.Collection(appSettingsCollection)

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	result := []*models.AppSettings{}

	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("AppSettings have not been found")
	}

	// There can be only one settings document in the DB.
	settings := result[0]

	return settings, nil
}

func (ar *AppSettingsRepository) Update(settings *models.AppSettings) error {
	collection := ar.db.Collection(appSettingsCollection)

	currentSettings, err := ar.GetAppSettings()
	if err != nil {
		return err
	}

	settings.UpdatedAt = models.NullableTimeNow()

	payload := bson.D{{"$set", settings}}

	if _, err := collection.UpdateByID(context.TODO(), currentSettings.ID, payload); err != nil {
		return err
	}

	return nil
}
