package db

import (
	"context"
	"encoding/json"
	"io/ioutil"
	db2 "kdownloader/pkg/db"
	"kdownloader/pkg/kemono"
	"os"
	"testing"
)

func TestMongoInsert(t *testing.T) {
	URI := ``
	DBName := `kdb`
	DataSet := `horosuke.json`
	jsonFile, err := os.Open(DataSet)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteVal, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var allMeta kemono.PosterAll
	err = json.Unmarshal(byteVal, &allMeta)
	if err != nil {
		panic(err)
	}

	dbMeta, dbPostMetas := db2.DBTypeConvert(&allMeta)

	cli, err := db2.InitMongo(context.Background(), URI, DBName)
	defer cli.Close()

	if err != nil {
		panic(err)
	}

	err = cli.InsertPosterPosts(dbMeta, dbPostMetas)

	if err != nil {
		panic(err)
	}
}

func TestMongoUpdate(t *testing.T) {
	URI, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		panic("MONGO_URI not found")
	}
	DBName := `kdb`
	DataSet := `horosuke.json`
	jsonFile, err := os.Open(DataSet)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteVal, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var allMeta kemono.PosterAll
	err = json.Unmarshal(byteVal, &allMeta)
	if err != nil {
		panic(err)
	}

	dbMeta, dbPostMetas := db2.DBTypeConvert(&allMeta)

	cli, err := db2.InitMongo(context.Background(), URI, DBName)
	defer cli.Close()

	if err != nil {
		panic(err)
	}

	err = cli.UpdatePosterPosts(dbMeta, dbPostMetas)

	if err != nil {
		panic(err)
	}
}

func TestMongoLinkQuery(t *testing.T) {
	URI, exists := os.LookupEnv("MONGO_URI")
	if !exists {
		panic("MONGO_URI not found")
	}
	DBName := `kdb`

	cli, err := db2.InitMongo(context.Background(), URI, DBName)
	defer cli.Close()

	if err != nil {
		panic(err)
	}

	result, err := cli.LinkQuery()

	if err != nil {
		panic(err)
	}

	for _, v := range result {
		for _, v1 := range v.PostFiles {
			println(v1)
		}
		for _, v1 := range v.PostDownloads {
			println(v1)
		}
	}
}
