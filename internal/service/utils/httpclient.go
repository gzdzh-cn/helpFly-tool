package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const apiEncrypt = "gofly@888"

// http请求-通用
// method=请求方式(GET POST PUT DELETE)，url=请求地址，data=请求数据,headerList=请求头数据
func RequestHttp(method, url, data string, headerList map[string]string) (string, error) {
	payload := strings.NewReader(data)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", err
	}
	for key, val := range headerList {
		req.Header.Add(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// 发送GET请求
func HttpGet(url_text string, data map[string]interface{}) (map[string]interface{}, error) {
	u, err := url.Parse(url_text)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	for k, v := range data {
		paras.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = paras.Encode()
	resp, err := http.Get(u.String())
	//关闭资源
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("request token err :" + err.Error())
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New("request token response json parse err :" + err.Error())
	} else {
		return jMap, nil
	}

}

// 发送POST请求
func HttpPost(url_text string, urldata map[string]interface{}, postdata map[string]interface{}, contentType string) (map[string]interface{}, error) {
	u, err := url.Parse(url_text)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	for k, v := range urldata {
		paras.Set(k, v.(string))
	}
	u.RawQuery = paras.Encode()
	//json序列化
	jsonData := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(jsonData)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(postdata); err != nil {
		return nil, errors.New("请求错误 :" + err.Error())
	}
	body := bytes.NewBufferString(string(jsonData.Bytes()))
	resp, erro := http.Post(u.String(), contentType, body)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if erro != nil {
		return nil, errors.New("请求错误 :" + erro.Error())
	}
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New(" 返回结果解析错误 :" + err.Error())
	} else {
		return jMap, nil
	}
}

// tool
// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func Get(url string) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return result.String()
}

/**
* 发送POST请求
* url：         请求地址
* data：        POST请求提交的数据 interface{}
* contentType： 请求体格式，如：application/json
* content：     请求放回的内容
 */
func Post(url string, data interface{}, contentType string) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)
	return string(result)
}

//****************与GoFly后端api对接-请求***********************

/**
*  post-api 请求
*  添加请求头验证
*  data map[string]interface{}
 */
func PostApi(url string, postdata map[string]interface{}, contentType string) (map[string]interface{}, error) {
	if contentType == "" {
		contentType = "application/json"
	}
	//json序列化
	jsonData := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(jsonData)
	jsonEncoder.SetEscapeHTML(false)
	if err := jsonEncoder.Encode(postdata); err != nil {
		return nil, errors.New("请求错误 :" + err.Error())
	}
	body := bytes.NewBufferString(string(jsonData.Bytes()))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("content-type", contentType)
	//处理api验证
	timestamp := time.Now().Unix()
	passstr := Md5(fmt.Sprintf("%v%v", apiEncrypt, timestamp))
	apiencoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v#%v", passstr, timestamp)))
	req.Header.Add("apiverify", apiencoded)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New(" 返回结果解析错误 :" + err.Error())
	} else {
		return jMap, nil
	}
}

/**
*  get-api请求
*  添加请求头验证
*  data map[string]interface{}
 */
func GetApi(urltext string, data map[string]interface{}, contentType string) (map[string]interface{}, error) {
	if contentType == "" {
		contentType = "application/json"
	}
	u, err := url.Parse(urltext)
	if err != nil {
		log.Fatal(err)
	}
	paras := &url.Values{}
	//设置请求参数
	for k, v := range data {
		paras.Set(k, fmt.Sprintf("%v", v))
	}
	u.RawQuery = paras.Encode()
	body := bytes.NewBufferString("")
	req, err := http.NewRequest("GET", u.String(), body)
	req.Header.Add("content-type", contentType)
	//处理api验证
	timestamp := time.Now().Unix()
	passstr := Md5(fmt.Sprintf("%v%v", apiEncrypt, timestamp))
	apiencoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v#%v", passstr, timestamp)))
	req.Header.Add("apiverify", apiencoded)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	jMap := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jMap)
	if err != nil {
		return nil, errors.New(" 返回结果解析错误 :" + err.Error())
	} else {
		return jMap, nil
	}
}
