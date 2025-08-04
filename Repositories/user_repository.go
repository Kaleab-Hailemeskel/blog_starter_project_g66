package repositories

import (
	"blog_starter_project_g66/Domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
   collection *mongo.Collection
}


func NewUserRepository(dbClient *MongoDBClient) *UserRepository {
    db := dbClient.Client.Database("user_db") 
    return &UserRepository{
        collection: db.Collection("users"), 
    }
}

func (r *UserRepository) Create(user *domain.User) error {
    _, err := r.collection.InsertOne(context.TODO(), user)
    return err
}

func (r *UserRepository) CheckUserExistance(userEmail string) bool {
	filter := bson.M{"email": userEmail}
	err := r.collection.FindOne(context.TODO(), filter).Err()
	return err == nil // true means user found


}

func (r *UserRepository) CloseDataBase() error{
   return r.collection.Database().Client().Disconnect(context.TODO())
}