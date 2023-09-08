package server

import (
	"encoding/json"
	"errors"
	"main/db"
	"main/tcpclient"
	"net/http"
	"time"
)

func handleTogglePost(w http.ResponseWriter, req *http.Request) error {
    request, err := decodeJson[ToggleRequest](req.Body) 

    if err != nil {
        return errors.New(getErrorMessageJSON(err.Error()))
    }

    //come back here and get response
    _, toggleErr := tcpclient.Toggle(request.Address)

    if toggleErr != nil {
        return errors.New(getErrorMessageJSON(err.Error()))
    }

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte("its fine"))

    return nil
}

func handleNamePost(w http.ResponseWriter, req *http.Request) error {
    // this looks like shit
    request, err := decodeJson[SetNameRequest](req.Body)
    if err != nil { return errors.New(getErrorMessageJSON(err.Error())) }

    macAddr, macErr := db.GetMacAddr(request.Address);
    if macErr != nil { return errors.New(getErrorMessageJSON(macErr.Error())) }

    name, nameErr := db.UpdateName(macAddr, request.Name)
    if nameErr != nil { return errors.New(getErrorMessageJSON(nameErr.Error())) }

    objectResponse := NameResponse { Name: name }

    jsonResponse, jsonErr := json.Marshal(&objectResponse)
    if jsonErr != nil { return errors.New(getErrorMessageJSON(jsonErr.Error())) }

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(jsonResponse))

    return nil
}

func handleNameGet(w http.ResponseWriter, req *http.Request) error {
    request, err := decodeJson[GetNameRequest](req.Body)
    if err != nil { return errors.New(getErrorMessageJSON(err.Error())) }

    //We're getting the MAC address from the database then getting the name
    //using the MAC address. This is probably a bad way of doing things. Too bad!
    macAddr, macErr := db.GetMacAddr(request.Address);
    if macErr != nil { return errors.New(getErrorMessageJSON(macErr.Error())) }

    name, nameErr := db.GetName(macAddr)
    if nameErr != nil { return errors.New(getErrorMessageJSON(nameErr.Error())) }

    objectResponse := NameResponse { Name: name }

    jsonResponse, jsonErr := json.Marshal(&objectResponse)
    if jsonErr != nil { return errors.New(getErrorMessageJSON(jsonErr.Error())) }

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(jsonResponse))

    return nil
}

func handleDevicesGet(w http.ResponseWriter, req *http.Request) error {
    devices := tcpclient.Devices()

    response := DevicesResponse { devices }
    
    jsonResponse, jsonErr := json.Marshal(&response)
    if jsonErr != nil { return errors.New(getErrorMessageJSON("UNABLE TO CONVERT DEVICE LIST TO JSON")) }

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(jsonResponse))

    return nil
}

func handleTimerGet(w http.ResponseWriter, req *http.Request) error {
    timerResponses := make([]TimerInfoResponse, 0)

    for _, timer := range timers {
        t := TimerInfoResponse {
            Address: timer.Address,
            Time: timer.FireTime.String(),
            Id: timer.Id,
        }

        timerResponses = append(timerResponses, t)
    }

    response := GetTimerResponse { Timers: timerResponses }
    responseJson, jsonErr := json.Marshal(&response)

    if jsonErr != nil { return errors.New(getErrorMessageJSON("unable to convert timers into json")) }

    w.WriteHeader(http.StatusAccepted)
    w.Write(responseJson)
    
    return nil
}

func handleTimerPost(w http.ResponseWriter, req *http.Request) error {
    request, _ := decodeJson[PostTimerRequest](req.Body)

    // don't want to check if everything is filed out here
    //if decodeErr != nil {
    //    return decodeErr
    //}

    switch request.Action {
        case 0: { 
            responseJson, addErr := timerAddAction(request)

            if addErr != nil { return addErr }
            
            w.WriteHeader(http.StatusAccepted)
            w.Write([]byte(responseJson))
        }

        case 1: {
            responseJson, rmErr := timerRemoveAction(request) 

            if rmErr != nil { return rmErr }

            w.WriteHeader(http.StatusAccepted)
            w.Write([]byte(responseJson))
        }

        default: {
            return errors.New(getErrorMessageJSON("invalid action id"))
        }
    }

    return nil
}

func timerAddAction(request PostTimerRequest) (string, error) {
    time, err := time.Parse(time.RFC3339, request.Time)

    if err != nil {
        return "", errors.New(getErrorMessageJSON("unable to parse time"))
    }

    success := addTimer(time, request.Address)

    switch success {
        case TIMER_OK: { break }

        case BEFORE_TIME: {
            return "", errors.New(getErrorMessageJSON("time sent is before the current time"))
        }

        default: {
            return "", errors.New(getErrorMessageJSON("something terribly wrong has happened"))
        }
    }

    response := TimerInfoResponse {
        Address: request.Address,
        Time: request.Time,
        Id: currentId,
    }

    jsonResponse, jsonErr := json.Marshal(&response)
    if jsonErr != nil { return "", errors.New(getErrorMessageJSON("unable to convert device list to json")) }

    return string(jsonResponse), nil
}

func timerRemoveAction(request PostTimerRequest) (string, error) {
    success := removeTimer(request.Id);

    switch success {
        case TIMER_OK: { break }

        case INVALID_TIMER_ID: {
            return "", errors.New(getErrorMessageJSON("the id provided is not in the schedule pool"))
        }

        default: {
            return "", errors.New(getErrorMessageJSON("something has gone awfully wrong"))
        }
    }

    response := TimerInfoResponse {
        Id: request.Id,
    }

    jsonResponse, jsonErr := json.Marshal(&response)

    if jsonErr != nil {
        return "", errors.New(getErrorMessageJSON("UNABLE TO CONVERT DEVICE LIST TO JSON"))
    }

    return string(jsonResponse), nil
}

func handleChangePinPost(w http.ResponseWriter, req *http.Request) error {
    request, decodeErr := decodeJson[PostChangePin](req.Body)

    if decodeErr != nil {
        return errors.New(decodeErr.Error())
    }

    response, changePinErr := tcpclient.ChangePin(request.Address, request.Pin, request.Action);
    if changePinErr != nil { return errors.New(changePinErr.Error()) }

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte(response))
        
    return nil
}
