package infrastructure

import (
	domain "blog_starter_project_g66/Domain"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
)

type AI_Interaction struct {
	Client                         *genai.Client
	Model                          *genai.GenerativeModel
	CountInteraction               int
	UserFirstInteractionTimestamps time.Time
	UserLastInteractionTimestamps  time.Time
	Ctx                            context.Context
}

func NewAICommentInteraction(clientID primitive.ObjectID) domain.IAIInteraction {
	ai_int := CreateNewAIInteraction(clientID)
	return &AICommentInteraction{
		AI_Interaction: ai_int,
	}
}
func NewAIBlogInteraction(clientID primitive.ObjectID) domain.IAIInteraction {
	ai_int := CreateNewAIInteraction(clientID)
	return &AIBlogInteraction{
		AI_Interaction: ai_int,
	}
}
func NewAIBlogFilterInteraction(clientID primitive.ObjectID) domain.IAIInteraction {
	ai_int := CreateNewAIInteraction(clientID)
	return &AIBlogFilterInteraction{
		AI_Interaction: ai_int,
	}
}

type AICommentInteraction struct {
	*AI_Interaction
}
type AIBlogInteraction struct {
	*AI_Interaction
}
type AIBlogFilterInteraction struct { //? UseCase: if the user doesn't know anyting about what might be the filter words but he/she knows what they want so the ai will help them by giving a filter object. it will be used to filter blogs based on that thing
	*AI_Interaction
}

func CreateNewAIInteraction(clientID primitive.ObjectID) *AI_Interaction {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		//! REMEMBER to switch back to environment variables for production.
		apiKey = "AIzaSyAsU5HGko_VqY1ZBcEak5UYK_I0c1W2dL4"
	}
	// if the user interacts today it's countInteraction and
	// the timestamps should be loaded from the database
	// so that the user ONLY get limited amout of prompt per day.
	// THIS FEATURE COULD BE USED FOR USERS WHEN THEY BECOME PAID USER

	// Create the Gemini client.
	UserFirstInteractionTimestamps := time.Now()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil
	}
	model := client.GenerativeModel("gemini-2.5-flash")
	return &AI_Interaction{
		Client:                         client,
		Ctx:                            ctx,
		Model:                          model,
		UserFirstInteractionTimestamps: UserFirstInteractionTimestamps,
		UserLastInteractionTimestamps:  UserFirstInteractionTimestamps,
		CountInteraction:               0,
	}
}

func (aint *AI_Interaction) IncrementInteractionCount() {
	aint.CountInteraction += 1
}
func (aint *AI_Interaction) CloseAIConnection() error {
	return aint.Client.Close()
}
func (aint *AI_Interaction) IsClientConnected() bool {
	// This is for the time being; the connection should be checked
	// from the other child AI interactions to use theirs connection.
	return aint.Client != nil
}
func (aint *AI_Interaction) GenerateContent(prompt string) (*domain.AIResponse, error) {
	aint.IncrementInteractionCount() // increament the interaction for observation
	resp, err := aint.Model.GenerateContent(aint.Ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		respText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

		var aiResponse domain.AIResponse
		if err := json.Unmarshal([]byte(respText), &aiResponse); err != nil {

			return nil, err
		}
		return &aiResponse, nil
	}
	return nil, fmt.Errorf("no content generated")
}

func (commInt *AICommentInteraction) ParseJsonBodyToDomain(aiResponse *domain.AIResponse) any {
	return aiResponse.MainResponse
}
func (commInt *AICommentInteraction) CallAIAndGetResponse(developerMessage string, userMessage string, jsonBodyStirng string) (*domain.AIResponse, error) {
	overAllPrompt := `
	Paraphrase the following without changing the context or sentiment. Correct grammar, sentence structure, and casing.
	Generate a plain string response with no special formatting or code block delimiters, in this JSON structure:
	{
		"main_response": // the paraphrased string
		"editorial_response": // the editorial message from the AI
		"is_nil_response": // a boolean indicating if main_response is null
	}
	Here is the paragraph:
` + userMessage

	return commInt.GenerateContent(overAllPrompt)
}

func (blgInt *AIBlogInteraction) ParseJsonBodyToDomain(aiResponse *domain.AIResponse) any {
	var generatedBlog domain.Blog

	if err := json.Unmarshal(aiResponse.MainResponse, &generatedBlog); err != nil {

		return nil
	}

	return &generatedBlog
}
func (blgInt *AIBlogInteraction) CallAIAndGetResponse(developerMessage string, userMessage string, jsonBodyStirng string) (*domain.AIResponse, error) {
	overAllPrompt := userMessage + "Today is: " + time.Now().Format("2006-01-02") + `
	.Generate a response in the following JSON format as a plain string, with no special formatting:
	{
		"main_response": {
			"title": // the blog's title
			"tags": // a list of up to 7 possible tags for the blog
			"author": "<Your_Name>" // a placeholder for the author's name
			"description": // the body of the blog
			"last_update_time": // the current timestamp
		}
		"editorial_response": // an editorial message from the AI
		"is_nil_response": // a boolean indicating if the main_response is null
	}
	Use the content from the following JSON to populate the fields:
` + jsonBodyStirng
	return blgInt.GenerateContent(overAllPrompt)
}
func (blgFilter *AIBlogFilterInteraction) ParseJsonBodyToDomain(aiResponse *domain.AIResponse) any {
	var generatedFilter domain.AIBlogFilter
	if err := json.Unmarshal(aiResponse.MainResponse, &generatedFilter); err != nil {

		return nil
	}
	return &generatedFilter
}
func (blgFilter *AIBlogFilterInteraction) CallAIAndGetResponse(developerMessage string, userMessage string, jsonBodyStirng string) (*domain.AIResponse, error) {
	overAllPrompt := developerMessage + jsonBodyStirng + "Today is: " + time.Now().Format("2006-01-02") + `
	Create a filter based on the following user prompt: 
` + userMessage + `
	Generate a plain string response with no special formatting, using this JSON structure:
	{
		"main_response": {
			"tags": // a list of tags from the user's prompt
			"after_date": // an ISO 8601 formatted date. If the user mentions "last week," calculate it from last Monday. If "last month," use the first day of the previous month. If "last year," use the first day of the previous year. For all other relative dates (e.g., "yesterday"), provide a date representing the start of the previous full day.
			"title": // the blog's title from the user prompt
			"author_name": // the blog author's name
		}
		"editorial_response": // an editorial message from the AI
		"is_nil_response": // a boolean value indicating if the main_response is null; make it false if any field is not null
	}
	Use an empty string or null for any field not mentioned in the user prompt.
`
	return blgFilter.GenerateContent(overAllPrompt)
}
