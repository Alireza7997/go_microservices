package app

import (
	"fmt"
	"log"

	"github.com/alireza/global"
	"github.com/gin-gonic/gin"
)

func Api() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.MaxMultipartMemory = global.CFG.MaxUploadSize << 20

	err := router.Run(fmt.Sprintf("%s:%s", global.CFG.Host, global.CFG.Port))
	if err != nil {
		log.Fatal(err)
	}

}
