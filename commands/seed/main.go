package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/joho/godotenv"
	"github.com/romanyx/polluter"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"go.uber.org/fx"
	"io/ioutil"
	"path/filepath"
)

func main() {
	godotenv.Load(filepath.Join(".env"))

	var (
		seed = flag.String("seed", "basic", "seed name")
	)
	flag.Parse()

	fmt.Println("Use seed: " + *seed)

	ctx := context.Background()
	app := di.NewApp(fx.Invoke(func(conn db.IConnector) {
		db := conn.GetDB()

		p := polluter.New(polluter.MySQLEngine(db))

		seedFile := fmt.Sprintf("%s.yml", *seed)
		seedBytes, err := ioutil.ReadFile(filepath.Join("seeds", seedFile))
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
	}))

	defer app.Stop(ctx)
	err := app.Start(ctx)
	if err != nil {
		panic(err)
	}
}
