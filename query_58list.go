package main

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo"
	"flag"
	"fmt"
)
const URL = "192.168.30.125:27017"

var (
	mgoSession *mgo.Session
	dataBase   = "bus"
)
var (
	logFileName = flag.String("log", "/var/log/spider.log", "Log file name")
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic(err)
		}
	}

	return mgoSession.Clone()
}

func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}
func GetJokes(){
	doc, err := goquery.NewDocument("http://nj.58.com/chuzu/")
	if err != nil{
		log.Fatal(err)
	}
	doc.Find(".des").Each(func(i int, s *goquery.Selection){
		fmt.Println(s.Html())
		q_insert := func(c *mgo.Collection) error {
			selector := bson.M{
				"html":s.Text(),
			}

			return c.Insert(selector)
		}
		witchCollection("pf_58zufang", q_insert)
	})


}

func main(){
	GetJokes()
}