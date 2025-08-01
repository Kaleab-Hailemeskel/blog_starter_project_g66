package repositories

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "blog_starter_project_g66/Domain"
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