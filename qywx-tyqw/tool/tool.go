package tool

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

// 判空.
func IsEmpty(target string) bool {
	if len(target) == 0 {
		return true
	}
	if len(strings.TrimSpace(target)) == 0 {
		return true
	}
	return false
}

// ============== accessToken ==========================
func GetAccessToken(corpid string, corpsecret string) (string, error) {
	log.Println("准备获取 accesstoken")
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	if IsEmpty(corpid) || IsEmpty(corpsecret) {
		log.Println("参数异常,", "corpid : ", corpid, "corpsecret: ", corpsecret)
		return "", errors.New("参数异常,请查看后台日志.")
	}
	url = url + "?corpid=" + corpid + "&corpsecret=" + corpsecret
	tokenMap, err := HttpGetRequest(url)
	if err != nil {
		log.Println("获取token异常:", err)
		return "", errors.New("获取token异常,查看服务日志")
	}
	jsonData, _ := json.Marshal(tokenMap)
	var tokenResp AccessTokenResp
	err = json.Unmarshal(jsonData, &tokenResp)
	if err != nil {
		log.Println("access_token解析异常:", err)
		return "", errors.New("token解析异常,查看解析日志")
	}
	if tokenResp.ErrCode != 0 {
		log.Println("token结果异常:", tokenResp)
		return "", errors.New("获取token返回状态码不正确")
	}
	return tokenResp.AccessToken, nil
}

type AccessTokenResp struct {
	ErrCode     int32  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int32  `json:"expires_in"`
}

type WxServerIpList struct {
	IpList  []string `json:"ip_list"`
	ErrCode int32    `json:"errcode"`
	ErrMsg  string   `json:"errmsg"`
}

// IsAccessTokenExpire 判断acesstoken是否过期
// 只有当没有过期,返回false.,其他都返回true
func IsAccessTokenExpire(accesstoken string) bool {
	// 通过调用企业微信的 指定接口.
	log.Println("准备校验 accesstoken")
	url := "https://qyapi.weixin.qq.com/cgi-bin/get_api_domain_ip?access_token="
	if IsEmpty(accesstoken) {
		log.Println("accesstoken 未赋值,接口调用异常")
		return true
	}
	url = url + accesstoken
	result, err := HttpGetRequest(url)
	if err != nil {
		log.Println("校验accessToken 异常,", err)
		return true
	}
	var wxIp WxServerIpList
	jsonData, _ := json.Marshal(result)
	err = json.Unmarshal(jsonData, &wxIp)
	if err != nil {
		log.Println("Unmarshal Error...", err)
		return true
	}
	// ErrCode 为 42001 时,表示过期.
	if wxIp.ErrCode == 0 {
		return false
	}
	return true
}

// 校验,并重新获取token.建议每次方法调用前都执行
func CheckAndGetAccessToken(token string, corpid string, corpsecret string) string {
	if IsEmpty(token) || IsAccessTokenExpire(token) {
		log.Println("token过期,需要重新获取")
		tokenNew, err := GetAccessToken(corpid, corpsecret)
		if err != nil {
			log.Println("获取token异常", err)
		}
		return tokenNew
	}
	return token
}

// ====================== http =======================================
// client 参数
var client = http.Client{
	Timeout: 10 * time.Second,
}

// HttpGetRequest 因为返回的是json,先用map保存.
// 针对的是 需要返回json数据的get 接口.
func HttpGetRequest(url string) (map[string]interface{}, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Println("返回状态不正常: ", resp)
		return nil, errors.New("接口响应不正常.")
	}
	defer resp.Body.Close()
	result := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("解码异常: ", resp)
		return nil, err
	}
	//decoder := json.NewDecoder(resp.Body)
	//err = decoder.Decode(&result)
	return result, nil
}

// HttpPostJson POST请求.获取json结果,结果以map形式返回.
// 不做逻辑处理,只是做一个post请求而已.
// 获取到结果后,先判 err . 不为空表示获取到了返回值.再做逻辑处理.
func HttpPostJson(url string, data interface{}, header map[string]string) (map[string]interface{}, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		log.Println("JOSN编码失败", err)
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		log.Println("创建请求失败", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	// 添加请求头
	if header != nil {
		for key, val := range header {
			request.Header.Add(key, val)
		}
	}
	log.Println("Post的请求体: \n", request)

	// 如果直接使用 http.Post() ,是无法添加Header
	// 在需要添加修改Header的地方,建议使用client.Do()
	response, err := client.Do(request)
	if err != nil {
		log.Println("接口返回Response: ", response)
		log.Println("发送请求失败", err)
		return nil, err
	}
	//log.Println("接口返回Response: ", response)
	defer response.Body.Close()
	result := make(map[string]interface{})
	err = json.NewDecoder(response.Body).Decode(&result)
	log.Println("Post请求返回的Result: \n", result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 判断ip是否在白名单中.
func IsWhiteIp(addr string, ips []string) bool {
	for _, ip := range ips {
		if addr == ip {
			return true
		}
	}
	return false
}
