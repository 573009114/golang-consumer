package controllers

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"time"
	"unsafe"

	"github.com/BurntSushi/toml"
	"gopkg.in/olivere/elastic.v6"
)

var (
	client *elastic.Client
	err    error
	ctx    = context.Background()
)

type Item struct {
	Msg            string    `json:"msg"`
	TraceID        string    `json:"trace_id"`
	ServiceID      string    `json:"service.id"`
	Module         string    `json:"module"`
	LogLevel       string    `json:"log_level"`
	Type           string    `json:"type"`
	AppName        string    `json:"app_name"`
	Path           string    `json:"path"`
	Caller         string    `json:"caller"`
	ServiceVersion string    `json:"service.version"`
	Version        string    `json:"@version"`
	Host           string    `json:"host"`
	Timestamp      time.Time `json:"timestamp"`
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func BytesToMyStruct(b []byte) *Item {
	return (*Item)(unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	))
}
func init() {
	//读取config文件内容
	var config cfg
	if _, err := toml.DecodeFile("./config/config.toml", &config); err != nil {
		fmt.Println(err)
	}

	client, err = elastic.NewClient(
		elastic.SetURL(config.Es.Endpoint),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(60*time.Second),
		elastic.SetGzip(true),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func BulkIndexToEs(kafka_topic string, data []byte) (err error) {

	msg := &Item{}

	// 序列化kafka接收的数据
	err = json.Unmarshal(data, &msg)
	if err != nil {
		log.Print(err)
	}

	// 获取日志文件名
	filesname := filepath.Base(msg.Path)

	//转换时间为CST时间
	timers := msg.Timestamp.Local()
	t := timers.Format("2006-01-02")

	//拼接索引
	indexname := filesname + "-" + t
	fmt.Println(indexname)

	//写入
	_, err = client.Index().
		Index(indexname).
		Type(kafka_topic).
		BodyJson(msg).
		Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func SearchEsBody(indexName string, field_name string, field_value string) (res int64) {
	//查询条件
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Should(elastic.NewMatchQuery(field_name, field_value))

	searchResult, err := client.Search().
		Index(indexName). // 设置索引名
		Query(boolQuery). // 设置查询条件
		From(0).          // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(10).         // 设置分页参数 - 每页大小
		Pretty(true).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	// fmt.Println(reflect.TypeOf(searchResult.TotalHits()))
	return searchResult.TotalHits()
}

// func main() {

// 	//查询es
// 	indexName := "logstash-20220201"
// 	field_name := "title"
// 	field_value := "111"
// 	r := SearchEsBody(indexName, field_name, field_value)

// 	//报警
// 	alarm_rule := int64(7)
// 	if r > alarm_rule {
// 		fmt.Println(r)
// 	}
// }
