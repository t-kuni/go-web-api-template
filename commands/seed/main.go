package main

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/joho/godotenv"
	"github.com/romanyx/polluter"
	"github.com/t-kuni/go-web-api-skeleton/wire"
	"io/ioutil"
	"path/filepath"
)

func main() {
	godotenv.Load(filepath.Join(".env"))

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	db := app.DBConnector.GetDB()
	p := polluter.New(polluter.MySQLEngine(db))

	seedBytes, err := ioutil.ReadFile(filepath.Join("seeds", "seeds.yml"))
	if err != nil {
		panic(err)
	}

	var seeds map[string]interface{}
	err = yaml.Unmarshal(seedBytes, &seeds)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		panic(err)
	}
	defer db.Exec("SET FOREIGN_KEY_CHECKS=1;")

	for table, _ := range seeds {
		_, err := db.Exec("TRUNCATE TABLE " + table)
		if err != nil {
			panic(err)
		}
	}

	if err := p.Pollute(bytes.NewReader(seedBytes)); err != nil {
		panic(err)
	}

	for table, _ := range seeds {
		fmt.Println("Seed: " + table)
	}
	fmt.Println("Seeding successfully!")
}
