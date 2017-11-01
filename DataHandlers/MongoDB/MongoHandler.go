package MongoDB

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const ConnectionStringKey = "ConnectionString"
const DatabaseKey = "Database"

type IDbHandler interface {
	Insert(collection string, valuesBson ...bson.M) error
	FindFirst(collection string, query bson.M) (map[string]interface{}, error)
	FindFMany(collection string, query bson.M) ([]map[string]interface{}, error)
	Update(collection string, query bson.M, update bson.M) error
	UpdateMany(collection string, query bson.M, update bson.M) error
	Delete(collection string, query bson.M) error
	Init(config map[string]string)
}

type MongoHandler struct {
	session        *mgo.Session
	configurations map[string]string
}

func (self *MongoHandler) Init(config map[string]string) {
	self.configurations = config
}

func (self *MongoHandler) init() {
	session, err := mgo.Dial(self.configurations[ConnectionStringKey])
	handleErrors(err)
	session.SetMode(mgo.Monotonic, true) // wtf is that?
	self.session = session
}

func handleErrors(error error) {
	if error != nil {
		log.Panic(error.Error())
	}
}

func (self *MongoHandler) Insert(collection string, valuesBson ...bson.M) error {
	self.init()
	// db name and collection name, will create if the collection does not exists
	c := self.session.DB(self.configurations[DatabaseKey]).C(collection)
	for i := range valuesBson {
		err := c.Insert(valuesBson[i])
		if err != nil {
			return err
		}
	}
	defer self.session.Close()
	return nil
}

func (self *MongoHandler) FindFirst(collection string, query bson.M) (map[string]interface{}, error) {
	self.init()
	c := self.session.DB(self.configurations[DatabaseKey]).C(collection)
	var res bson.M
	err := c.Find(query).One(&res)
	defer self.session.Close()
	return ConvertBsonMToMap(res), err
}

func (self *MongoHandler) FindFMany(collection string, query bson.M) ([]map[string]interface{}, error) {
	self.init()
	var resValue []bson.M
	err := self.session.DB(self.configurations[DatabaseKey]).C(collection).Find(query).All(&resValue)
	//handleErrors(err)
	if resValue != nil {
		mapResult := make([]map[string]interface{}, len(resValue))
		for i := range resValue {
			mapResult[i] = ConvertBsonMToMap(resValue[i])
		}
		defer self.session.Close()
		return mapResult, err
	}
	return nil, nil
}

func (self *MongoHandler) Update(collection string, query bson.M, update bson.M) error {
	self.init()
	difQuery := bson.M{"$set": update}
	err := self.session.DB(self.configurations[DatabaseKey]).C(collection).Update(query, difQuery)
	//handleErrors(err)
	defer self.session.Close()
	return err
}

func (self *MongoHandler) UpdateMany(collection string, query bson.M, update bson.M) error {
	self.init()
	difQuery := bson.M{"$set": update}
	_, err := self.session.DB(self.configurations[DatabaseKey]).C(collection).UpdateAll(query, difQuery)
	//handleErrors(err)
	defer self.session.Close()
	return err
}

func (self *MongoHandler) Delete(collection string, query bson.M) error {
	self.init()
	err := self.session.DB(self.configurations[DatabaseKey]).C(collection).Remove(query)
	//handleErrors(err)
	defer self.session.Close()
	return err
}

func ConvertBsonMToMap(input bson.M) map[string]interface{} {
	if input == nil {
		return nil
	}
	output := make(map[string]interface{})
	for i := range input {
		output[i] = input[i]
	}
	return output
}
