package load

import (
	"fmt"
	"log"
	"microservice/auth/global"
	"microservice/config"
	g "microservice/gateway/global"
	"microservice/pkg/color"
	"microservice/pkg/loader"
)

var cfg = &config.Config{}

func init() {
	initConfig()
	initMicroservices()
	info()
}

func initConfig() {
	if err1, err2 := loader.ParseYaml("../env.yaml", cfg, false), loader.ParseYaml("../env.yml", cfg, false); err1 != nil || err2 != nil {
		if err1 != nil {
			log.Fatalln(err1)
		} else if err2 != nil {
			log.Fatalln(err2)
		}
	}

	if _, ok := cfg.Microservices[global.Name]; ok {
		cfg.CurrentMicroservice = cfg.Microservices[global.Name]
	} else {
		log.Fatalf("Microservice definition for %s not found", global.Name)
	}

	global.CFG = cfg
	global.SecretKeyBytes = []byte(cfg.SecretKey)
}

func initMicroservices() {
	if auth, ok := g.CFG.Microservices["auth"]; ok {
		g.AuthService = &auth
	} else {
		log.Fatalln("auth microservice is not defined")
	}
}

func info() {
	fmt.Println(color.Cyan, fmt.Sprintf("\n==%sSystem Info%s==%s\n", color.Yellow, color.Cyan, color.Reset))
	fmt.Printf("Name:\t\t\t%s%s%s\n", color.Blue, g.Name, color.Reset)
	fmt.Printf("Version:\t\t%s%s%s\n", color.Blue, g.Version, color.Reset)
	// TODO: Check active/inactive microservices and print their status
	if g.CFG.Debug {
		fmt.Printf("Debug:\t\t\t%s%v%s\n", color.Red, g.CFG.Debug, color.Reset)
	} else {
		fmt.Printf("Debug:\t\t\t%s%v%s\n", color.Green, g.CFG.Debug, color.Reset)
	}
	fmt.Printf("Address:\t\thttp://%s:%s\n", g.CFG.Gateway.IP, g.CFG.Gateway.Port)
	fmt.Printf("Allowed Origins:\t%v\n", g.CFG.AllowOrigins)
	if g.CFG.AllowHeaders != "" {
		fmt.Printf("Extra Allowed Headers:\t%v\n", g.CFG.AllowHeaders)
	}
	fmt.Print(color.Cyan, "===============\n\n", color.Reset)
}
