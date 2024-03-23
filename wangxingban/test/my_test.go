package test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getConnection() (*mongo.Collection, error) {
	dsn := "mongodb://root:root@localhost:27022"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	collection := client.Database("test_db").Collection("test_col")
	return collection, nil
}

func getConnection2() (*mongo.Collection, error) {
	dsn := "mongodb://root:root@localhost:27022"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	collection := client.Database("237_asm").Collection("ip_record")
	return collection, nil
}

type person struct {
	// omitempty 选项表示如果 ID 字段是零值则在 BSON 中省略这个字段，让 MongoDB 自动生成 _id。
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func Test1(t *testing.T) {
	collection, err := getConnection()
	if err != nil {
		t.Fatal(err)
	}

	ps := []interface{}{
		person{ID: "1", Name: "Tom", Age: 10},
		person{ID: "2", Name: "Tony", Age: 20},
		person{ID: "3", Name: "Boby", Age: 15},
	}

	insertManyResult, err := collection.InsertMany(context.Background(), ps)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("insertManyResult: %+v\n", insertManyResult.InsertedIDs)
}

func Test2(t *testing.T) {
	collection, err := getConnection()
	if err != nil {
		t.Fatal(err)
	}

	filter := bson.M{}
	deleteResult, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		t.Fatal(err)
	}
	println(deleteResult.DeletedCount)
}

func Test3(t *testing.T) {
	collection, err := getConnection()
	if err != nil {
		t.Fatal(err)
	}

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetProjection(bson.M{"name": 1, "age": 1, "_id": 0})

	cur, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer cur.Close(context.Background())

	var ps []person
	if err := cur.All(context.Background(), &ps); err != nil {
		t.Fatal(err)
	}
	fmt.Println(ps)
}

func Test4(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true, // 使用完整时间戳
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			funcName := path.Base(f.Function)
			return funcName, filename + ":" + strconv.Itoa(f.Line)
		},
	})
	logrus.SetReportCaller(true)

	logrus.Info("这是一个 info 消息")
	logrus.Warn("这是一个 warning 消息")
	logrus.Error("这是一个 error 消息")

	// 设置日志级别
	logrus.SetLevel(logrus.WarnLevel)

	// 这条日志将不会显示，因为它的级别低于 WarnLevel
	logrus.Info("这条日志不会被记录")
}

func Test5(t *testing.T) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Info("这是一个 info 消息")
	logrus.Warn("这是一个 warning 消息")
	logrus.Error("这是一个 error 消息")

	// 设置日志级别
	logrus.SetLevel(logrus.WarnLevel)

	// 这条日志将不会显示，因为它的级别低于 WarnLevel
	logrus.Info("这条日志不会被记录")
}

func Test6(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
	})

	// 打开日志文件
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal("打开日志文件失败:", err)
	}
	defer file.Close()

	// 设置 logrus 的输出为文件和标准输出
	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)

	// 记录一条日志
	logrus.Info("这条日志将同时出现在控制台和文件中")
}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	frame := entry.Caller
	filename := path.Base(frame.File)

	// 自定义日志格式
	return []byte(fmt.Sprintf("[%s]|%s|%s:%d: %s\n",
		entry.Level, entry.Time.Format(time.RFC3339),
		filename, frame.Line,
		entry.Message)), nil
}

func Test7(t *testing.T) {
	logrus.SetFormatter(new(CustomFormatter))
	logrus.SetReportCaller(true)

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal("打开日志文件失败:", err)
	}
	defer file.Close()
	// 设置 logrus 的输出为文件和标准输出
	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)

	logrus.Info("这是一个 info 消息")
}

func Test8(t *testing.T) {
	// 创建一个新的日志记录器实例
	logger := logrus.New()

	// 为这个实例设置日志级别
	logger.SetLevel(logrus.WarnLevel)

	// 使用这个实例进行日志记录
	logger.Info("这条信息不会被记录，因为级别是 Warn")
	logger.Warn("这是一个警告")
}

func Test9(t *testing.T) {
	collection, err := getConnection2()
	if err != nil {
		t.Fatal(err)
	}

	filter := bson.M{"ip_address": "180.166.127.193"}
	findOptions := options.FindOne()
	findOptions.SetProjection(bson.M{
		"ip.province":  1,
		"ip.city":      1,
		"ip.os_type":   1,
		"ip.latitude":  1,
		"ip.longitude": 1,
		"_id":          0,
	})

	var result map[string]any
	err = collection.FindOne(context.Background(), filter, findOptions).Decode(&result)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		t.Error(err)
	}
	t.Log(result)
}

func Test10(t *testing.T) {
	t.Log(fmt.Sprintf("%.4f", 31.2222))
	t.Log(strconv.FormatFloat(31.2222, 'f', 4, 64))
}

func Test11(t *testing.T) {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	htmlContent := string(bs)
	doc := soup.HTMLParse(htmlContent)
	title := doc.Find("title").Text()

	t.Log(title)
}

func Test12(t *testing.T) {
	var a any
	b, ok1 := a.(map[string]any)
	c, ok2 := b["ip"]
	d, ok3 := c.(map[string]any)
	t.Log(b, ok1, c, ok2, d, ok3)
}
