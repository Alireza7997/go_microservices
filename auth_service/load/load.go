package load

import (
	"fmt"
	"log"
	"microservice/auth/global"
	"microservice/config"
	"microservice/pkg/database"
	"microservice/pkg/loader"

	migrate "github.com/rubenv/sql-migrate"
)

var cfg = &config.Config{}

func init() {
	initConfig()
	initDatabase()
	migrateLatestChanges()
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
	global.SecretKeyBytes = []byte(cfg.SecretKey)
	global.CFG = cfg
}

func initDatabase() {
	var err error
	global.AllSQLCons, global.DB, err = database.New(cfg.CurrentMicroservice.Databases, cfg.Debug)
	if err != nil {
		log.Fatalf("error while creating database : %v", err)
	}

	if global.CFG.Debug {
		_, ok := global.AllSQLCons["test"]
		if !ok {
			log.Fatalln("the 'test' db is not defined")
		}
	} else {
		_, ok := global.AllSQLCons["main"]
		if !ok {
			log.Fatalln("the 'main' db is not defined")
		}
	}

}

func migrateLatestChanges() {
	db, err := global.DB()
	if err != nil {
		panic(err)
	}
	mainOrTest := "test"
	if !global.CFG.Debug {
		mainOrTest = "main"
	}
	migrations := &migrate.FileMigrationSource{
		Dir: fmt.Sprintf("migrations/%s/", mainOrTest),
	}

	n, err := migrate.Exec(db, global.CFG.CurrentMicroservice.Databases[mainOrTest].Type, migrations, migrate.Up)
	if err != nil {
		log.Fatalln(err)
	}

	if n > 0 {
		fmt.Println("\n==Migrations==")
		fmt.Printf("Applied %d migrations!\n", n)
	}
}
