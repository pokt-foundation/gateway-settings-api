package controllers

import (
	"context"
	"gateway-settings-api/configs"
	"gateway-settings-api/models"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const SALT_ROUNDS = 10

var usersCollection *mongo.Collection = configs.GetCollection(configs.DB, configs.EnvUsersCollectionName())
var loadBalancersCollection *mongo.Collection = configs.GetCollection(configs.DB, configs.EnvLoadBalancersCollectionName())

// Login get user and password
func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	email := input.Email
	password := []byte(input.Password)

	var userResult models.User

	userFilter := bson.D{{Key: "email", Value: email}}
	err := usersCollection.FindOne(ctx, userFilter).Decode(&userResult)

	if err != nil {
		return err
	}

	dbHashedPassword := []byte(userResult.Password)

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(dbHashedPassword, password)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var lbResults []models.LoadBalancer

	lbFilter := bson.D{{Key: "user", Value: userResult.Id}}
	cursor, err := loadBalancersCollection.Find(ctx, lbFilter)

	if err = cursor.All(ctx, &lbResults); err != nil {
		log.Fatal(err)
	}

	var applicationIDs []string

	for _, lb := range lbResults {
		appIds := lb.ApplicationIDs
		applicationIDs = append(applicationIDs, appIds...)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["applicationIDs"] = applicationIDs
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(configs.EnvJWTSigningKey()))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}
