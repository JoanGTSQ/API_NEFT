package controllers

import (
	"github.com/gin-gonic/gin"
	"neft.web/logger"
	"neft.web/models"

	"net"
	"net/http"
)

type Devices struct {
	db models.DeviceDB
}

func NewDevices(db models.DeviceDB) *Devices {
	return &Devices{
		db: db,
	}
}

// define struct for mac address json
type MacAddress struct {
	Id string
}

func (db *Devices) RetrieveByMac() gin.HandlerFunc {
	return func(context *gin.Context) {
		var mac MacAddress
		ifas, err := net.Interfaces()
		if err != nil {
			logger.Warning.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				mac = MacAddress{Id: a}
				break
			}
		}
		_, err = db.db.ByMac(mac.Id)
		switch err {
		case models.ERR_NOT_FOUND:
			context.JSON(http.StatusBadRequest, gin.H{"error": "You don't have permission to enter beta url"})
			context.Abort()
		case nil:

		default:
			logger.Warning.Println(err)
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}
