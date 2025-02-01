package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mymmrac/telego"
	"github.com/wneessen/go-mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	tu "github.com/mymmrac/telego/telegoutil"
	"pureheroky.com/backend/database"
	"pureheroky.com/backend/models"
	"pureheroky.com/backend/utils"
)

type ControllerImpl struct {
	mailClient *mail.Client
	bot        *telego.Bot
}

func ControllerService(mailClient *mail.Client, bot *telego.Bot) *ControllerImpl {
	return &ControllerImpl{
		mailClient: mailClient,
		bot:        bot,
	}
}

func (c *ControllerImpl) GetUserData(ctx *fiber.Ctx) error {
	collection := database.MongoDB.Collection("users")

	var userDoc models.UserData
	err := collection.FindOne(ctx.Context(), bson.M{}).Decode(&userDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.Status(fiber.StatusNotFound).JSON(models.UserResponse{
				Data:   models.UserData{},
				Status: 404,
			})
		}
		log.Printf("Error finding user data: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user data",
		})
	}

	return ctx.JSON(models.UserResponse{
		Data:   userDoc,
		Status: 200,
	})
}

func (c *ControllerImpl) SetUserData(ctx *fiber.Ctx) error {
	var newUserData models.UserData
	if err := ctx.BodyParser(&newUserData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	collection := database.MongoDB.Collection("users")
	filter := bson.M{}
	update := bson.M{"$set": newUserData}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx.Context(), filter, update, opts)
	if err != nil {
		log.Printf("Error upserting user data: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upsert user data",
		})
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}

func (c *ControllerImpl) GetSkills(ctx *fiber.Ctx) error {
	collection := database.MongoDB.Collection("skills")

	var skillsDoc struct {
		ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		Skills []string           `bson:"skills"     json:"skills"`
	}

	err := collection.FindOne(ctx.Context(), bson.M{}).Decode(&skillsDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.JSON(models.SkillsResponse{
				Skills: []string{},
				Status: 200,
			})
		}
		log.Printf("Error finding skills: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch skills",
		})
	}

	return ctx.JSON(models.SkillsResponse{
		Skills: skillsDoc.Skills,
		Status: 200,
	})
}

func (c *ControllerImpl) AddSkill(ctx *fiber.Ctx) error {
	var request struct {
		Skill string `json:"skill"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if request.Skill == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Skill cannot be empty",
		})
	}

	collection := database.MongoDB.Collection("skills")

	filter := bson.M{}
	update := bson.M{"$push": bson.M{"skills": request.Skill}}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx.Context(), filter, update, opts)
	if err != nil {
		log.Printf("Error adding skill: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add skill",
		})
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}

func (c *ControllerImpl) GetProjects(ctx *fiber.Ctx) error {
	collection := database.MongoDB.Collection("projects")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error finding projects: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch projects",
		})
	}
	defer cursor.Close(context.Background())

	var projects []models.Project
	if err = cursor.All(context.Background(), &projects); err != nil {
		log.Printf("Error decoding projects: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode projects",
		})
	}

	return ctx.JSON(models.ProjectResponse{
		Projects: projects,
		Status:   200,
	})
}

func (c *ControllerImpl) CreateProject(ctx *fiber.Ctx) error {
	var project models.Project
	if err := ctx.BodyParser(&project); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	collection := database.MongoDB.Collection("projects")

	res, err := collection.InsertOne(ctx.Context(), project)
	if err != nil {
		log.Printf("Error inserting project: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to insert project",
		})
	}

	return ctx.JSON(fiber.Map{
		"insertedId": res.InsertedID,
		"status":     "ok",
	})
}

func (c *ControllerImpl) GetLatestCommits(ctx *fiber.Ctx) error {
	token := os.Getenv("GIT_TOKEN")

	if token == "" {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "GitHub token not found in environment GITHUB_TOKEN",
		})
	}

	repos, err := utils.FetchUserRepos(token)
	if err != nil {
		log.Printf("Error fetching user repos: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user repos",
		})
	}

	var allCommits []models.CommitInfo
	for _, repo := range repos {
		commits, err := utils.FetchCommitsForRepo(token, repo.Name)
		if err != nil {
			log.Printf("Error fetching commits for repo %s, %v\n", repo.Name, err)
			continue
		}
		for _, commit := range commits {
			allCommits = append(allCommits, models.CommitInfo{
				ProjectName: repo.Name,
				Branch:      repo.DefaultBranch,
				Date:        commit.Commit.Author.Date,
				Message:     commit.Commit.Message,
			})
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  200,
		"commits": allCommits,
	})
}

func (c *ControllerImpl) SendRequest(ctx *fiber.Ctx) error {
	var request models.UserRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	messageText := fmt.Sprintf("New request from %s - <code>%s</code>\n\n%s", request.UserName, request.UserEmail, request.UserMessage)
	chatIDStr := os.Getenv("CHAT_ID")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Println("Ошибка преобразования CHAT_ID:", err)
		return err
	}

	message := tu.Message(tu.ID(chatID), messageText)
	message.ParseMode = telego.ModeHTML

	_, err = c.bot.SendMessage(message)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "There was error trying to send a telegram request: " + err.Error(),
		})
	}

	userAnswer := fmt.Sprintf("Your message was: \n\n%s\n\nI will contact you shortly", request.UserMessage)
	newMessage := mail.NewMsg()
	newMessage.From(os.Getenv("EMAIL_FROM"))
	newMessage.To(request.UserEmail)
	newMessage.Subject("Thanks for your request! - pureheroky.com")
	newMessage.SetBodyString(mail.TypeTextPlain, userAnswer)

	if err := c.mailClient.DialAndSend(newMessage); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "There was error trying to send an email: " + err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status": 200,
	})
}
