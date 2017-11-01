package ServiceAPIFactory

import (
	"errors"
	"log"
	"reflect"
	"path/filepath"
	"os"
	"context"
)

const DefaultRegistrationFile = "/RegistryService.json"

func GetServiceInstance(ctx context.Context,t reflect.Type) interface{} {
	var registrationConfigPath string
	//this debug will be checked from the execution parameters to check if the code is on debugging mode
	debug := false
	// get the working directory
	if (!debug) {
		wd, err := os.Getwd()
		if (err != nil) {
			log.Fatal(err.Error())
		}
		_, folder := filepath.Split(wd)
		registrationConfigPath = folder + DefaultRegistrationFile
	}
	// all the microservices apis implement the Make method.
	// im invoking the Make method and passing the current context and registration config path and returning the instance
	instance := reflect.New(t).Elem().Interface()

	var method = reflect.ValueOf(instance).MethodByName("Make")
	if (method == reflect.Value{}) {
		log.Fatal(errors.New("Service of type " + t.String() + " does not implement make method"))
	}

	result := method.Call([]reflect.Value{reflect.ValueOf(ctx),reflect.ValueOf(registrationConfigPath)})
	if len(result) > 0 {
		return result[0].Interface()
	}
	return instance
}
