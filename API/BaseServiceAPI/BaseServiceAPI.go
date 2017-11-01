package BaseServiceAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"local/UbotTrade/Global"
	"local/UbotTrade/StaticConfiguration"
	"log"
	"net/http"
	"strconv"
	"context"
	"github.com/nu7hatch/gouuid"
	"local/UbotTrade/Global/GeneralMicroserviceBehavior"
)

var BasicConfig map[string]interface{}

type ServiceAPI interface {
	Make(registrationConfig string) interface{}
}

type BaseServiceAPI struct {
	ServiceInformation *Global.ServiceInformation
	ServiceKey         string
	Ctx            context.Context
}

type IBaseServiceAPI interface {
	SendRequest(req interface{}, method string, responseObj interface{}) (err error)
	SendRequestWithUrl(requestMetadata map[string]string, req interface{}, method string, responseObj interface{}, url string) (err error)
	GetServiceMetrics() (response interface{}, err error)
	GetRequestMetrics() (response Global.ServiceMetricsResponse, err error)
}

//load the static config of each microservice to locate registry service
func LoadConfig(configPath string)  {
	if BasicConfig == nil || len(BasicConfig) > 0 {
		var err error
		BasicConfig, err = StaticConfiguration.ReadConfiguration(configPath)
		if err != nil {
			log.Println(err.Error())
		}
	}
}


func Make(ctx context.Context, basicConfigPath string, serviceKey string, serviceInfo *Global.ServiceInformation) (instance *BaseServiceAPI) {

	instance = &BaseServiceAPI{}
	instance.ServiceKey = serviceKey
	instance.ServiceInformation = serviceInfo

	LoadConfig(basicConfigPath)

	instance.Ctx = ctx
	return instance
}
//the generic Send request is called from all applications accept registry service
func (base *BaseServiceAPI) SendRequest(req interface{}, method string, responseObj interface{}) (err error) {
	if base.ServiceInformation != nil {
		url := base.ServiceInformation.Url + ":" + strconv.Itoa(base.ServiceInformation.Port)

		requestMetadata := HandleRequestMetadata(base.Ctx)
		GeneralMicroserviceBehavior.RegisterMicroserviceHop(requestMetadata, *base.ServiceInformation)

		err = base.SendRequestWithUrl(requestMetadata, req, method, responseObj, url)
	} else {
		err = errors.New("service information is nil")
	}
	return
}
// the actual send request method, accessed only by SendRequest and by registry service.
func (base *BaseServiceAPI) SendRequestWithUrl(requestMetadata map[string]string, req interface{}, method string, responseObj interface{}, url string) (err error) {
	//convert request to json
	reqJson, err := json.Marshal(req)

	if(err != nil){
		log.Println(err.Error())
		return
	}

	//json to bytes
	var jsonStr = []byte(reqJson)

	//create request
	httpReq, err := http.NewRequest("POST", url+"/"+method, bytes.NewBuffer(jsonStr))
	httpReq.Header.Set("Content-Type", "application/json")

	// all the request metadata headers are added to the request here, including request id and user id
	for k, v := range requestMetadata {
		if (httpReq.Header.Get(k) == "" ) {
			httpReq.Header.Add(k, v)
		} else {
			httpReq.Header.Set(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	defer resp.Body.Close()

	//read request body
	body, _ := ioutil.ReadAll(resp.Body)
	if body == nil {
		log.Fatal("response is nill")
	} else {
		//convert the response body to the ref of the object that was passed as responseObj
		json.Unmarshal([]byte(body), responseObj)
	}
	return
}

func (base *BaseServiceAPI) GetServiceMetrics() (response interface{}, err error) {
	requestMetadata := HandleRequestMetadata(base.Ctx)
	err = base.SendRequestWithUrl(requestMetadata, nil, "metrics", &response, base.ServiceInformation.Url+":"+strconv.Itoa(base.ServiceInformation.Port))
	return
}

func (base *BaseServiceAPI) GetRequestMetrics() (response Global.ServiceMetricsResponse, err error) {
	requestMetadata := HandleRequestMetadata(base.Ctx)
	err = base.SendRequestWithUrl(requestMetadata, nil, "requestmetrics", &response, base.ServiceInformation.Url+":"+strconv.Itoa(base.ServiceInformation.Port))
	return
}
// create the basic request metadata and add request id
func HandleRequestMetadata(ctx context.Context) map[string]string {
	var requestMetadata map[string]string
	if (ctx == nil || ctx.Value(Global.ContextKeyRequestMetadata) == nil || ctx == context.Background()) {
		requestMetadata = make(map[string]string)
		requestId, err := uuid.NewV4()
		if (err != nil) {
			log.Panic(err)
		}

		requestMetadata[Global.RequestMetadataKeyRequestId] = requestId.String()
	} else {
		requestMetadata = ctx.Value(Global.ContextKeyRequestMetadata).(map[string]string)
	}

	return  requestMetadata
}
