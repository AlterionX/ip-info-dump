package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	errorTypes "github.com/AlterionX/ip-info-dump/infosource/base"
	"github.com/AlterionX/ip-info-dump/infosource"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func httpStatusResponse(status int) events.APIGatewayProxyResponse {
    return events.APIGatewayProxyResponse {
        StatusCode: status,
        Body: http.StatusText(status),
    }
}

func clientErrorResponse(status int) events.APIGatewayProxyResponse {
    return httpStatusResponse(status)
}

func serverErrorResponse() events.APIGatewayProxyResponse {
    return httpStatusResponse(http.StatusInternalServerError)
}

func RunIPLookup(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    log.Printf("Running ip lookup for: %v and %v", ctx, req)

    data, err := infosource.GetInfo(req.QueryStringParameters["q"], infosource.GetAllSources())
    if err == nil {
        log.Printf("Running GetInfo resulted in the following data: %v", data)
    } else {
        if errors.Is(err, errorTypes.BadArgument) {
            log.Printf("Attempting to look up domain/IP %q failed due to %q.", req.Body, err.Error());
            return clientErrorResponse(http.StatusBadRequest), nil
        } else {
            // GetInfo is only supposed to return known values, so if something else happens...
            return serverErrorResponse(), nil
        }
    }

    result, err := json.Marshal(data)
    if err != nil {
        log.Printf("Marshalling returned values from functions failed due to %q.", err.Error());
        return serverErrorResponse(), nil
    } else {
        log.Printf("Returned data was marshalled into %q.", result)
    }

    response := httpStatusResponse(http.StatusOK)
    response.Body = string(result)
    return response, nil
}

func main() {
    lambda.Start(RunIPLookup)
}
