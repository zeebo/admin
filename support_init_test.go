package admin

import (
	"flag"
	"fmt"
	"launchpad.net/mgo"
	"log"
	"os"
	"os/exec"
)

var (
	load       = flag.Bool("load", false, "Run mongoimport on the json file for the database")
	export     = flag.Bool("export", false, "Export json files")
	exit       = flag.Bool("exit", false, "Exit after loading")
	sessionurl = flag.String("s", "localhost", "Mongo url for the test database")
	session    *mgo.Session
)

func load_collection(collection string) error {
	cmd := exec.Command("mongoimport", "--drop", "-d", "admin_test", "-c", collection, fmt.Sprintf("json/admin_test.%s.json", collection))
	return cmd.Run()
}

func export_collection(collection string) error {
	file, err := os.Create(fmt.Sprintf("json/admin_test.%s.json", collection))
	if err != nil {
		return err
	}
	defer file.Close()
	cmd := exec.Command("mongoexport", "-d", "admin_test", "-c", collection)
	cmd.Stdout = file

	return cmd.Run()
}

func init() {
	flag.Parse()

	types := []string{"T", "T2", "T5", "T6"}

	//Import: mongoimport --drop -d admin_test -c T admin_test.json
	//Export: mongoexport -d admin_test -c T > admin_test.json

	//before commit:
	//mongoexport -d admin_test -c T > admin_test.json
	//go test -load
	//git commit -a -m 'msg'

	if *export {
		for _, t := range types {
			if err := export_collection(t); err != nil {
				log.Fatalf("Error exporting %s: %s", t, err)
			}
		}
		log.Fatal("Data exported successfully.")
	}

	if *load {
		for _, t := range types {
			if err := load_collection(t); err != nil {
				log.Fatalf("Error loading %s: %s", t, err)
			}
		}
		if *exit {
			log.Fatal("Data imported sucessfully.")
		}
	} else {
		//be sure to load T6 no matter what
		if err := load_collection("T6"); err != nil {
			log.Fatalf("Error loading %s: %s", "T6", err)
		}
	}

	var err error
	session, err = mgo.Mongo(*sessionurl)
	if err != nil {
		log.Fatal("Cannot use that session: %s", err)
	}

}
