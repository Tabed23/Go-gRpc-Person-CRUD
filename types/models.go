package types

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type Person struct {
	Name                 string                  `bson:"name" json:"name"`
	ID 					 primitive.ObjectID      `bson:"_id" json:"_id,omitempty"`
	Email                string                  `bson:"email" json:"email"`
	Phones               []*Person_PhoneNumber   `bson:"phones" json:"phones"`
	CreatedAt            *timestamppb.Timestamp  `bson:"created_at" json:"created_at"`
	LastUpdated          *timestamppb.Timestamp  `bson:"last_updated" json:"last_updated"`
}
type Person_PhoneNumber struct {
	Number               string           `bson:"number" json:"number"`
	Type                 Person_PhoneType `bson:"type" json:"type"`
}

type AddressBook struct {
	People               []*Person `bson:"people" json:"people"`
}
type Person_PhoneType int 
const (
	number Person_PhoneType = iota  + 1 
	MOBILE 
    HOME 
    WORK 
)