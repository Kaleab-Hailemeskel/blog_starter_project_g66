package repositories

import (
	"blog_starter_project_g66/Delivery/controllers"
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
	dbName   = "blogbd"
	collName = "blogCollection"
	pageSize = 10 // Define your page size
)

type BlogDB struct {
	Coll   mongo.Collection
	Contxt context.Context
	Client *mongo.Client
}

func NewBlogDataBaseService() domain.IUserRepository {
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

func (bldb *BlogDB) CreateBlog(blog *domain.Blog) error {
	// Set LastUpdate timestamp
	blog.LastUpdate = time.Now()
	// If BlogID is empty, generate a new ObjectID and assign its hex string to BlogID
	if blog.BlogID == "" {
		blog.BlogID = primitive.NewObjectID().Hex()
	}
	blogDTO := controllers.ChangeToDTOBlog(blog)
	// Convert BlogID string to primitive.ObjectID for insertion as _id
	objID, err := primitive.ObjectIDFromHex(blogDTO.BlogID)
	if err != nil {
		return fmt.Errorf("invalid BlogID format: %w", err)
	}

	// Create a BSON document for insertion
	doc := bson.M{
		"_id":         objID,
		"owner_id":    blogDTO.OwnerID,
		"title":       blogDTO.Title,
		"tags":        blogDTO.Tags,
		"author":      blogDTO.Author,
		"description": blogDTO.Description,
		"last_update": blogDTO.LastUpdate,
	}

	_, err = bldb.Coll.InsertOne(bldb.Contxt, doc)
	if err != nil {
		return fmt.Errorf("error creating blog: %w", err)
	}
	return nil
}

func (bldb *BlogDB) DeleteBlogByID(blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return fmt.Errorf("invalid BlogID format: %w", err)
	}

	filter := bson.M{"_id": objID}
	result, err := bldb.Coll.DeleteOne(bldb.Contxt, filter)
	if err != nil {
		return fmt.Errorf("error deleting blog with ID %s: %w", blogID, err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no blog found with ID %s to delete", blogID)
	}
	return nil
}

func (bldb *BlogDB) UpdateBlogByID(blogID string, updatedBlog *domain.Blog) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return fmt.Errorf("invalid BlogID format: %w", err)
	}

	filter := bson.M{"_id": objID}

	// Create an update document using $set to update specific fields
	updateDoc := bson.M{
		"$set": bson.M{
			"owner_id":    updatedBlog.OwnerID,
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

func (bldb *BlogDB) GetAllBlogs(pageNumber int) (*[]domain.Blog, error) {
	if pageNumber < 1 {
		pageNumber = 1 // Ensure page number is at least 1
	}
	skip := int64((pageNumber - 1) * pageSize)
	limit := int64(pageSize)

	// Sort by last_update in descending order (latest first)
	findOptions := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "last_update", Value: -1}})

	cursor, err := bldb.Coll.Find(bldb.Contxt, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting all blogs: %w", err)
	}
	defer cursor.Close(bldb.Contxt) // Ensure cursor is closed

	var blogs []domain.Blog
	for cursor.Next(bldb.Contxt) {
		var blog domain.Blog
		var rawBSON bson.M
		if err := cursor.Decode(&rawBSON); err != nil {
			return nil, fmt.Errorf("error decoding raw BSON for blog: %w", err)
		}

		// Manually map _id (primitive.ObjectID) to BlogID (string)
		if id, ok := rawBSON["_id"].(primitive.ObjectID); ok {
			blog.BlogID = id.Hex()
		} else {
			return nil, fmt.Errorf("could not convert _id to primitive.ObjectID")
		}

		// Now, unmarshal the rest of the fields directly.
		// Use bson.Unmarshal for robust conversion, or scan fields one by one
		// if you need more control or have complex types.
		// For simplicity, directly map common fields that don't need _id conversion.
		// You might want to remap the rawBSON to a temp struct or use a custom unmarshaler
		// for Blog if this becomes too cumbersome.

		// For now, let's remap to a temporary struct with primitive.ObjectID and then convert
		// or iterate over the fields in rawBSON if the structure is known.
		// Simpler approach: use a temporary struct for decoding.
		tempBlog := struct {
			ID          primitive.ObjectID `bson:"_id,omitempty"`
			OwnerID     string             `bson:"owner_id"`
			Title       string             `bson:"title"`
			Tags        []string           `bson:"tags"`
			Author      string             `bson:"author"`
			Description string             `bson:"description"`
			LastUpdate  time.Time          `bson:"last_update"`
		}{}

		bsonBytes, err := bson.Marshal(rawBSON)
		if err != nil {
			return nil, fmt.Errorf("error marshaling raw BSON for blog: %w", err)
		}
		if err := bson.Unmarshal(bsonBytes, &tempBlog); err != nil {
			return nil, fmt.Errorf("error unmarshaling blog: %w", err)
		}

		blog.BlogID = tempBlog.ID.Hex() // Ensure BlogID is correctly populated
		blog.OwnerID = tempBlog.OwnerID
		blog.Title = tempBlog.Title
		blog.Tags = tempBlog.Tags
		blog.Author = tempBlog.Author
		blog.Description = tempBlog.Description
		blog.LastUpdate = tempBlog.LastUpdate

		blogs = append(blogs, blog)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through blog cursor: %w", err)
	}
	return &blogs, nil
}
func (bldb *BlogDB) GetAllBlogsByFilter(url_filter domain.Filter, pageNumber int) (*[]domain.Blog, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	skip := int64((pageNumber - 1) * pageSize)
	limit := int64(pageSize)

	filter := bson.M{}

	if url_filter.Tag != "" {
		// Use $regex for case-insensitive partial match on tags array
		filter["tags"] = bson.M{"$regex": primitive.Regex{Pattern: url_filter.Tag, Options: "i"}}
	}

	if !url_filter.AfterDate.IsZero() {
		// Example: blogs created on or after the given date (using LastUpdate)
		filter["last_update"] = bson.M{"$gte": url_filter.AfterDate}
	} 
	// TODO: popularityValue should be calculated from the Popularity dataBase OR SHOULD MERGE THEM IN TO ONE
	// Assuming popularity is now part of the 'Blog' struct, we'll map a hypothetical 'popularity' field.
	// If you want to filter by 'Popularity' from the previous struct, you need to add it to the new Blog struct.
	// For this example, I'll assume `Popularity` is still relevant or you add it back.
	// Since your new struct doesn't have `Popularity`, I'll comment this out or you need to add it to `Blog`.
	// If you *do* add it, uncomment and use a field like `blog.Popularity`.
	// if popularityValue > 0 {
	// 	filter["popularity"] = bson.M{"$gte": popularityValue}
	// }

	findOptions := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "last_update", Value: -1}})

	cursor, err := bldb.Coll.Find(bldb.Contxt, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting filtered blogs: %w", err)
	}
	defer cursor.Close(bldb.Contxt)

	var blogs []domain.Blog
	for cursor.Next(bldb.Contxt) {
		var blog domain.Blog
		var rawBSON bson.M
		if err := cursor.Decode(&rawBSON); err != nil {
			return nil, fmt.Errorf("error decoding raw BSON for filtered blog: %w", err)
		}

		if id, ok := rawBSON["_id"].(primitive.ObjectID); ok {
			blog.BlogID = id.Hex()
		} else {
			return nil, fmt.Errorf("could not convert _id to primitive.ObjectID for filtered blog")
		}

		tempBlog := struct {
			ID          primitive.ObjectID `bson:"_id,omitempty"`
			OwnerID     string             `bson:"owner_id"`
			Title       string             `bson:"title"`
			Tags        []string           `bson:"tags"`
			Author      string             `bson:"author"`
			Description string             `bson:"description"`
			LastUpdate  time.Time          `bson:"last_update"`
		}{}

		bsonBytes, err := bson.Marshal(rawBSON)
		if err != nil {
			return nil, fmt.Errorf("error marshaling raw BSON for filtered blog: %w", err)
		}
		if err := bson.Unmarshal(bsonBytes, &tempBlog); err != nil {
			return nil, fmt.Errorf("error unmarshaling filtered blog: %w", err)
		}

		blog.BlogID = tempBlog.ID.Hex()
		blog.OwnerID = tempBlog.OwnerID
		blog.Title = tempBlog.Title
		blog.Tags = tempBlog.Tags
		blog.Author = tempBlog.Author
		blog.Description = tempBlog.Description
		blog.LastUpdate = tempBlog.LastUpdate

		blogs = append(blogs, blog)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through filtered blog cursor: %w", err)
	}
	return &blogs, nil
}
func (bldb *BlogDB) CheckBlogExistance(blogID string) bool {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return false // Invalid ID format means it cannot exist
	}

	filter := bson.M{"_id": objID}
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
