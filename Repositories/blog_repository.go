package repositories

import (
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
	mainBlogDbName           = "blogbd_test"
	mainBlogDbCollName       = "blogCollection"
	blogPopularityDbName     = "blogPop_test"
	blogPopularityDbCollName = "blogPop_collection"
	pageSize                 = 10
	ASC                      = 1
	DESC                     = -1
)

type BlogDB struct {
	Coll   mongo.Collection
	Contxt context.Context
	Client *mongo.Client
}

func NewBlogDataBaseService() domain.IBlogRepository {
	connection, err := Connect()

	if err != nil {
		log.Fatal("can't initailize ", mainBlogDbName, " Database")
	}

	collection := connection.Client.Database(mainBlogDbName).Collection(mainBlogDbCollName)

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
func (bldb *BlogDB) FindBlogByID(blogID primitive.ObjectID) (*domain.Blog, error) {
	filter := bson.M{"_id": blogID}
	var blog domain.Blog
	err := bldb.Coll.FindOne(bldb.Contxt, filter).Decode(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, nil
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
func (bldb *BlogDB) GetAllBlogsByFilter(url_filter *domain.Filter, pageNumber int) ([]*domain.BlogDTO, error) {
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

	var blogs []*domain.BlogDTO
	for cursor.Next(bldb.Contxt) {
		var blogDTO domain.BlogDTO // Change target type to DTO for decoding

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
	log.Println("Disconnected from Blog MongoDB.")
	return nil
}

// ! BLOG POUPULARITY STARTS HERE, DON'T FORGET TO CREATE IT'S OWN FILE TO MOVE IT THERE IF NECCESSARY
type PopularityDB struct {
	Coll   mongo.Collection
	Contxt context.Context
	Client *mongo.Client
}

func NewBlogPopularityDataBaseService() domain.IPopularityRepository {
	connection, err := Connect()

	if err != nil {
		log.Fatal("can't initailize ", blogPopularityDbName, " Database")
	}

	collection := connection.Client.Database(blogPopularityDbName).Collection(blogPopularityDbCollName)

	return &PopularityDB{
		Coll:   *collection,
		Contxt: context.TODO(),
		Client: connection.Client,
	}

}

func (bldb *PopularityDB) CheckUserLikeBlogID(blogID primitive.ObjectID, userID primitive.ObjectID) bool {
	filterUserInBlog := bson.M{
		"blog_id": blogID,
		"likes":   userID.Hex(),
	}
	count, err := bldb.Coll.CountDocuments(bldb.Contxt, filterUserInBlog)
	if err != nil || count == 0 {
		return false
	}
	return true
}
func (bldb *PopularityDB) CheckUserDisLikeBlogID(blogID primitive.ObjectID, userID primitive.ObjectID) bool {
	filterUserInBlog := bson.M{
		"blog_id":  blogID,
		"dislikes": userID.Hex(),
	}
	count, err := bldb.Coll.CountDocuments(bldb.Contxt, filterUserInBlog)
	if err != nil || count == 0 {
		return false
	}
	return true
}
func (bldb *PopularityDB) UserLikeBlogByID(blogID primitive.ObjectID, userID primitive.ObjectID, revert bool) error {
	// Filter to find the document
	likeFilter := bson.M{"blog_id": blogID}

	// Update statement using the $pull operator
	update := bson.M{}

	if revert {
		update["$pull"] = bson.M{
			"likes": userID.Hex(), // This tells MongoDB to remove the user from the likes array
		}
	} else {
		update["$push"] = bson.M{
			"likes": userID.Hex(), // This tells MongoDB to append the user.Hex() to the likes array
		}
	}
	updateRes, err := bldb.Coll.UpdateOne(bldb.Contxt, likeFilter, update)
	if err != nil {
		return err
	} else if updateRes.ModifiedCount == 0 {
		return fmt.Errorf("no blog found with ID %s to update OR no like were made %t", blogID, revert)
	}
	return nil
}
func (bldb *PopularityDB) UserDisLikeBlogByID(blogID primitive.ObjectID, userID primitive.ObjectID, revert bool) error {
	// Filter to find the document
	dislikeFilter := bson.M{"blog_id": blogID}

	// Update statement using the $pull operator
	update := bson.M{}

	if revert {
		update["$pull"] = bson.M{
			"dislikes": userID.Hex(), // This tells MongoDB to remove the user from the dislikes array
		}
	} else {
		update["$push"] = bson.M{
			"dislikes": userID.Hex(), // This tells MongoDB to append the user.Hex() to the dislikes array
		}
	}
	updateRes, err := bldb.Coll.UpdateOne(bldb.Contxt, dislikeFilter, update)
	if err != nil {
		return err
	} else if updateRes.ModifiedCount == 0 {
		return fmt.Errorf("no blog found with ID %s to update OR no dislike were made %t", blogID, revert)
	}
	return nil
}
func (bldb *PopularityDB) CommentBlogByID(blogID primitive.ObjectID, commentDTO *domain.CommentDTO) error {
	commentFilter := bson.M{
		"blog_id":      blogID,
		"comments._id": commentDTO.UserID,
	}
	updatedComment := bson.M{
		"$set": bson.M{
			"comments.$._id":       commentDTO.UserID,
			"comments.$.user_name": commentDTO.UserName,
			"comments.$.comment":   commentDTO.Comment,
		},
	}
	//? if the comment was found in that blog it will update it other wise it will create a new comment
	opts := options.Update().SetUpsert(true)
	result, err := bldb.Coll.UpdateOne(bldb.Contxt, commentFilter, updatedComment, opts)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("no blog found with ID %s to update OR no comment were made", blogID)
	}
	return nil
}
func (bldb *PopularityDB) CreateBlogPopularity(blogID primitive.ObjectID) error {
	_, err := bldb.Coll.InsertOne(bldb.Contxt, domain.PopularityDTO{BlogID: blogID})
	if err != nil {
		return err
	}
	return nil
}
func (bldb *PopularityDB) IncreaseBlogViewByID(blogID primitive.ObjectID) error {
	filter := bson.M{"blog_id": blogID}

	increaseByOne := bson.M{
		"$inc": bson.M{"view_count": 1},
	}
	updateres, err := bldb.Coll.UpdateOne(bldb.Contxt, filter, increaseByOne)
	if err != nil {
		return err
	} else if updateres.ModifiedCount == 0 {
		return fmt.Errorf("no blog found with ID %s to update OR no view increase were made", blogID)
	}
	return nil
}

func (bldb *PopularityDB) BlogPostLikeCountByID(blogID primitive.ObjectID) (int, error) {
	// Define the aggregation pipeline to match the document and count the 'likes' array.
	pipeline := mongo.Pipeline{
		bson.D{
			{
				Key: "$match",
				Value: bson.D{
					{
						Key:   "blog_id",
						Value: blogID,
					},
				},
			},
		},
		bson.D{
			{
				Key: "$project",
				Value: bson.D{
					{
						Key: "count",
						Value: bson.D{
							{
								Key:   "$size",
								Value: "$likes",
							},
						},
					},
				},
			},
		},
	}

	// Execute the aggregation
	cursor, err := bldb.Coll.Aggregate(bldb.Contxt, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(bldb.Contxt)
	type CountResult struct {
		Count int `bson:"count"`
	}
	// Unmarshal the result
	var results []CountResult
	if err = cursor.All(bldb.Contxt, &results); err != nil {
		return 0, err
	}

	// Return the count, handling the case where no document was found
	if len(results) == 0 {
		return 0, nil
	}
	return results[0].Count, nil
}
func (bldb *PopularityDB) BlogPostDisLikeCountByID(blogID primitive.ObjectID) (int, error) {
	// Define the aggregation pipeline to match the document and count the 'likes' array.
	pipeline := mongo.Pipeline{
		bson.D{
			{
				Key: "$match",
				Value: bson.D{
					{
						Key:   "blog_id",
						Value: blogID,
					},
				},
			},
		},
		bson.D{
			{
				Key: "$project",
				Value: bson.D{
					{
						Key: "count",
						Value: bson.D{
							{
								Key:   "$size",
								Value: "$dislikes",
							},
						},
					},
				},
			},
		},
	}

	// Execute the aggregation
	cursor, err := bldb.Coll.Aggregate(bldb.Contxt, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(bldb.Contxt)
	type CountResult struct {
		Count int `bson:"count"`
	}
	// Unmarshal the result
	var results []CountResult
	if err = cursor.All(bldb.Contxt, &results); err != nil {
		return 0, err
	}

	// Return the count, handling the case where no document was found
	if len(results) == 0 {
		return 0, nil
	}
	return results[0].Count, nil
}
func (bldb *PopularityDB) BlogPostCommentCountByID(blogID primitive.ObjectID) (int, error) {
	// Define the aggregation pipeline to match the document and count the 'likes' array.
	pipeline := mongo.Pipeline{
		bson.D{
			{
				Key: "$match",
				Value: bson.D{
					{
						Key:   "blog_id",
						Value: blogID,
					},
				},
			},
		},
		bson.D{
			{
				Key: "$project",
				Value: bson.D{
					{
						Key: "count",
						Value: bson.D{
							{
								Key:   "$size",
								Value: "$comments",
							},
						},
					},
				},
			},
		},
	}

	// Execute the aggregation
	cursor, err := bldb.Coll.Aggregate(bldb.Contxt, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(bldb.Contxt)
	type CountResult struct {
		Count int `bson:"count"`
	}
	// Unmarshal the result
	var results []CountResult
	if err = cursor.All(bldb.Contxt, &results); err != nil {
		return 0, err
	}

	// Return the count, handling the case where no document was found
	if len(results) == 0 {
		return 0, nil
	}
	return results[0].Count, nil
}
func (bldb *PopularityDB) GetPopularityBlogByID(blogID primitive.ObjectID) (*domain.PopularityDTO, error) {
	filter := bson.M{"blog_id": blogID}
	var popBlog domain.PopularityDTO
	err := bldb.Coll.FindOne(bldb.Contxt, filter).Decode(&popBlog)
	if err != nil {
		return nil, err
	}
	return &popBlog, nil
}
func (bldb *PopularityDB) GetPopularityByFilter(order int, pageNumber int) ([]*domain.PopularityDTO, error) {
	skip := int64((pageNumber - 1) * pageSize)
	limit := int64(pageSize)
	if !(order == ASC || order == DESC) {
		order = DESC
	}
	filter := bson.M{}
	findOptions := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "popularity_value", Value: order}})

	cursor, err := bldb.Coll.Find(bldb.Contxt, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting filtered Popularity: %w", err)
	}
	defer cursor.Close(bldb.Contxt)

	var blogs []*domain.PopularityDTO
	for cursor.Next(bldb.Contxt) {
		var blogDTO domain.PopularityDTO // Change target type to DTO for decoding

		if err := cursor.Decode(&blogDTO); err != nil {
			return nil, fmt.Errorf("error decoding filtered Popularity DTO: %w", err)
		}

		blogs = append(blogs, &blogDTO)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through filtered Popularity cursor: %w", err)
	}
	return blogs, nil
}

func (bldb *PopularityDB) CloseDataBase() error {
	if bldb.Client == nil {
		return nil // Nothing to close
	}
	if err := bldb.Client.Disconnect(bldb.Contxt); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %w", err)
	}
	log.Println("Disconnected from Popularity MongoDB.")
	return nil
}
