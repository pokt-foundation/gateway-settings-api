package controllers

import (
	"context"
	"errors"
	"gateway-settings-api/configs"
	"gateway-settings-api/models"
	"gateway-settings-api/responses"
	"log"
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

	id, _ := primitive.ObjectIDFromHex(application.Id)

	// Does the application exist?
	filter := bson.D{{Key: "_id", Value: id}}

	var result models.Application
	err := applicationsCollection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(responses.ContractAllowlistResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Application not found"}})
		}
		panic(err)
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

	var operations []mongo.WriteModel

	for _, blockchainContractAllowlist := range application.GatewaySettings.ContractsAllowlist {
		// Does the application have this blockchain in the allowlist?
		filter = bson.D{{Key: "_id", Value: id}, {Key: "gatewaySettings.whitelistContracts.blockchain_id", Value: blockchainContractAllowlist.BlockchainID}}

		err = applicationsCollection.FindOne(ctx, filter).Decode(&result)

		if err == mongo.ErrNoDocuments {
			filter = bson.D{{Key: "_id", Value: id}}
			update := bson.M{"$push": bson.M{"gatewaySettings.whitelistContracts": bson.M{"blockchain_id": blockchainContractAllowlist.BlockchainID, "contracts": blockchainContractAllowlist.Contracts}}}

			operation := mongo.NewUpdateOneModel()

			operation.SetFilter(filter)
			operation.SetUpdate(update)

			operations = append(operations, operation)
		} else {
			update := bson.D{{Key: "$push", Value: bson.D{{Key: "gatewaySettings.whitelistContracts.$.contracts", Value: bson.D{{Key: "$each", Value: blockchainContractAllowlist.Contracts}}}}}}

			operation := mongo.NewUpdateOneModel()

			operation.SetFilter(filter)
			operation.SetUpdate(update)

			operations = append(operations, operation)
		}
	}

	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true)

	bulkResult, bulkErr := applicationsCollection.BulkWrite(ctx, operations, &bulkOption)

	if bulkErr != nil {
		log.Printf("%v", bulkErr)
		return c.Status(http.StatusInternalServerError).JSON(responses.ContractAllowlistResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": bulkErr.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.ContractAllowlistResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": bulkResult}})
}
