package main

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo"
	_"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"flag"
	"os"
	"github.com/henrylee2cn/pholcus/common/mahonia"
)
var MongoUrl string
var MongoCollection string

type TomConfi struct {
	Title string
	Mongodb _Mongodb
	Log _Log
	SpiderTarget _SpiderTarget
}
type _Mongodb struct {
	Ip string
	Port string
	Username string
	Password string
	Database string
	Collection string

}
type _Log struct{
	Path string
	FileName string
}
type _SpiderTarget struct{
	Url string
}

var  config TomConfi
var (
	mgoSession *mgo.Session
	dataBase   string
)



func getSession365() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(MongoUrl)
		if err != nil {
			panic(err)
		}
	}

	return mgoSession.Clone()
}

func witchCollection365(collection string, s func(*mgo.Collection) error) error {
	session := getSession365()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}
func GetJokes365(){
	doc, err := goquery.NewDocument(config.SpiderTarget.Url)
	if err != nil{
		log.Fatal(err)
	}
	doc.Find(".mainBody .list_con table .list_item").Each(func(i int, s *goquery.Selection){
		//fmt.Println(s.Find(".des .room").Html())
		/*doc,err := s.Find(".img_list img").Attr("lazy_src")
		if(!err){
			return
		}*/

                fmt.Println(mahonia.NewDecoder("GB2312").ConvertString(s.Find("a").Text()))
		q_insert := func(c *mgo.Collection) error {
			selector := bson.M{
				"html":mahonia.NewDecoder("GB2312").ConvertString(s.Find("a").Text()),
			}

			return c.Insert(selector)
		}
		witchCollection365(MongoCollection, q_insert)
	})


}

func main(){
	//init toml config
        _,err := toml.DecodeFile("config.toml",&config)
	if(err!=nil){
		fmt.Println(err.Error())
		return

	}else{
		logFileName := flag.String("log", config.Log.Path, config.Log.FileName)
		dataBase = config.Mongodb.Database
		MongoUrl = config.Mongodb.Ip+":"+config.Mongodb.Port
		MongoCollection = config.Mongodb.Collection
		flag.Parse()
		logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if logErr != nil {
			fmt.Println("no log file access...exit...")
			//os.Exit(1)
		}
		log.SetOutput(logFile)
	}



	GetJokes365()
}