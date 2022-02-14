package repositories

import (
	"context"
	"time"

	"github.com/el-mike/dogecrack/shepherd/internal/common/models"
	"github.com/el-mike/dogecrack/shepherd/internal/persist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection = "users"

// - UserRepository - MongoDB-backed repository for users.
type UserRepository struct {
	db *mongo.Database
}

// NewUserRepository - returns new UserRepository instance.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: persist.GetDatabase(),
	}
}

// GetByName - gets a single User with given name.
func (ur *UserRepository) GetByName(name string) (*models.User, error) {
	collection := ur.db.Collection(usersCollection)

	filter := bson.M{"name": name}

	result := collection.FindOne(context.TODO(), filter)

	user := &models.User{}

	if err := result.Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetById - returns a single User with given id.
func (ur *UserRepository) GetById(id string) (*models.User, error) {
	collection := ur.db.Collection(usersCollection)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}

	result := collection.FindOne(context.TODO(), filter)

	user := &models.User{}

	if err := result.Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Insert - adds a single User to the DB.
func (ur *UserRepository) Insert(user *models.User) error {
	collection := ur.db.Collection(usersCollection)

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}
