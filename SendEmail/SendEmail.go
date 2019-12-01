package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

func main() {
	TimeSettle()

}
func getPm25() int{

	client := &http.Client{}
	resp, err := client.Get("https://api.waqi.info/feed/shanghai/?token=demo")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	m := make(map[string]interface{})
	json.Unmarshal([]byte(body), &m)
	a :=m["data"].(map[string]interface{})["iaqi"].(map[string]interface{})["pm25"].(map[string]interface{})["v"]
	//fmt.Println(m["data"].(map[string]interface{})["iaqi"].(map[string]interface{})["pm25"].(map[string]interface{})["v"])
	//fmt.Println(reflect.TypeOf(m))
	b :=int(a.(float64))
	return b
}

//发送邮件
func sendToMail(user,password,host,to,subject,body,mailtype string) error {
	hp := strings.Split(host,":")
	auth := smtp.PlainAuth("",user,password,hp[0])
	var content_type string
	if mailtype =="html" {
		content_type = "Content_Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content_Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To:" + to +"\r\nFrom: " + user + "<"+
		user + ">\r\nSubject: "+ subject + "\r\n" +
		content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to,";")
	err := smtp.SendMail(host,auth,user,send_to,msg)
	return err
}

func sendEmail(subject,body string)  {
	user := "tobey1024@126.com"
	pwd := "wmt126yxsqm"
	host := "smtp.126.com:25"
	to := "1832541104@qq.com"//可以用;隔开发送多个
	fmt.Println("send email")
	err := sendToMail(user,pwd,host,to,subject,body,"html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}
// 定时结算（一天发一次）
func TimeSettle() {
	d := time.Duration(time.Minute)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		//currentTime := time.Now()
		pm25 :=getPm25()
		if pm25>50 {
		//	if currentTime.Hour() == 8 { // 8点发送
				Sendinfo(pm25)
				time.Sleep(time.Hour)
			//}
		}
		<-t.C
	}
}

func Sendinfo(pm25 int) {
	subject, body := "空气质量预警","今日pm25高达"+  strconv.Itoa(pm25) +",请尽量减少户外活动"

	sendEmail(subject, body)
}