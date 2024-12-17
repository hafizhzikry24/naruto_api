package models

type Personal struct {
	Birthdate   string `json:"birthdate" bson:"birthdate"`
	Sex         string `json:"sex" bson:"sex"`
	Status      string `json:"status" bson:"status"`
	Height      string `json:"height" bson:"height"`
	Weight      string `json:"weight" bson:"weight"`
	BloodType   string `json:"bloodType" bson:"bloodType"`
	Occupation  string `json:"occupation" bson:"occupation"`
	Affiliation string `json:"affiliation" bson:"affiliation"`
	Clan        string `json:"clan" bson:"clan"`
}

type Rank struct {
	NinjaRank string `json:"ninjaRank" bson:"ninjaRank"`
}

type Debut struct {
	Anime     string `json:"anime" bson:"anime"`
	AppearsIn string `json:"appearsIn" bson:"appearsIn"`
}

type Character struct {
	Name     string   `json:"name" bson:"name"`
	Slug     string   `json:"slug" bson:"slug"`
	Images   []string `json:"images" bson:"images"`
	Personal Personal `json:"personal" bson:"personal"`
	Rank     Rank     `json:"rank" bson:"rank"`
	Debut    Debut    `json:"debut" bson:"debut"`
	Jutsu    []string `json:"jutsu" bson:"jutsu"`
}
