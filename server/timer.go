package server

import (
	"main/tcpclient"
	"time"
)

type TimerStatus uint64

const (
    TIMER_OK TimerStatus = 0
    BEFORE_TIME TimerStatus = 1
    INVALID_TIMER_ID TimerStatus = 2
)

type TimerTracker struct {
    Id int
    Address string
    Timer *time.Timer
    FireTime time.Time
}

var timers []TimerTracker
var currentId int

func addTimer(fireTime time.Time, address string) TimerStatus {
    if fireTime.Before(time.Now()) { return BEFORE_TIME }

    currentId += 1

    duration := fireTime.Sub(time.Now())

    timer := time.NewTimer(duration)
    id := currentId

    tracker := TimerTracker {
        Id: id,
        Address: address,
        Timer: timer,
        FireTime: fireTime,
    }

    timers = append(timers, tracker)

    go func(timer *time.Timer, address string, id int) {
        <-timer.C
        tcpclient.Toggle(address)
        //come back later and add checking for if the toggle worked
        removeTimer(id)
    }(timer, address, id)

    return TIMER_OK
}

func removeTimer(id int) TimerStatus {
    index := -1

    for i, timer := range timers {
        if timer.Id == id {
            index = i
            break
        }
    } 

    if index == -1 { return INVALID_TIMER_ID }

    timers[index].Timer.Stop()

    timers[index] = timers[len(timers)-1]
    timers = timers[:len(timers)-1]

    return TIMER_OK
}

