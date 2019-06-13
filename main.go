/**
Author:       yuyongpeng@hotmail.com
Github:       https://github.com/yuyongpeng/
Date:         2019-06-13 23:20:27
LastEditors:
LastEditTime: 2019-06-13 23:20:27
Description:
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"strconv"
	"strings"
)

/**
获得验证码
*/
func fetchSec(mobile string) bool {
	params := map[string]string{
		"mobile": mobile,
		"smsType": "sms",
	}
	url := "http://api-release-planet.hard-chain.cn/ucenter/account/captcha"
	request := gorequest.New()
	resp, body, errs := request.Post(url).
		Type("multipart").
		SendMap(params).
		End()
	if errs != nil {
		fmt.Println(errs)
	}
	fmt.Println(resp.Status)
	var ret map[string]interface{}
	if err := json.Unmarshal([]byte(body), &ret); err != nil {
		fmt.Println(err)
	}
	if ret["code"].(float64) == 200 {
		return true
	}else{
		return false
	}
}

/*
登录
*/
func login(mobile, sec string) (uid float64, token string , e error) {
	url := "http://api-release-planet.hard-chain.cn/ucenter/account/login"
	body := map[string]string {
		"captcha": sec,
		"mobile": mobile,
		"app_id": "10000",
	}
	request := gorequest.New()
	resp, resbody, errs := request.Post(url).SendMap(body).End()
	if errs != nil {
		fmt.Println(errs)
	}
	fmt.Println(resp.Status)
	fmt.Println(resbody)
	var ret map[string]interface{}
	if err := json.Unmarshal([]byte(resbody), &ret); err != nil {
		fmt.Println(err)
	}
	if ret["code"].(float64) == 200 {
		data := ret["data"].(map[string]interface{})
		id := data["id"].(float64)
		token := data["ucenter_token"].(string)
		return id, token, nil
	}else{
		return 0,"" , fmt.Errorf("登录失败")
	}
}
/**
抓取碎片
 */
func fetchStar(uid , token string) {
	url := "http://api-release-planet.hard-chain.cn/fans/mine/gatherDebris"
	request := gorequest.New()
	resp, body, errs := request.Get(url).
		Set("user-id", uid).
		Set("ucenter-token", token).
		End()
	if errs != nil {
		fmt.Print(errs)
	}
	fmt.Println(resp.Status)
	fmt.Println(body)
}

func main() {
	//参考 version 2 https://my.oschina.net/zengsai/blog/3719
	//login("13552885937", "666666")
	//os.Exit(1)
	f := bufio.NewReader(os.Stdin) //读取输入的内容
	fmt.Print("请输入登录用的手机号>")
	var Input string
	mobile, _ := f.ReadString('\n') //定义一行输入的内容分隔符。
	if fetchSec(mobile) == true {
		fmt.Println("发送成功，请查收")
	}
	fmt.Print("请查看手机获取验证码，并输入>")
	sec, _ := f.ReadString('\n')
	id, token , err := login(mobile, sec)
	if 	err != nil {
		fmt.Println(err)
	}

	fetchStar(strconv.FormatFloat(id, 'E', -1, 64), token)


	fmt.Println("字符串长度", len(Input))
	for i := 0; i < len(Input); i++ {
		if i >= len(Input)-2 { //最后一个字符,输出数字
			fmt.Print(Input[i])
		} else {
			fmt.Print(string(Input[i]))
		}
	}
	//windows平台操作
	//分隔符'\n' Input是 xxx\r\n    编码是1310
	//分隔符'\r' Input是 xxx\r  编码是13
	Input = strings.Replace(Input, "\n", "", -1)
	Input = strings.Replace(Input, "\r", "", -1)
	fmt.Println("")
	fmt.Println("字符串长度", len(Input))
	for i := 0; i < len(Input); i++ {
		fmt.Print(string(Input[i]))
	}
}