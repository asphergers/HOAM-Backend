package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"reflect"
)

func getErrorMessageJSON(message string) string {
    response := &ErrorResponse { Error: message }

    responseMarshal, _ := json.Marshal(response)

    return string(responseMarshal)
}

func decodeJson[T any](rawJson io.ReadCloser) (T, error) {
    var response T
    resq, err := ioutil.ReadAll(rawJson)

    if err != nil {
        resp := getErrorMessageJSON("unable to read request body")
        return response, errors.New(resp)
    }

    unmarshalErr := json.Unmarshal(resq, &response)

    if unmarshalErr != nil {
        resp := getErrorMessageJSON("cannot decode because of JSON syntax error or mismatching type")
        return response, errors.New(resp)
    }

    invalidStructErr := validStruct[T](response);

    if invalidStructErr != nil {
        resp := getErrorMessageJSON(invalidStructErr.Error())
        return response, errors.New(resp)
    }

    return response, nil
}

func validStruct[T any](obj T) error {
    structType := reflect.TypeOf(obj)

    if structType.Kind() != reflect.Struct {
        //this should literally never happen but im keeping it in just incase
        return errors.New("somehow the request did not get parsed into a struct")
    }

    structValue := reflect.ValueOf(obj)
    filedNumbers := structValue.NumField()

    for i := 0; i < filedNumbers; i++ {
        field := structValue.Field(i)
        filedName := structType.Field(i).Name

        //isValid := field.IsValid() && !field.IsZero()

        if !field.IsValid() {
            return errors.New("field " + filedName + " is not filled")
        }
    }

    return nil
}

