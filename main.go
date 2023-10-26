package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/dailypenn/content-recommendation/dbutils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no subcommand provided; options are import and update")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	switch os.Args[1] {
	case "import":
		importCmd := flag.NewFlagSet("import", flag.ExitOnError)
		numWorkers := importCmd.Int("w", -1, "number of workers to use when fetching articles")
		importCmd.Parse(os.Args[2:])
		if *numWorkers < 0 {
			log.Fatal("number of workers must be greater than 0")
		}
		dbutils.ImportArticles(*numWorkers, client.Database("Cluster"))
	case "update":
		dbutils.UpdateArticles(client.Database("Cluster"))
	default:
		log.Fatal("unknown subcommand")
	}

}
