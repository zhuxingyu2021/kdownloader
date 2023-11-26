package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"kdownloader/pkg/utils"
)

const PosterMetaCollectionName string = "poster_meta"
const PostsMetaCollectionName string = "posts_meta"

type MongoClientCtx struct {
	mongoClient *mongo.Client
	dbName      string

	ctx context.Context
}

func InitMongo(ctx context.Context, URI string, dbName string) (*MongoClientCtx, error) {
	// 设置客户端选项
	clientOptions := options.Client().ApplyURI(URI)

	// 连接到 MongoDB
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

func (c *MongoClientCtx) Close() error {
	return c.mongoClient.Disconnect(c.ctx)
}

func (c *MongoClientCtx) insertPostMetas(postMetas []*DBPostMeta) error {
	postsCollection := c.mongoClient.Database(c.dbName).Collection(PostsMetaCollectionName)
	var interfaceSlice []interface{} = make([]interface{}, len(postMetas))

	for i, v := range postMetas {
		interfaceSlice[i] = v
	}

	insertManyResult, err := postsCollection.InsertMany(c.ctx, interfaceSlice)
	if err != nil {
		return err
	}

	utils.Logger.Info("MongoAction", zap.String("method", "InsertMany"),
		zap.Any("docs", insertManyResult.InsertedIDs))

	return nil
}

func (c *MongoClientCtx) InsertPosterPosts(meta *DBPosterMeta, postMetas []*DBPostMeta) error {
	posterCollection := c.mongoClient.Database(c.dbName).Collection(PosterMetaCollectionName)

	insertOneResult, err := posterCollection.InsertOne(c.ctx, meta)
	if err != nil {
		return err
	}

	utils.Logger.Info("MongoAction", zap.String("method", "InsertSingle"),
		zap.Any("docs", insertOneResult.InsertedID))

	return c.insertPostMetas(postMetas)
}

func (c *MongoClientCtx) UpdatePosterPosts(meta *DBPosterMeta, postMetas []*DBPostMeta) error {
	posterCollection := c.mongoClient.Database(c.dbName).Collection(PosterMetaCollectionName)
	postsCollection := c.mongoClient.Database(c.dbName).Collection(PostsMetaCollectionName)

	filter := bson.M{"posterinfo.platform": meta.PosterInfo.Platform, "posterinfo.userid": meta.PosterInfo.Userid}

	var results []DBPosterMeta
	cur, err := posterCollection.Find(c.ctx, filter)
	if err != nil {
		return err
	}

	for cur.Next(c.ctx) {
		var elem DBPosterMeta
		err := cur.Decode(&elem)
		if err != nil {
			return err
		}
		results = append(results, elem)
	}

	if len(results) > 1 {
		return fmt.Errorf("Duplicated posters")
	} else if len(results) == 0 {
		return c.InsertPosterPosts(meta, postMetas)
	} else {
		searchResult := results[0]
		nowPostsIds := map[string]bool{} // 数据库中现有的posts id

		for _, v := range searchResult.PostRef {
			nowPostsIds[v.PostId] = true
		}

		missPosterIds := map[string]bool{} // 新增的posts id
		for _, v := range meta.PostRef {
			_, exists := nowPostsIds[v.PostId]
			if !exists { // 在数据库中没有
				missPosterIds[v.PostId] = true
			}
		}

		if len(missPosterIds) > 0 {
			newPostMetas := []*DBPostMeta{}

			for _, v := range postMetas {
				_, miss := missPosterIds[v.PostInfoMeta.PostId]
				if miss {
					newPostMetas = append(newPostMetas, v)
				}
			}

			filter := bson.M{
				"id": searchResult.ID,
			}
			update := bson.M{
				"$set": bson.M{
					"id":        meta.ID,
					"fetchtime": meta.FetchTime,
					"postref":   meta.PostRef,
				},
			}
			updateOneResult, err := posterCollection.UpdateOne(c.ctx, filter, update)
			if err != nil {
				return err
			}

			utils.Logger.Info("MongoAction", zap.String("method", "UpdateSingle"),
				zap.Any("docs", updateOneResult.UpsertedID),
				zap.Int64("matched_count", updateOneResult.MatchedCount),
				zap.Int64("modified_count", updateOneResult.ModifiedCount),
				zap.Int64("upserted_count", updateOneResult.UpsertedCount))

			filter = bson.M{
				"postsinfoid": searchResult.ID,
			}
			update = bson.M{
				"$set": bson.M{
					"postsinfoid": meta.ID,
				},
			}
			updateMultiResult, err := postsCollection.UpdateMany(c.ctx, filter, update)
			if err != nil {
				return err
			}

			utils.Logger.Info("MongoAction", zap.String("method", "UpdateMany"),
				zap.Any("docs", updateMultiResult.UpsertedID),
				zap.Int64("matched_count", updateMultiResult.MatchedCount),
				zap.Int64("modified_count", updateMultiResult.ModifiedCount),
				zap.Int64("upserted_count", updateMultiResult.UpsertedCount))

			return c.insertPostMetas(newPostMetas)
		}
	}

	return nil
}

func (c *MongoClientCtx) LinkQuery() ([]*DBLinkQueryResult, error) {
	type linkQueryResult struct {
		ID             primitive.ObjectID `bson:"_id"`
		FileInDataBase bool               `bson:"fileindatabase"`
		PostFiles      []string           `bson:"postfiles"`
		PostDownloads  []string           `bson:"postdownloads"`
	}

	postsCollection := c.mongoClient.Database(c.dbName).Collection(PostsMetaCollectionName)

	filter := bson.M{"fileindatabase": false}
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"_id": 1}) // 1 表示升序，-1 表示降序
	findOptions.SetLimit(100)

	var results []*DBLinkQueryResult
	cur, err := postsCollection.Find(c.ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		var elem linkQueryResult
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		if elem.FileInDataBase {
			return nil, fmt.Errorf("Invalid query in LinkQuery")
		}
		results = append(results, &DBLinkQueryResult{
			DBQueryID:     elem.ID.Hex(),
			PostFiles:     elem.PostFiles,
			PostDownloads: elem.PostDownloads,
		})
	}

	utils.Logger.Info("MongoAction", zap.String("method", "LinkQuery"),
		zap.Int("results_size", len(results)))

	return results, nil
}
