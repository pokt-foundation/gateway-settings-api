package controllers

import (
	"context"
	"errors"
	"gateway-settings-api/configs"
	"gateway-settings-api/models"
	"gateway-settings-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var applicationsCollection *mongo.Collection = configs.GetCollection(configs.DB, configs.EnvAppsCollectionName())
var validate = validator.New()

func AddContractToAllowlist(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var application models.Application
	defer cancel()

	if err := c.BodyParser(&application); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ContractAllowlistResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&application); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.ContractAllowlistResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	applicationIDs, ok := claims["applicationIDs"].([]interface{})

	if !ok {
		return errors.New("applicationIDs is not an array")
	}

	var allowed bool

	for _, appId := range applicationIDs {
		if appId.(string) == application.Id {
			allowed = true
		}
	}

	if !allowed {
		return c.Status(http.StatusUnauthorized).JSON(responses.ContractAllowlistResponse{Status: http.StatusUnauthorized, Message: "This application doesn't belong to your user.", Data: nil})
	}

	id, _ := primitive.ObjectIDFromHex(application.Id)

	var operations []mongo.WriteModel

	for _, blockchainContractAllowlist := range application.GatewaySettings.ContractsAllowlist {
		filter := bson.D{{Key: "_id", Value: id}, {Key: "gatewaySettings.whitelistContracts.blockchain_id", Value: blockchainContractAllowlist.BlockchainID}}
		update := bson.D{{Key: "$push", Value: bson.D{{Key: "gatewaySettings.whitelistContracts.$.contracts", Value: bson.D{{Key: "$each", Value: blockchainContractAllowlist.Contracts}}}}}}

		operation := mongo.NewUpdateOneModel()

		operation.SetFilter(filter)
		operation.SetUpdate(update)
		operation.SetUpsert(true)

		operations = append(operations, operation)
	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)

	result, err := applicationsCollection.BulkWrite(ctx, operations, &bulkOption)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.ContractAllowlistResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.ContractAllowlistResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}
