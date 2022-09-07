package controllers

import (
	"neft.web/models"
  	"github.com/gin-gonic/gin"

  "net"
  "fmt"
)

type Devices struct {
	db models.DeviceDB
}

func NewDevices(db models.DeviceDB) *Devices {
	return &Devices{
		db: db,
	}
}
//define struct for mac address json
type MacAddress struct {
    Id string
}
func (db *Devices) RetrieveByMac() gin.HandlerFunc {
	return func(context *gin.Context) {
    mac := &MacAddress{Id: ""}
     ifas, err := net.Interfaces()
     if err != nil {
          fmt.Println("hola ",mac)
         return
     }
     for _, ifa := range ifas {
         a := ifa.HardwareAddr.String()
         if a != "" {
             mac = &MacAddress{Id: a}
             fmt.Println("adios ",mac)
             break
         }
     }
    device, err := db.db.ByMac(mac.Id)
    fmt.Println(device)
    context.Next()
    }
}