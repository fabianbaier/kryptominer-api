package core

import (
	"fmt"
	"github.com/fabianbaier/kryptominer-api/config"
	"github.com/fabianbaier/kryptominer-api/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type updateETHWalletRequest struct {
	ActiveWorkers   uint64  `json:"activeWorkers"`
	AverageHashrate float64 `json:"averageHashrate"`
	Balance         float64 `json:"balance"`
	Time            int64   `json:"time"`
}

type apiResponse struct {
	Status string
	Message string
}

func Start(appConfig config.Config) {
	port := fmt.Sprintf(":%d", appConfig.FlagAPIPort)
	logrus.Infof("Starting Kryptominer API on 0.0.0.0%s", port)
	r := gin.Default()
	db, err := storage.GetConnection(appConfig.DBURL)
	if err != nil {

	}

	r.Use(func(c *gin.Context) {
		c.Set("store", db)
		c.Next()
	})

	r.GET("/api/v1/eth/wallets/:address", getETHWalletHandler)
	r.POST("/api/v1/eth/wallets/:address", updateETHWalletHandler)
	r.Run(port)
}

func getETHWalletHandler(c *gin.Context) {
	db, ok := c.Get("store")
	if !ok {
		return
	}
	ethAddress := c.Param("address")
	w, err := db.(*storage.Connection).GetWallet(ethAddress)
	if err != nil {
		message := apiResponse{
			Status:  "FAILED",
			Message: fmt.Sprintf("%s", err),
		}
		c.JSON(500, message)
		return
	} else {
		c.JSON(200, w)
	}
}

func updateETHWalletHandler(c *gin.Context) {
	var req updateETHWalletRequest
	err := c.BindJSON(&req)
	if err != nil {
		// bail out
		logrus.Warningf("Update Wallet Error: Parsing: %s", err)
		return
	}
	addr := c.Param("address")
	payload := storage.Wallet{
		ActiveWorkers:   req.ActiveWorkers,
		Address: addr,
		AverageHashrate: req.AverageHashrate,
		Time:            req.Time,
		Balance:         req.Balance,
	}

	db, ok := c.Get("store")
	if !ok {
		return
	}

	err = db.(*storage.Connection).InsertWalletState(payload)
	if err != nil {
		message := apiResponse{
			Status:  "FAILED",
			Message: fmt.Sprintf("%s", err),
		}
		c.JSON(500, message)
		return
	}
	c.JSON(200, payload)
}
