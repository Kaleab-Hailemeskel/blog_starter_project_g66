package repositories

import (
	conv "blog_starter_project_g66/Delivery/converter"
	domain "blog_starter_project_g66/Domain"
	"blog_starter_project_g66/config"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	collection *mongo.Collection
	Contxt     context.Context
	Client     *mongo.Client
}

func NewUserRepository() domain.IUserRepository {
	connection, err := Connect()
	mainBlogDbName := config.BLOG_DB
	UserDataBaseName := config.USER_DB
	UserCollectionName := config.USER_COLLECTION_NAME
	if err != nil {
		log.Fatal("can't initailize ", mainBlogDbName, " Database")
	}

	collection := connection.Client.Database(UserDataBaseName).Collection(UserCollectionName)

	return &UserRepository{
		collection: collection,
		Contxt:     context.TODO(),
		Client:     connection.Client,
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	log.Println("==>   inside create")
	userDTO := conv.ChangeToDTOUser(user)
	_, err := r.collection.InsertOne(r.Contxt, userDTO)
	return err
}
func (r *UserRepository) CheckUserExistance(userEmail string) bool {
	filter := bson.M{"email": userEmail}
	err := r.collection.FindOne(r.Contxt, filter).Err()
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
func (r *UserRepository) GetUserByID(userID string) (*domain.UserDTO, error) {
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
	return r.Client.Disconnect(r.Contxt)
}
func (r *UserRepository) UpdatePassword(userID, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userid": userID}
	update := bson.M{"$set": bson.M{"password": hashedPassword}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
func (r *UserRepository) UpdateRole(email, role string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"role": role}}

	_, err := r.collection.UpdateOne(r.Contxt, filter, update)
	return err
}
func (repo *UserRepository) UpdateUserByEmail(email string, dto *domain.UpdateProfileDTO) (*domain.UserDTO, error) {
	filter := bson.M{"email": email}

	updateFields := bson.M{}
	if dto.UserName != "" {
		updateFields["username"] = dto.UserName
	}
	if dto.PersonalBio != "" {
		updateFields["personal_bio"] = dto.PersonalBio
	}
	if dto.ProfilePic != "" {
		updateFields["profile_pic"] = dto.ProfilePic
	}
	if dto.PhoneNum != "" {
		updateFields["phone_num"] = dto.PhoneNum
	}
	if dto.TelegramHandle != "" {
		updateFields["telegram_handle"] = dto.TelegramHandle
	}
	if dto.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updateFields["password"] = string(hashedPassword)
	}

	update := bson.M{"$set": updateFields}

	var updatedUser domain.UserDTO
	err := repo.collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedUser)

	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}
