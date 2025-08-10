package repositories

import (
	domain "blog_starter_project_g66/Domain"
	"blog_starter_project_g66/config"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserOTPRepository struct {
	collection *mongo.Collection
}

func NewUserOTPRepository(dbClient *MongoDBClient) *UserOTPRepository {
	userDB := config.USER_DB
	db := dbClient.Client.Database(userDB)
	return &UserOTPRepository{
		collection: db.Collection(config.USER_OTP_COLLECTION_NAME),
	}
}

func (r *UserOTPRepository) StoreOTP(entry domain.UserUnverified) error {
	filter := bson.M{"email": entry.Email}
	update := bson.M{"$set": entry}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func (r *UserOTPRepository) FindOTP(email string) (*domain.UserUnverified, error) {
	var entry domain.UserUnverified
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *UserOTPRepository) DeleteOTP(email string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"email": email})
	return err
}
