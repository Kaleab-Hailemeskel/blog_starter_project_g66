package repositories

import (
	domain "blog_starter_project_g66/Domain"
	"context"
	"errors"
	"time"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "golang.org/x/crypto/bcrypt"
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
func (r *UserRepository) FindByEmail(email string) (*domain.UserDTO, error) {
	var user domain.UserDTO
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil

}

func (r *UserRepository) GetUserByID(userID string) (*domain.UserDTO, error){
	var user domain.UserDTO
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepository) CloseDataBase() error {
	return r.collection.Database().Client().Disconnect(context.TODO())
}

func (r *UserRepository) UpdatePassword(userID, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userid": userID}
	update := bson.M{"$set": bson.M{"password": hashedPassword}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
} 

// func (r *UserRepository) CreateSuperAdmin() error{

// 	email := "superadmin@gmail.com"
// 	username := "superadmin"
// 	password := "123456789ADm@"

// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	user := bson.M{
// 		"username": username,
// 		"email": email,
// 		"password": string(hashedPassword),
// 		"role": "SUPER_ADMIN",
// 	}
// 	_, err := r.collection.InsertOne(context.TODO(), user)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *UserRepository) UpdateRole(email, role string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"role":role}}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}