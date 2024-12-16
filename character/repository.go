package character

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"my-gin-app/models"
)

type Repository interface {
	Create(character *models.Character) error
	FindBySlug(slug string) (*models.Character, error)
	UpdateBySlug(slug string, update bson.M) error
	DeleteBySlug(slug string) error
	ListCharacters(skip int64, limit int64) ([]models.Character, error)
	CountCharacters() (int64, error)
}

type MongoRepository struct {
	Collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &MongoRepository{
		Collection: collection,
	}
}

func (r *MongoRepository) Create(character *models.Character) error {
	_, err := r.Collection.InsertOne(context.Background(), character)
	return err
}

func (r *MongoRepository) FindBySlug(slug string) (*models.Character, error) {
	var character models.Character
	err := r.Collection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("character not found")
		}
		return nil, err
	}
	return &character, nil
}

func (r *MongoRepository) UpdateBySlug(slug string, update bson.M) error {
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"slug": slug}, bson.M{"$set": update})
	return err
}

func (r *MongoRepository) DeleteBySlug(slug string) error {
	_, err := r.Collection.DeleteOne(context.Background(), bson.M{"slug": slug})
	return err
}

func (r *MongoRepository) ListCharacters(skip int64, limit int64) ([]models.Character, error) {
	var characters []models.Character
	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetSkip(skip).SetLimit(limit)
	}

	cursor, err := r.Collection.Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var character models.Character
		if err := cursor.Decode(&character); err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	return characters, nil
}

func (r *MongoRepository) CountCharacters() (int64, error) {
	return r.Collection.CountDocuments(context.Background(), bson.D{})
}
