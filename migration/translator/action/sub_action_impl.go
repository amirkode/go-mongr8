package action

// import (
// 	"context"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // init subaction in MongoDB manner

// func NewCreateCollectionSubAction(ctx context.Context, db *mongo.Database, collectionName string, option *bson.M) SubAction {
// 	execFuc := func() {
// 		opts := options.CreateCollection()
// 		db.CreateCollection(ctx, collectionName, bson.M{})
// 	}
// }