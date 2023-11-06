package repository

import (
	"context"
	"server/mongo"

)
type Report struct {
	Title string      `json:"title"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func InsertRepot(report *Report) (interface{}, error) {
	db, err := mongo.NewMongoDB(mongo.MongoDBURL)
	if err != nil {
		return nil, err
	}
	insertID, err := db.InsertOne(context.Background(), "report", &report)
	if err != nil {
		return nil, err
	}
	return insertID, nil
}

func GetAllReport() ([]Report, error) {
	var taskArr []Report
	db, err := mongo.NewMongoDB(mongo.MongoDBURL)
	if err != nil {
		return nil, err
	}
	err = db.QueryAll(context.Background(), "report", nil, 0, 10, &taskArr)
	if err != nil {
		return nil, err
	}
	return taskArr, nil
}