package repositories

import (
	"blog_starter_project_g66/Domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongouserRepository struct {
	collection *mongo.Collection
}

func (r *MongouserRepository) FindByEmail(email string) (*domain.UserDTO, error) {
	var user domain.UserDTO
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil

}

func (r *MongouserRepository) UpdatePassword(userID, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userid": userID}
	update := bson.M{"$set": bson.M{"password": hashedPassword}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
} 
