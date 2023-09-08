package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const FILE string = "./hoam.db"
const CREATE_FILE string = "db/create"
var DB *sql.DB

func Start() {
    db, err := sql.Open("sqlite3", FILE)
    if err != nil { fmt.Println("unable to find hoam.db file, exiting"); os.Exit(0) }

    DB = db

    create, fileOpenErr := os.ReadFile(CREATE_FILE);
    if fileOpenErr != nil { fmt.Println("unable to open create file, exiting"); os.Exit(0) }

    _, DBExecErr := DB.Exec(string(create));
    if DBExecErr != nil {
        fmt.Println("unable to create database schema: " + DBExecErr.Error())
        fmt.Println("exiting")
        os.Exit(0)
    }
}

func AddDevice(macAddr string, ip_addr string, name string) (string, error) {
    //come back and finish this later
    command := `
        INSERT OR REPLACE INTO devices (mac, ip, name) VALUES ("`+ macAddr +`", "`+ ip_addr +`", "`+ name +`")
    `

    _, addError := DB.Exec(command);
    if addError != nil { fmt.Println("unable to add name: " + addError.Error()); return "", addError }

    return "", nil 
}

func UpdateIpAddr(macAddr string, ip_addr string) (string, error) {
    command := `
        UPDATE devices
        SET ip = "` + ip_addr + `"
        WHERE mac = "` + macAddr + `"
    `

    result, updateErr := DB.Exec(command);
    if updateErr != nil { fmt.Println("unable to update ip address: " + updateErr.Error()); return "", updateErr }

    rows, rowsErr := result.RowsAffected();
    if rowsErr != nil { fmt.Println("unable to see rows effected" + rowsErr.Error()); return "", rowsErr }

    if rows == 0 { 
        fmt.Println("no rows effected when updating ip");
        return "", errors.New("unable to update name, invalid mac address")
    }

    return "updated ip addr", nil
}

func UpdateName(macAddr string, name string) (string, error) {
    command := `
        UPDATE devices
        SET name = '`+ name +`'
        WHERE mac = '`+ macAddr +`'
    `

    result, updateNameErr := DB.Exec(command);
    if updateNameErr != nil { return "", updateNameErr }

    rows, rowsErr := result.RowsAffected();
    if rowsErr != nil { return "", rowsErr }

    if rows == 0 { 
        fmt.Println("no rows effected when updating name");
        return "", errors.New("unable to update name, invalid mac address")
    }

    return name, nil
}

func GetName(macAddr string) (string, error) {
    query := "SELECT name FROM devices WHERE mac=?"
    row := DB.QueryRow(query, macAddr)

    var name string

    queryErr := row.Scan(&name)
    if queryErr == sql.ErrNoRows { return "", errors.New("unable to index mac addresss") }

    return name, nil
}

func GetMacAddr(ipAddr string) (string, error) {
    query := "SELECT mac FROM devices WHERE ip=?"

    row := DB.QueryRow(query, ipAddr)

    var macAddr string
    queryErr := row.Scan(&macAddr)
    if queryErr == sql.ErrNoRows { return "", errors.New("no mac address matching this ip") }

    return macAddr, nil
}

func UpdateDevice(macAddr string, ipAddr string) (string, error) {
    query := "SELECT count(*) FROM devices WHERE mac=?"
    row := DB.QueryRow(query, macAddr)

    var count int64
    queryErr := row.Scan(&count)
    if queryErr == sql.ErrNoRows { return "", queryErr }

    if count > 0 {
        _, updateIpErr := UpdateIpAddr(macAddr, ipAddr)
        if updateIpErr != nil { return "", updateIpErr }

        return "updated device", nil
    }

    _, addErr := AddDevice(macAddr, ipAddr, "unnamed")
    if addErr != nil { return "", addErr }
    
    return "added new device", nil
}
