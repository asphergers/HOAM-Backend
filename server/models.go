package server

import "main/tcpclient"

// api/toggle
type ToggleRequest struct {
    Address string
}

type ToggleResponse struct {
    message string
}

// api/name
type GetNameRequest struct {
    Address string
}

type SetNameRequest struct {
    Address string
    Name string
}

type NameResponse struct {
    Name string
}

// api/devices
type DevicesResponse struct {
    Devices []tcpclient.Device
}

// api/timer
type GetTimerResponse struct {
    Timers []TimerInfoResponse 
}

type GetTimerRequest struct {

}

type PostTimerRequest struct {
    Address string
    Time string
    Action int
    Id int
}

type TimerInfoResponse struct {
    Address string
    Time string
    Id int
}

type TimerRemoveResponse struct {
    RemovedId int
}

// api/changepin
type PostChangePin struct {
    Address string
    Pin int
    Action int
}

type ChangePinResponse struct {

}

// generic error response, should probably add an error number field or something

type ErrorResponse struct {
    Error string
}
