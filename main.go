package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
)

var (
	apiKey string
	url    string
	image  string
	name   string
)

func main() {

	flag.StringVar(&apiKey, "api_key", "", "1panel API key")
	flag.StringVar(&url, "url", "", "1panel API URL http://192.168.1.22:21733")
	flag.StringVar(&image, "image", "", "docker image name")
	flag.StringVar(&name, "name", "", "docker container name")
	flag.Parse()
	println(apiKey, url, image, name)
	if len(apiKey) == 0 || len(url) == 0 || len(image) == 0 || len(name) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	token := generateToken(apiKey)
	println(token)
	client := resty.New()
	client.SetHeader("1Panel-Token", token)
	client.SetHeader("1Panel-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	resp, err := client.R().SetBody(map[string]interface{}{
		"forcePull": true,
		"image":     image,
		"name":      name,
	}).Post(url + "/api/v1/containers/upgrade")
	if err != nil {
		panic(err)
	}
	type Resp struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
	var r Resp
	err = sonic.Unmarshal(resp.Body(), &r)
	if err != nil {
		panic(err)
	}
	if r.Code != 200 {
		panic(fmt.Sprintf("error: %s", r.Message))
	}
	fmt.Println(resp.String())
}

func generateToken(apiKey string) string {
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("1panel%s%d", apiKey, timestamp)
	return md5Sum(data)
}

func md5Sum(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
