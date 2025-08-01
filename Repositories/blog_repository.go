package repositories

import (
	"blog_starter_project_g66/Delivery/controllers"
	conv "blog_starter_project_g66/Delivery/converter"
	domain "blog_starter_project_g66/Domain"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName   = "blogbd_test"
	collName = "blogCollection"
	pageSize = 10
)

type BlogDB struct {
	Coll   mongo.Collection
	Contxt context.Context
	Client *mongo.Client
}

func NewBlogDataBaseService() domain.IBlogRepository {
	connection, err := Connect()

	if err != nil {
		log.Fatal("can't initailize ", dbName, " Database")
	}

	collection := connection.Client.Database(dbName).Collection(collName)

	return &BlogDB{
		Coll:   *collection,
		Contxt: context.TODO(),
		Client: connection.Client,
	}

}

func (bldb *BlogDB) CreateBlog(blog *domain.Blog, userID primitive.ObjectID) error {

	blog.LastUpdate = time.Now()
	blogDTO := conv.ChangeToDTOBlog(blog) // Convert domain.Blog to controllers.BlogDTO
	blogDTO.OwnerID = userID
	_, err := bldb.Coll.InsertOne(bldb.Contxt, blogDTO) // Insert the DTO into the collection
	if err != nil {
		return fmt.Errorf("error creating blog: %w", err)
	}
	return nil
}

func (bldb *BlogDB) DeleteBlogByID(blogID primitive.ObjectID) error {
	filter := bson.M{"_id": blogID}
	result, err := bldb.Coll.DeleteOne(bldb.Contxt, filter)
	if err != nil {
		return fmt.Errorf("error deleting blog with ID %s: %w", blogID, err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no blog found with ID %s to delete", blogID)
	}
	return nil
}

func (bldb *BlogDB) UpdateBlogByID(blogID primitive.ObjectID, updatedBlog *domain.Blog) error {
	// I don't want to check the existance of the blog before updating, b/c the updateOne will tell us if there weren't any changes mamde afterall
	filter := bson.M{"_id": blogID}

	// Create an update document using $set to update specific fields.
	updateDoc := bson.M{
		"$set": bson.M{
			"title":       updatedBlog.Title,
			"tags":        updatedBlog.Tags,
			"author":      updatedBlog.Author,
			"description": updatedBlog.Description,
			"last_update": time.Now(), // Update timestamp on every update
		},
	}

	result, err := bldb.Coll.UpdateOne(bldb.Contxt, filter, updateDoc)
	if err != nil {
		return fmt.Errorf("error updating blog with ID %s: %w", blogID, err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no blog found with ID %s to update or no changes were made", blogID)
	}
	return nil
}

func (bldb *BlogDB) GetAllBlogsByFilter(url_filter *domain.Filter, pageNumber int) ([]*controllers.BlogDTO, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	skip := int64((pageNumber - 1) * pageSize)
	limit := int64(pageSize)

	filter := bson.M{}
	{
		if url_filter.Tag != "" {
			// Use $regex for case-insensitive partial match on tags array
			filter["tags"] = bson.M{"$regex": primitive.Regex{Pattern: url_filter.Tag, Options: "i"}}
		}
		if url_filter.AuthorName != "" {
			filter["author"] = bson.M{"$regex": primitive.Regex{Pattern: url_filter.AuthorName, Options: "i"}}
		}
		if url_filter.Title != "" {
			filter["title"] = bson.M{"$regex": primitive.Regex{Pattern: url_filter.Title, Options: "i"}}
		}
		if !url_filter.AfterDate.IsZero() {
			filter["last_update"] = bson.M{"$gte": url_filter.AfterDate}
		}

	}

	findOptions := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "last_update", Value: -1}})

	cursor, err := bldb.Coll.Find(bldb.Contxt, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting filtered blogs: %w", err)
	}
	defer cursor.Close(bldb.Contxt)

	var blogs []*controllers.BlogDTO
	for cursor.Next(bldb.Contxt) {
		var blogDTO controllers.BlogDTO // Change target type to DTO for decoding

		if err := cursor.Decode(&blogDTO); err != nil {
			return nil, fmt.Errorf("error decoding filtered blog DTO: %w", err)
		}

		blogs = append(blogs, &blogDTO)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through filtered blog cursor: %w", err)
	}
	return blogs, nil
}

func (bldb *BlogDB) CheckBlogExistance(blogID primitive.ObjectID) bool {

	filter := bson.M{"_id": blogID}
	count, err := bldb.Coll.CountDocuments(bldb.Contxt, filter)
	if err != nil {
		log.Printf("Error checking blog existence for ID %s: %v", blogID, err)
		return false
	}
	return count > 0
}

func (bldb *BlogDB) CloseDataBase() error {
	if bldb.Client == nil {
		return nil // Nothing to close
	}
	if err := bldb.Client.Disconnect(bldb.Contxt); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w", err)
	}
	log.Println("Disconnected from MongoDB.")
	return nil
}
