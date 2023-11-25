package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PosterMetaCollectionName string = "poster_meta"
const PostsMetaCollectionName string = "posts_meta"

type MongoClientCtx struct {
	mongoClient *mongo.Client
	dbName      string

	ctx context.Context
}

func InitMongo(URI string, dbName string) (*MongoClientCtx, error) {
	// 设置客户端选项
	clientOptions := options.Client().ApplyURI(URI)

	// 连接到 MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoClientCtx{
		mongoClient: client,
		dbName:      dbName,
		ctx:         ctx,
	}, nil
}

func (c *MongoClientCtx) InsertPosterPosts(meta *DBPosterMeta, postMetas []*DBPostMeta) error {
	posterCollection := c.mongoClient.Database(c.dbName).Collection(PosterMetaCollectionName)
	postsCollection := c.mongoClient.Database(c.dbName).Collection(PostsMetaCollectionName)

	insertOneResult, err := posterCollection.InsertOne(c.ctx, meta)
	if err != nil {
		return err
	}

	fmt.Println("Inserted a single document: ", insertOneResult.InsertedID)

	var interfaceSlice []interface{} = make([]interface{}, len(postMetas))

	for i, v := range postMetas {
		interfaceSlice[i] = v
	}

	var insertManyResult *mongo.InsertManyResult
	insertManyResult, err = postsCollection.InsertMany(c.ctx, interfaceSlice)
	if err != nil {
		return err
	}

	fmt.Println("Inserted many documents: ", insertManyResult.InsertedIDs)

	return nil
}
