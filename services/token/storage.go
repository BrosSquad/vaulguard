package token

import (
	"context"

	"github.com/BrosSquad/vaulguard/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Storage interface {
	Get(interface{}) (models.TokenDto, error)
	Create(token *models.TokenDto) (*models.TokenDto, error)
}

func NewSqlStorage(db *gorm.DB) Storage {
	return sqlStorage{
		db:    db,
		cache: make(map[uint]models.Token),
	}
}

func NewMongoStorage(ctx context.Context, client *mongo.Collection) Storage {
	return mongoStorage{
		ctx:    ctx,
		client: client,
		cache:  make(map[primitive.ObjectID]models.Token),
	}
}

type sqlStorage struct {
	db    *gorm.DB
	cache map[uint]models.Token
}

func (s sqlStorage) Get(idOrObjectId interface{}) (models.TokenDto, error) {
	id := uint(idOrObjectId.(uint64))
	var token models.Token
	if _, ok := s.cache[id]; !ok {
		tx := s.db.Joins("Application").First(&token, id)
		if err := tx.Error; err != nil {
			return models.TokenDto{}, err
		}
		s.cache[id] = token
	} else {
		token = s.cache[id]
	}
	return models.TokenDto{
		ID:            token.ID,
		Value:         token.Value,
		ApplicationId: token.ApplicationId,
		CreatedAt:     token.CreatedAt,
		UpdatedAt:     token.UpdatedAt,
		Application: models.ApplicationDto{
			ID:        token.Application.ID,
			Name:      token.Application.Name,
			CreatedAt: token.Application.CreatedAt,
			UpdatedAt: token.Application.UpdatedAt,
		},
	}, nil
}

func (s sqlStorage) Create(tokenDto *models.TokenDto) (*models.TokenDto, error) {
	token := models.Token{
		Value:         tokenDto.Value,
		ApplicationId: tokenDto.ApplicationId.(uint),
		CreatedAt:     tokenDto.CreatedAt,
		UpdatedAt:     tokenDto.UpdatedAt,
	}
	if err := s.db.Create(&token).Error; err != nil {
		return nil, err
	}

	tokenDto.ID = token.ID

	return tokenDto, nil
}

type mongoStorage struct {
	ctx    context.Context
	client *mongo.Collection
	cache  map[primitive.ObjectID]models.Token
}

func (m mongoStorage) Get(idOrObjectID interface{}) (models.TokenDto, error) {
	objectID := idOrObjectID.(primitive.ObjectID)
	filter := bson.A{
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "applications"},
				{Key: "localField", Value: "ApplicationId"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "app"},
			}},
		},
		bson.D{
			{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}}},
		},
		bson.D{
			{Key: "$unwind", Value: "$app"},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "ID", Value: "$_id"},
				{Key: "Value", Value: "$Value"},
				{Key: "ApplicationId", Value: "$ApplicationId"},
				{Key: "Application", Value: bson.D{
					{Key: "Id", Value: "$app._id"},
					{Key: "Name", Value: "$app.Name"},
					{Key: "CreatedAt", Value: "$app.CreatedAt"},
					{Key: "UpdatedAt", Value: "$app.UpdatedAt"},
				}},
				{Key: "CreatedAt", Value: "$CreatedAt"},
				{Key: "UpdatedAt", Value: "$UpdatedAt"},
			}},
		},
	}

	cursor, err := m.client.Aggregate(m.ctx, filter)

	if err != nil {
		return models.TokenDto{}, err
	}

	defer cursor.Close(m.ctx)
	token := models.TokenDto{}

	for cursor.Next(m.ctx) {
		err := cursor.Decode(&token)
		if err != nil {
			return models.TokenDto{}, err
		}
	}

	return token, nil
}

func (m mongoStorage) Create(token *models.TokenDto) (*models.TokenDto, error) {
	inserted, err := m.client.InsertOne(m.ctx, bson.M{
		"Value":         token.Value,
		"ApplicationId": token.ApplicationId.(primitive.ObjectID),
		"CreatedAt":     token.CreatedAt,
		"UpdatedAt":     token.UpdatedAt,
	})

	if err != nil {
		return nil, err
	}

	token.ID = inserted.InsertedID

	return token, nil
}
