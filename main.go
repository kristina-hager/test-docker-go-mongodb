package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Person struct {
	Name  string
	Phone string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello there you beautiful world:  I was built on ", runtime.GOOS, " with a CPU ", runtime.GOARCH, "\n")

	log.Print("db uri: ", os.Getenv("DB_PORT_27017_TCP_ADDR"))
	//log.Print("env vars: ", os.Environ())

	//KH: seems strange to me to rely on DB_PORT (name the commandline gave to linked container)
	//instructions online call for relying on MONGOHQ_URL
	//is there some way to programmatically set MONGOHQ_URL to the same val as appropriate env var from docker CLO?
	uri := os.Getenv("DB_PORT_27017_TCP_ADDR")
	if uri == "" {
		fmt.Fprint(w, "no db connection port provided")
		log.Print("read-no db connect port provided")
		return
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		fmt.Fprint(w, "Can't connect to mongo, go error: ", err)
		log.Fatal("read-Can't connect to mongo, go error: ", err)
		return
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	c := sess.DB("test").C("people")

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		fmt.Fprint(w, "Fatal error: ", err, "\n")
		log.Print("find fail: ", err)
		return
	}

	fmt.Fprint(w, "retrieved phone number: ", result.Phone, "\n")
}

func loadDatabase() {
	log.Print("---start write to database---")
	//KH: seems strange to me to rely on DB_PORT (name the commandline gave to linked container)
	//instructions online call for relying on MONGOHQ_URL
	//is there some way to programmatically set MONGOHQ_URL to the same val as appropriate env var from docker CLO?
	uri := os.Getenv("DB_PORT_27017_TCP_ADDR")
	if uri == "" {
		log.Fatal("no db connect port provided")
		return
	}

	sess, err := mgo.Dial(uri)
	if err != nil {
		log.Fatal("write-Can't connect to mongo, go error", err)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	//KH hard-coded names.. fix
	c := sess.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Bob", "+55 53 1234 1234"})
	if err != nil {
		log.Fatal(err)
	}
	log.Print("---wrote to database---")

}

func main() {
	loadDatabase()
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
