package load

import (
	"fmt"
	"log"
	"github.com/Alireza7997/go_microservices/config"
	g "github.com/Alireza7997/go_microservices/gateway/global"
	"github.com/Alireza7997/go_microservices/pkg/color"
	"github.com/Alireza7997/go_microservices/pkg/loader"
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

	if ms, ok := cfg.Microservices[g.Name]; ok {
		cfg.CurrentMicroservice = ms
	}

	g.CFG = cfg
}

func initMicroservices() {
	for _, name := range []string{"auth", "chat", "greet"} {
		if _, ok := g.CFG.Microservices[name]; !ok {
			log.Fatalf("microservice definition for '%s' not found", name)
		}
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
