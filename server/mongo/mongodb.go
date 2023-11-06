// mongogdb.go

package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var (
	MongoDBURL string = "mongodb://localhost:27017/frontend-monitor"
)

type Mongodb struct {
	url      string
	database string
	client   *mongo.Client
}

func NewMongoDB(url string) (*Mongodb, error) {
	cs, err := connstring.Parse(url)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(context.Background(), mongoOptions.Client().ApplyURI(url).SetConnectTimeout(10*time.Second))
	if err != nil {
		return nil, err
	}

	db := &Mongodb{
		url:      url,
		database: cs.Database,
		client:   client,
	}

	return db, nil
}

func (db *Mongodb) C(collection string) *mongo.Collection {
	return db.client.Database(db.database).Collection(collection)
}

func (db *Mongodb) Ping(ctx context.Context) error {
	return db.client.Ping(ctx, readpref.Primary())
}

func (db *Mongodb) InsertOne(ctx context.Context, collection string, doc interface{}) (id interface{}, err error) {
	var result *mongo.InsertOneResult

	result, err = db.C(collection).InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (db *Mongodb) QueryAll(ctx context.Context, collection string, filter interface{}, skip int, limit int, result interface{}) error {
	if filter == nil {
		filter = bson.M{}
	}

	opts := mongoOptions.Find()

	opts.SetSort(bson.M{"id": -1}) // 按id倒序
	opts.SetSkip(int64(skip))      // 跳过skip条
	if limit > 0 {
		opts.SetLimit(int64(limit)) // 限制limit条
	}

	cursor, err := db.C(collection).Find(ctx, filter, opts)
	if err != nil {
		return nil
	}
	return cursor.All(ctx, result)
}

func (db *Mongodb) UpdateOne(ctx context.Context, collection string, id int64, doc interface{}) error {
	// 根据id更新 doc 所有字段
	_, err := db.C(collection).UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": doc})
	return err
}