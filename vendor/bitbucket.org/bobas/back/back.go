package back

import (
	"net/smtp"

	"os"

	"bitbucket.org/bobas/back/domain"
	"bitbucket.org/bobas/back/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/nildev/lib/registry"
)

// CreateEndpoint is used to register new API endpoint. This will be used by our web app.
//
// Sample request:
//
// curl -X POST -d'{"sampleResponse":"{\"approve\":true}", "handlers":[{"metadata":{"type":"email"}, "value":["kiril@pirozenko.com"]}]}'  http://127.0.0.1:8080/api/v1/endpoints/first -H "Authorization: Bearer <token>" -v
//
// The path of the endpoint will be: /endpoints/{name}
// @method POST
// @path /endpoints/{name}
// @protected
func CreateEndpoint(user *registry.User, name string, sampleResponse string, handlers []domain.Handler) (endpointId string, err error) {
	endp, err := domain.MakeEndpoint(name, user.GetEmail(), sampleResponse, handlers)
	if err != nil {
		log.Errorf("%s", err)
		return "", err
	}

	writer, err := storage.MakeRepositoryFromEnv()
	if err != nil {
		log.Errorf("%s", err)
		return "", err
	}
	err = writer.Save(endp)

	if err != nil {
		log.Errorf("%s", err)
		return "", err
	}
	return endp.Id, err
}

// GetRequestData will return data required to render a form for Bob.
// This will be used by our web app.
//
// Sample request: curl -X GET http://127.0.0.1:8080/api/v1/requests/8ed643fe-9ecc-42b5-4d53-7c7ebf9063ff -H "Authorization: Bearer <token>" -v
//
//
// Response will have format:
// {
//     "sampleResponse" : "\{ ... \}"
// }
//
// @method GET
// @path /requests/{token}
func GetRequestData(token string) (result *domain.RequestResponse, err error) {
	reader, err := storage.MakeRepositoryFromEnv()
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	e, err := reader.GetById(token)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}

	result = &domain.RequestResponse{
		SampleResponse: e.SampleResponse,
	}

	return result, nil
}

// SaveResponse will save response to request. This will be used by our web app.
// Post data will contain form filled by Bob. Request body:
// {
//     "data" : { ... }
// }
//
// Sample request: curl -X POST -d'{"data":{"approve":false}}'  http://127.0.0.1:8080/api/v1/requests/8ed643fe-9ecc-42b5-4d53-7c7ebf9063ff -H "Authorization: Bearer token" -v
//
// @method POST
// @path /requests/{token}
func SaveResponse(token string, data interface{}) (result *string, err error) {
	log.Infof("Got data [%+v]", data)
	// token here is endpoint id
	r := domain.MakeRequest(token, data)
	writer, err := storage.MakeRequestRepositoryFromEnv()
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	err = writer.Save(r)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	result = &r.Id
	return result, nil
}

// HandleEndpointRequest is used by programmer to request data from Bob.
// Optionally you can provide data in the request. Payload format:
// {
//     "data" : { ... }
// }
//
// curl -X POST -d'{}'  http://127.0.0.1:8080/api/v1/endpoints/a280a122-728f-47ad-71f5-f99732575b77/responses -H "Authorization: Bearer <token>" -v
//
//
// Response will contain a link to an endpoint where developer can
// retrieve response data.
//
// @method POST
// @path /endpoints/{eid}/responses
// @protected
func HandleEndpointRequest(user *registry.User, eid string, data interface{}) (result bool, err error) {
	from := "to.bartkus@gmail.com"
	pass := os.Getenv("ND_PASS")
	reader, err := storage.MakeRepositoryFromEnv()
	if err != nil {
		log.Errorf("%s", err)
		return false, err
	}
	e, err := reader.GetById(eid)
	if err != nil {
		log.Errorf("%s", err)
		return false, err
	}

	for _, h := range e.Handlers {
		for _, email := range h.Value {
			msg := "From: " + from + "\n" +
				"To: " + string(email) + "\n" +
				"Subject: Hello! It's me Bob! \n\n" +
				"Please enter data http://127.0.0.1:8080/requests/" + e.Id

			err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{string(email)}, []byte(msg))
			if err != nil {
				log.Errorf("%s", err)
			}
		}
	}

	return true, nil
}

// ResponseData will return data filled by users.
//
// @method GET
// @path /responses/{eid}
// @protected
func ResponseData(eid string) (results []domain.Request, err error) {
	reader, err := storage.MakeRequestRepositoryFromEnv()
	if err != nil {
		log.Errorf("%s", err)
		return []domain.Request{}, err
	}
	results, err = reader.GetByEndpointId(eid)
	if err != nil {
		log.Errorf("%s", err)
		return []domain.Request{}, err
	}

	return results, nil
}
