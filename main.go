package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fabianbaier/kryptominer-api/config"
	"github.com/fabianbaier/kryptominer-api/core"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Kryptominer API v0.0.1")
	fmt.Println("Copyright 2017")

	appConfig := config.Configuration()
	if appConfig.Verbose {
		log.SetLevel(log.DebugLevel)

		// print memory usage every 60 seconds
		go func() {
			for {
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				log.Printf("Alloc = %v MB, TotalAlloc = %v MB, Sys = %v MB, NumGC = %v", m.Alloc>>20, m.TotalAlloc>>20, m.Sys>>20, m.NumGC)
				time.Sleep(600 * time.Second)
			}
		}()
	}
	core.Start(appConfig)
	log.Println("Process finished gracefully")
}
