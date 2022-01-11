package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Jobdata struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Comapny    string             `json:"company,omitempty"`
	From       string             `json:"from,omitempty"`
	To         string             `json:"To,omitempty"`
	Experience string             `json:"experience,omitempty"`
}

type UserInfo struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Address    string `json:"address"`
	Loginid    string `json:"loginid"`
	Password   string `json:"password"`
	Jobdetails Jobdata
}
