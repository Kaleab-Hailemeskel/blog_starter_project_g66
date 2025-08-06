package repositories

import (
	domain "blog_starter_project_g66/Domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshTokenRepository struct {
	Collection *mongo.Collection
}

func NewRefreshTokenRepository(dbClient *MongoDBClient) *RefreshTokenRepository {
	db := dbClient.Client.Database("user_db")
	return &RefreshTokenRepository{
		Collection: db.Collection("refresh_token"),
	}
}

func (r *RefreshTokenRepository) Save(token *domain.RefreshToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": token.UserID}
	update := bson.M{"$set": token}
	opts := options.Update().SetUpsert(true)

	_, err := r.Collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *RefreshTokenRepository) GetByToken(token string) (*domain.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rt domain.RefreshToken
	err := r.Collection.FindOne(ctx, bson.M{"token": token}).Decode(&rt)
	return &rt, err
}

func (r *RefreshTokenRepository) Delete(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Collection.DeleteOne(ctx, bson.M{"token": token})
	return err
}
