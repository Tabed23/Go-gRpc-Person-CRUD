package service

import (
	"context"
	"log"
	"persons-gRpc/db"
	"persons-gRpc/pb"
	"persons-gRpc/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PersonServer struct {
	pb.UnimplementedPersonServiceServer
}

var (
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	query  = db.DB
)

func (*PersonServer) CreatePerson(ctx context.Context, in *pb.Person) (*pb.PersonResponse, error) {
	ps := in.GetPhones()
	phones := make([]*types.Person_PhoneNumber, 0)
	for _, p := range ps {
		phones = append(phones, &types.Person_PhoneNumber{
			Number: p.Number,
			Type:   types.Person_PhoneType(p.Type),
		})
	}
	per := types.Person{
		ID:          primitive.NewObjectID(),
		Name:        in.GetName(),
		Email:       in.GetEmail(),
		Phones:      phones,
		CreatedAt:   timestamppb.Now(),
		LastUpdated: timestamppb.Now(),
	}

	_, err := query.Database("PersonGrpc").Collection("persons").InsertOne(ctx, &per)
	if err != nil {
		log.Fatalf("fail to insert data into monog [%v]", err)
	}

	return &pb.PersonResponse{
		Status: "created",
	}, nil

}
func (*PersonServer) GetPersonDetail(ctx context.Context, in *pb.GetPersonDetails) (*pb.Person, error) {
	oid, err := primitive.ObjectIDFromHex(string(in.GetId()))
	if err != nil {
		return nil, err
	}
	pr := types.Person{}
	result := query.Database("PersonGrpc").Collection("persons").FindOne(ctx, bson.M{"_id": oid})
	if err := result.Decode(&pr); err != nil {
		log.Fatalf("cannot decode the person  [%v]", err)
	}
	phones := make([]*pb.Person_PhoneNumber, 0)
	for _, p := range pr.Phones {
		phones = append(phones, &pb.Person_PhoneNumber{
			Number: p.Number,
			Type:   pb.Person_PhoneType(p.Type),
		})
	}
	return &pb.Person{
		Id:          pr.ID.Hex(),
		Name:        pr.Name,
		Email:       pr.Email,
		Phones:      phones,
		CreatedAt:   pr.CreatedAt,
		LastUpdated: pr.LastUpdated,
	}, nil

}
func (*PersonServer) GetAddresBook(ctx context.Context, in *pb.Empty) (*pb.AddressBook, error) {
	proto := make([]*pb.Person, 0)
	phones := make([]*pb.Person_PhoneNumber, 0)
	cursor, err := query.Database("PersonGrpc").Collection("persons").Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("cannot get the data %v", err)
	}
	defer cursor.Close(context.Background())
	
	for cursor.Next(context.Background()) {
		var per types.Person
		if err := cursor.Decode(&per); err != nil {
			log.Fatalf("cannot decode the person %v", err)
		}
		for _, numbers := range per.Phones{
			phones = append(phones, &pb.Person_PhoneNumber{
				Number: numbers.Number,
				Type: pb.Person_PhoneType(numbers.Type),
			})
		}
		proto =append(proto, &pb.Person{
			Name: per.Name,
			Id: per.ID.Hex(),
			Email: per.Email,
			Phones: phones,
			CreatedAt: per.CreatedAt,
			LastUpdated: per.LastUpdated,
		})

	}	
	return &pb.AddressBook{People: proto}, nil

}
func (*PersonServer) DeletePerson(ctx context.Context, in *pb.DeletePerosonID) (*pb.DeletePersonResponse, error) {
	oid, err := primitive.ObjectIDFromHex(string(in.GetId()))
	if err != nil {
		return nil, err
	}
	_, err = query.Database("PersonGrpc").Collection("persons").DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		log.Fatalf("cannot delete the given id {%v}", err)
	}
	return &pb.DeletePersonResponse{
		IsDeleted: true,
		DeletedAt: timestamppb.Now(),
	}, nil

}
func (*PersonServer) UpdatePerson(ctx context.Context, in *pb.Person) (*pb.Person, error) {

	oid, err := primitive.ObjectIDFromHex(string(in.GetId()))
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"name":         in.GetName(),
		"email":        in.GetEmail(),
		"phones":       in.GetPhones(),
		"created_at":   in.GetCreatedAt(),
		"last_updated": timestamppb.Now(),
	}

	result := query.Database("PersonGrpc").Collection("persons").FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))
	resp := types.Person{}
	if err = result.Decode(&resp); err != nil {
		log.Fatalf("cannot update the person %v", err)
	}
	phones := make([]*pb.Person_PhoneNumber, 0)
	for _, p := range resp.Phones {
		phones = append(phones, &pb.Person_PhoneNumber{
			Number: p.Number,
			Type:   pb.Person_PhoneType(p.Type),
		})
	}
	return &pb.Person{
		Id:          resp.ID.Hex(),
		Name:        resp.Name,
		Email:       resp.Email,
		Phones:      phones,
		CreatedAt:   resp.CreatedAt,
		LastUpdated: resp.LastUpdated,
	}, nil
}
