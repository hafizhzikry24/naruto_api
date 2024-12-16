package tailedbeast

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"my-gin-app/models"
)

type Repository interface {
	Create(beast *models.TailedBeast) error
	FindBySlug(slug string) (*models.TailedBeast, error)
	UpdateBySlug(slug string, update bson.M) error
	DeleteBySlug(slug string) error
	ListBeasts(skip int64, limit int64) ([]models.TailedBeast, error)
	CountBeasts() (int64, error)
}

type MongoRepository struct {
	Collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) Repository {
	return &MongoRepository{
		Collection: collection,
	}
}

func (r *MongoRepository) Create(beast *models.TailedBeast) error {
	_, err := r.Collection.InsertOne(context.Background(), beast)
	return err
}

func (r *MongoRepository) FindBySlug(slug string) (*models.TailedBeast, error) {
	var beast models.TailedBeast
	err := r.Collection.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&beast)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("tailed beast not found")
		}
		return nil, err
	}
	return &beast, nil
}

func (r *MongoRepository) UpdateBySlug(slug string, update bson.M) error {
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"slug": slug}, bson.M{"$set": update})
	return err
}

func (r *MongoRepository) DeleteBySlug(slug string) error {
	_, err := r.Collection.DeleteOne(context.Background(), bson.M{"slug": slug})
	return err
}

func (r *MongoRepository) ListBeasts(skip int64, limit int64) ([]models.TailedBeast, error) {
	var beasts []models.TailedBeast
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
		var beast models.TailedBeast
		if err := cursor.Decode(&beast); err != nil {
			return nil, err
		}
		beasts = append(beasts, beast)
	}

	return beasts, nil
}

func (r *MongoRepository) CountBeasts() (int64, error) {
	return r.Collection.CountDocuments(context.Background(), bson.D{})
}
