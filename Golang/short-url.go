package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "strings"

    "github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

var (
    alphabet = []byte("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
    l        = &Log{}
)

func ShortUrl(url string) string {
    md5Str := getMd5Str(url)
    var tempVal int64
    var result [4]string
    for i := 0; i < 4; i++ {
        tempSubStr := md5Str[i*8 : (i+1)*8]
        hexVal, _ := strconv.ParseInt(tempSubStr, 16, 64)
        tempVal = 0x3FFFFFFF & hexVal
        var index int64
        tempUri := []byte{}
        for i := 0; i < 6; i++ {
            index = 0x0000003D & tempVal
            tempUri = append(tempUri, alphabet[index])
            tempVal = tempVal >> 5
        }
        result[i] = string(tempUri)
    }
    return result[0]
}

func getMd5Str(str string) string {
    m := md5.New()
    m.Write([]byte(str))
    c := m.Sum(nil)
    return hex.EncodeToString(c)
}

type Log struct {
}

func (log *Log) Infof(format string, a ...interface{}) {
    log.log("INFO", format, a...)
}

func (log *Log) Info(msg string) {
    log.log("INFO", "%s", msg)
}

func (log *Log) Errorf(format string, a ...interface{}) {
    log.log("ERROR", format, a...)
}

func (log *Log) Error(msg string) {
    log.log("ERROR", "%s", msg)
}

func (log *Log) Fatalf(format string, a ...interface{}) {
    log.log("FATAL", format, a...)
}

func (log *Log) Fatal(msg string) {
    log.log("FATAL", "%s", msg)
}

func (log *Log) log(level, format string, a ...interface{}) {
    var cstSh, _ = time.LoadLocation("Asia/Shanghai")
    ft := fmt.Sprintf("%s %s %s\n", time.Now().In(cstSh).Format("2006-01-02 15:04:05"), level, format)
    fmt.Printf(ft, a...)
}

func handler(w http.ResponseWriter, r *http.Request) {
    l := &Log{}
    l.Infof("Hello world received a request, url: %s", r.URL.Path)
    l.Infof("url:%s ", r.URL)
    //if r.URL.Path == "/favicon.ico" {
    //    http.NotFound(w, r)
    //    return
    //}

    urls := strings.Split(r.URL.Path, "/")
    originUrl := getOriginUrl(urls[len(urls)-1])
    http.Redirect(w, r, originUrl, http.StatusMovedPermanently)
}

func new(w http.ResponseWriter, r *http.Request) {
    l.Infof("Hello world received a request, url: %s", r.URL)
    l.Infof("url:%s ", r.URL)
    originUrl, ok := r.URL.Query()["origin-url"]
    if !ok {
        l.Errorf("no origin-url params found")
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Bad request!"))
        return
    }

    surl := ShortUrl(originUrl[0])
    save(surl, originUrl[0])
    fmt.Fprint(w, surl)

}

func getOriginUrl(surl string) string {
    endpoint := os.Getenv("OTS_TEST_ENDPOINT")
    tableName := os.Getenv("TABLE_NAME")
    instanceName := os.Getenv("OTS_TEST_INSTANCENAME")
    accessKeyId := os.Getenv("OTS_TEST_KEYID")
    accessKeySecret := os.Getenv("OTS_TEST_SECRET")
    client := tablestore.NewClient(endpoint, instanceName, accessKeyId, accessKeySecret)

    getRowRequest := &tablestore.GetRowRequest{}
    criteria := &tablestore.SingleRowQueryCriteria{}

    putPk := &tablestore.PrimaryKey{}
    putPk.AddPrimaryKeyColumn("id", surl)
    criteria.PrimaryKey = putPk

    getRowRequest.SingleRowQueryCriteria = criteria
    getRowRequest.SingleRowQueryCriteria.TableName = tableName
    getRowRequest.SingleRowQueryCriteria.MaxVersion = 1

    getResp, _ := client.GetRow(getRowRequest)
    colmap := getResp.GetColumnMap()
    return fmt.Sprintf("%s", colmap.Columns["originUrl"][0].Value)
}

func save(surl, originUrl string) {
    endpoint := os.Getenv("OTS_TEST_ENDPOINT")
    tableName := os.Getenv("TABLE_NAME")
    instanceName := os.Getenv("OTS_TEST_INSTANCENAME")
    accessKeyId := os.Getenv("OTS_TEST_KEYID")
    accessKeySecret := os.Getenv("OTS_TEST_SECRET")
    client := tablestore.NewClient(endpoint, instanceName, accessKeyId, accessKeySecret)

    putRowRequest := &tablestore.PutRowRequest{}
    putRowChange := &tablestore.PutRowChange{}
    putRowChange.TableName = tableName

    putPk := &tablestore.PrimaryKey{}
    putPk.AddPrimaryKeyColumn("id", surl)
    putRowChange.PrimaryKey = putPk

    putRowChange.AddColumn("originUrl", originUrl)
    putRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
    putRowRequest.PutRowChange = putRowChange

    if _, err := client.PutRow(putRowRequest); err != nil {
        l.Errorf("putrow failed with error: %s", err)
    }
}

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/new", new)
    port := os.Getenv("PORT")
    if port == "" {
        port = "9090"
    }

    if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
        log.Fatalf("ListenAndServe error:%s ", err.Error())
    }

}
