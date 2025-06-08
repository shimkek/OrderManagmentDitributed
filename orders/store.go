package main

import (
	"context"

	"github.com/shimkek/omd-common/api"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	DBName   = "orders"
	CollName = "orders"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{db}
}

func (s *store) Create(ctx context.Context, o Order) (bson.ObjectID, error) {
	col := s.db.Database(DBName).Collection(CollName)

	newOrder, err := col.InsertOne(ctx, o)
	return newOrder.InsertedID.(bson.ObjectID), err
}

func (s *store) Update(ctx context.Context, orderID string, p *api.Order) error {
	col := s.db.Database(DBName).Collection(CollName)

	updateFields := bson.M{}
	if p.PaymentLink != "" {
		updateFields["paymentLink"] = p.PaymentLink
	}
	if p.Status != "" {
		updateFields["status"] = p.Status
	}
	ObjectID, err := bson.ObjectIDFromHex(orderID)
	if err != nil {
		return err
	}
	_, err = col.UpdateOne(ctx,
		bson.M{"_id": ObjectID},
		bson.M{"$set": updateFields},
	)

	// for i, order := range orders {
	// 	if order.OrderID == orderID {
	// 		if p.Status != "" {
	// 			orders[i].Status = p.Status
	// 		}
	// 		if p.PaymentLink != "" {
	// 			orders[i].PaymentLink = p.PaymentLink
	// 		}
	// 	}
	// }
	return err
}
func (s *store) Get(ctx context.Context, orderID string) (*Order, error) {
	col := s.db.Database(DBName).Collection(CollName)

	var o Order
	hexID, err := bson.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, err
	}
	err = col.FindOne(ctx, bson.M{
		"_id": hexID,
	}).Decode(&o)

	return &o, err
}
func (s *store) Delete(ctx context.Context) error {
	return nil
}
