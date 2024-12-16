package models

type TailedBeast struct {
	Name        string   `json:"name" bson:"name"`
	Slug        string   `json:"slug" bson:"slug"`
	Images      []string `json:"images" bson:"images"`
	Rank        string   `json:"rank" bson:"rank"`
	Abilities   []string `json:"abilities" bson:"abilities"`
	Personality string   `json:"personality" bson:"personality"`
}
