package controllers

import (
	"context"
	"gateway-settings-api/configs"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const SALT_ROUNDS = 10

var usersCollection *mongo.Collection = configs.GetCollection(configs.DB, "TestUsers")
var loadBalancersCollection *mongo.Collection = configs.GetCollection(configs.DB, "TestLoadBalancers")

// Login get user and password
func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type LoginInput struct {
		Email    string `json:"identity"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	email := input.Email
	password := []byte(input.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, SALT_ROUNDS)
	if err != nil {
		panic(err)
	}

	var userResult bson.M

	userFilter := bson.D{{Key: "email", Value: email}}
	err = usersCollection.FindOne(ctx, userFilter).Decode(&userResult)

	if err != nil {
		return err
	}

	dbHashedPassword := []byte(userResult["password"].(string))

	// Comparing the password with the hash
	err = bcrypt.CompareHashAndPassword(hashedPassword, dbHashedPassword)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var lbResult bson.M

	id, _ := primitive.ObjectIDFromHex(userResult["_id"].(string))

	lbFilter := bson.D{{Key: "_id", Value: id}}
	err = usersCollection.FindOne(ctx, lbFilter).Decode(&lbResult)

	if err != nil {
		return err
	}

	applicationIDs := lbResult["applicationIDs"].([]string)

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
