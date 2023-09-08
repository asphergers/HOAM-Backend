package server

import (
	"fmt"
	"errors"
	"net/http"
)

func middle(w *http.ResponseWriter, req **http.Request) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

    if (*req).Method == http.MethodOptions {
        (*w).WriteHeader(http.StatusOK)
    }
}

func toggle(w http.ResponseWriter, req *http.Request) {
    middle(&w, &req)

    var err error

    switch req.Method {
        case http.MethodPost: { err = handleTogglePost(w, req) }
        
        default: {
            w.WriteHeader(http.StatusBadRequest)
            err = errors.New(getErrorMessageJSON("invalid request method"))
        }
    }

    if err != nil {
        w.Write([]byte(err.Error()))
        return
    }
}

func name(w http.ResponseWriter, req *http.Request) {
    middle(&w, &req)

    var err error

    switch req.Method {
        case http.MethodGet: { err = handleNameGet(w, req) }
        case http.MethodPost: { err = handleNamePost(w, req) }
        case http.MethodOptions: { w.WriteHeader(http.StatusOK) }

        default: {
            w.WriteHeader(http.StatusBadRequest)
            err = errors.New(getErrorMessageJSON("invalid request method"))           
        }
    }

    if err != nil {
        w.Write([]byte(err.Error()))
        return
    }
}

func devices(w http.ResponseWriter, req *http.Request) {
    middle(&w, &req)

    var err error

    switch req.Method {
        case http.MethodGet: { err = handleDevicesGet(w, req) }

        default: {
            w.WriteHeader(http.StatusBadRequest)
            err = errors.New(getErrorMessageJSON("Invalid method"))
        }
    }

    if err != nil {
        w.Write([]byte(err.Error()))
        return
    }
}

func timer(w http.ResponseWriter, req *http.Request) {
    middle(&w, &req)

    var err error
    
    switch req.Method {
        case http.MethodGet: { err = handleTimerGet(w, req) }
        case http.MethodPost: { err = handleTimerPost(w, req) }

        default: {
            w.WriteHeader(http.StatusBadRequest)
            err = errors.New(getErrorMessageJSON("invalid method")) 
        }
    }

    if err != nil {
        w.Write([]byte(err.Error()))
        return
    }
}

func changepin(w http.ResponseWriter, req *http.Request) {
    middle(&w, &req)

    var err error

    switch req.Method {
        case http.MethodPost: { err = handleChangePinPost(w, req) }

        default: {
            w.WriteHeader(http.StatusBadRequest)
            err = errors.New(getErrorMessageJSON("invalid method"))
        }
    }

    if err != nil {
        w.Write([]byte(err.Error()))
        return
    }
}

func StartHttp() {
    fmt.Println("starting https server... ")

    mux := http.NewServeMux()

    mux.HandleFunc("/toggle", toggle)
    mux.HandleFunc("/name", name)
    mux.HandleFunc("/devices", devices)
    mux.HandleFunc("/timer", timer)
    mux.HandleFunc("/changepin", changepin)

    http.ListenAndServe(":8090", mux)
}
