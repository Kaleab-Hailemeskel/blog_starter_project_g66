package infrastructure

import (
	domain "blog_starter_project_g66/Domain"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	log.Println("😜while gemini working on it")
	log.Println(resp)
	log.Println(err)
	log.Println("-------------")
	if err != nil {
		return nil, err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		respText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
		log.Println("✔ ", respText)
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
  overAllPrompt := userMessage + `
pharaphrase the following without leaving the context, and feeling. ` + developerMessage + `
correct: the grammar sentense structure and casing
Generate a response that is a plain string. Do not use any special formatting or code block delimiters.
{
  "main_response": // this contain the strring after pahraphrasing
  "editorial_response": // this should contain the editorial message from the ai response
  "is_nil_response": // this should contain a boolean value weather the main_response is null or not
}
   Here is the pharagraph:
` + jsonBodyStirng

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
  overAllPrompt := developerMessage + userMessage + `
        .The json format for the blog looks like the following:
    ` + jsonBodyStirng + `

        I want you to send me a response ONLY in the type of the following json string stucture, Generate a response that is a plain string. Do not use any special formatting or code block delimiters:
    {
            "main_response":    {
                    "title": // the title of the blog
                    "tags": // list of possible tags for the blog, not more than 7,
                    "author": // create a place holder here like <Your_Name>,
                    "description": // the description part should be the body of the blog,
                    "last_update_time": // the current time stamp.
            }
            "editorial_response": // this should contain the editorial message from the ai response
            "is_nil_response": // this should contain a boolean value weather the main_response is null or not
    }
    `
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
    Create a filter using the information provided by the following prompt:
` + userMessage + `
    .Generate a response that is a plain string. Do not use any special formatting or code block delimiters. The json filter structure looks like this:
    {
    "main_response": {
            "tags": // list of possible tags from the user prompt
            "after_date": // provide the ISO 8601 formatted date representing the start of the last full day. For example, if today is August 6, 2025, the value should be "2025-08-05T00:00:00Z". This ensures the filter includes posts from the previous day but not the current one. if the user has some thing about last week: calculate it from last week monday
            last month: calculate it from the first day of last month
            last year calculate it from the first day of the last year
                "title": // if the user gives a title for the blog
            "author_name": // possible name of the blog author
                }  
                "editorial_response": // this should contain the editorial message from the ai response
                "is_nil_response": // this should contain a boolean value whether the main_response is null or not; if there is at least one non-null value, make it false
            }
    Keep in mind that if the prompt didn't mention any of the fields, use an empty string or null value as needed when returning the result.`

  return blgFilter.GenerateContent(overAllPrompt)
}
