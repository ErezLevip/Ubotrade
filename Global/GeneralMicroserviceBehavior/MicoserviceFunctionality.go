package GeneralMicroserviceBehavior

import (
	"github.com/erezlevip/Ubotrade/Global"
	"encoding/json"
	"context"
	"log"
	"net/http"
)

func RegisterMicroserviceHop(requestMetadata map[string]string, serviceInfo Global.ServiceInformation) {

	var err error
	if (len(requestMetadata) == 0 || requestMetadata[Global.RequestMetadataKeyHops] == "") {
		requestMetadata, err = appendHopToMetadata(requestMetadata, nil, serviceInfo)
	} else {
		servicesInfoStr := requestMetadata[Global.RequestMetadataKeyHops]
		var servicesInfo []Global.ServiceInformation
		json.Unmarshal([]byte(servicesInfoStr), &servicesInfo)
		requestMetadata, err = appendHopToMetadata(requestMetadata, servicesInfo, serviceInfo)
	}

	if (err != nil) {
		log.Panic(err.Error())
	}

	return
}

func appendHopToMetadata(metadata map[string]string, servicesInfo []Global.ServiceInformation, serviceInfo Global.ServiceInformation) (map[string]string, error) {
	if (servicesInfo == nil) {
		servicesInfo = make([]Global.ServiceInformation, 0)
	}
	servicesInfo = append(servicesInfo, serviceInfo)

	jdata, err := json.Marshal(servicesInfo)
	if (err != nil) {
		log.Println(err.Error())
		return nil, err
	}

	metadata[Global.RequestMetadataKeyHops] = string(jdata)
	return metadata, err
}

func DecodeRequestWithMetadata(request interface{}, r *http.Request) (genericRequest map[string]interface{}, err error) {

	genericRequest = make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		//return nil, err
		log.Fatal(err.Error())
	}

	genericRequest["Request"] = request

	header := r.Header.Get(Global.RequestMetadataHeader)
	if (header != "") {
		var metadata map[string]string
		err = json.Unmarshal([]byte(header), &metadata)
		if (err != nil) {
			log.Println(err.Error())
			return
		}

		genericRequest[Global.ContextKeyRequestMetadata] = metadata
	}
	return genericRequest, err
}

func GetRequestFromGenericRequest(genericRequest interface{}, requestInstance interface{}) {
	if (genericRequest == nil) {
		return
	}

	var requestData = genericRequest.(map[string]interface{})["Request"]
	jdata, err := json.Marshal(requestData)
	if (err != nil) {
		log.Panic(err.Error())
		return
	}
	err = json.Unmarshal(jdata, requestInstance)
	if (err != nil) {
		log.Panic(err.Error())
		return
	}
	return
}

func GetContextFromGenericRequest(ctx context.Context, genericRequest interface{}) context.Context {
	if (ctx == nil) {
		ctx = context.Background()
	}

	if (genericRequest != nil) {
		var metaData = genericRequest.(map[string]interface{})[Global.ContextKeyRequestMetadata]
		if (metaData != nil) {
			requestMetadata := metaData.(map[string]string)
			ctx = context.WithValue(ctx, Global.ContextKeyRequestMetadata, requestMetadata)
		}
	}

	return ctx
}
