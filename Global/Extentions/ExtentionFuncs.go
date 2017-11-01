package Extentions

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

func ConvertBsonToStruct(bsonData bson.M, data interface{}) error {
	jData, err := bson.Marshal(bsonData)
	err = bson.Unmarshal([]byte(jData), &data)
	return err
}
func ConvertMapToStruct(bsonData map[string]interface{}, data interface{}) error {
	jData, err := json.Marshal(bsonData)
	err = json.Unmarshal([]byte(jData), &data)
	return err
}
