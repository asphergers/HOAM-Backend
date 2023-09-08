package tcpclient

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

    "main/db"
)

type Device struct {
    Address string
    Name string
    Mac string
}

const PORT = 225
const DEFAULT_TIMEOUT = 1000

var Connections []Device

func Start() {
    addrs := getIPv4Addrs()
    devices := pingNetwork(addrs)

    fmt.Println("\ndevices on port 225")
    for _, device := range devices {
        port := PORT 
        addr := fmt.Sprintf("%s:%d", device.String(), port)

        fmt.Println(addr)
    }
}

//address and message should be swapped here to follow the rest of the program, too bad!
func sendMessage(message string, address string) (string, error) {
    if !strings.Contains(message, " ") {
        message += " "
    }

    addr := address + ":" + strconv.Itoa(PORT)
    connection, err := net.DialTimeout("tcp", addr, 2000*time.Millisecond)

    if err != nil {
        return "", errors.New("UNABLE TO CONNECT TO TCP SERVER")
    }

    connection.Write([]byte(message))

    responseBuffer := make([]byte, 1024)
    _, tcpResponseError := connection.Read(responseBuffer)

    if tcpResponseError != nil {
        return "", errors.New("CONNECTION TO THE TCP SERVER WAS MADE BUT SERVER DID NOT RESPONSE IN TIME")
    }

    trimmedResponseBuffer := bytes.Trim(responseBuffer, "\x00")
    
    return parseEPSResponse(string(trimmedResponseBuffer));
}

func ChangeName(newName string, address string) (string, error) {
    command := "SETNAME " + newName

    response, err := sendMessage(command, address)

    if err != nil {
        return "", errors.New(err.Error())
    }

    return response, nil
}

func GetName(address string) (string, error) {
    command := "GETNAME"

    response, err := sendMessage(command, address)

    if err != nil {
        return "", errors.New(err.Error())
    }

    return response, nil
}

func Toggle(address string) (string, error) {
    command := "TOGGLE "

    response, err := sendMessage(command, address)

    if err != nil {
        return "", errors.New(err.Error())
    }

    return response, nil
}

func Devices() []Device {
    addrs := getIPv4Addrs()
    validAddrs := pingNetwork(addrs)

    var devices []Device

    for _, addr := range validAddrs {
        macAddr, updateErr := updateDeviceDB(addr.String());
        if updateErr != nil { continue }
        
        deviceName, getNameErr := db.GetName(macAddr);
        if getNameErr != nil { continue }

        device := Device { Address: addr.String(), Name: deviceName, Mac: macAddr }
        devices = append(devices, device)
    }

    return devices
}

func ChangePin(address string, pin int, action int) (string, error) {
    command := "CHANGEPIN " + strconv.Itoa(pin) + " " + strconv.Itoa(action)

    response, err := sendMessage(command, address)

    if err != nil {
        return "", errors.New(err.Error())
    }

    return response, nil
}

func Shake(address string) (string, error) {
    command := "SHAKE";

    response, err := sendMessage(command, address);
    if err != nil { return "", errors.New(err.Error()) }

    return response, nil
}

func updateDeviceDB(ipAddr string) (string, error) {
    macAddr, shakeErr := Shake(ipAddr)
    if shakeErr != nil { return "", errors.New("unable to get mac address") }

    response, updateErr := db.UpdateDevice(macAddr, ipAddr)
    if updateErr != nil { 
        fmt.Println("update error: " + updateErr.Error())
        return "", updateErr 
    }

    fmt.Println("update device response: " + response);

    return macAddr, nil
}

func parseEPSResponse(response string) (string, error){
    left, right, found := strings.Cut(response, ":");

    if !found { return "", errors.New("unable to parse esp32 response")}

    if left == "ERR" { return "", errors.New(right) }

    return right, nil
}

func getIPv4Addrs() []net.IP {
    _, ipv4Net, err := net.ParseCIDR("192.168.1.0/25")

    if err != nil { log.Fatal(err) }

    mask := binary.BigEndian.Uint32(ipv4Net.Mask)
    start := binary.BigEndian.Uint32(ipv4Net.IP)

    finish := (start & mask) | (mask ^ 0xffffffff)

    var addrs []net.IP

    for i := start; i <= finish; i++ {
        ip := make(net.IP, 4)
        binary.BigEndian.PutUint32(ip, i)
        addrs = append(addrs, ip)
    }

    return addrs
}

func pingNetwork(addrs []net.IP) []net.IP {
    var valid_addrs []net.IP

    wg := &sync.WaitGroup{}

    for _, addr := range addrs {
        wg.Add(1)
        go func(addr net.IP) {
            if !isDevice(addr) {
                wg.Done()
                return
            } 

            //COME BACK AND USE THIS WHEN DB IS SETUP
            _, err := Shake(addr.String());
            if err != nil {
                fmt.Println("failed handshake: " + addr.String())
                wg.Done()
                return
            }

            fmt.Println("\nvalid address at:", addr.String())
            valid_addrs = append(valid_addrs, addr)

            wg.Done()
        }(addr)
    }

    wg.Wait()

    return valid_addrs
}

func isDevice(addr net.IP) bool {
    a := fmt.Sprintf("%s:%d", addr.String(), PORT)

    conn, err := net.DialTimeout("tcp", a, DEFAULT_TIMEOUT*time.Millisecond)
    if err != nil { return false }

    conn.Close()

    return true
}

func handShakeDevice(addr *net.IP) bool {
    message := "SHAKE"

    _, err := sendMessage(message, addr.String())
    if err != nil { return false }

    return true
}
