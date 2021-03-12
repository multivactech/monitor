package mail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"

	"github.com/multivactech/monitor/config"
)

type Mail struct {
	port    int
	host    string
	email   string
	passwd  string
	toEmail string
	header  map[string]string
}

func (mail *Mail) Init() {
	mail.host = "smtp.exmail.qq.com"
	mail.port = 465
	mail.email = config.Config.MailConfig.Sender
	mail.passwd = config.Config.MailConfig.Passwd
	mail.toEmail = config.Config.MailConfig.Receiver

	mail.header = map[string]string{
		"From":         fmt.Sprintf("测试网监控<%v>", mail.email),
		"To":           mail.toEmail,
		"Content-Type": "text/html; charset=UTF-8",
	}
}

func (mail *Mail) Send(subject, body string) {

	mail.header["Subject"] = subject
	message := ""
	for k, v := range mail.header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		mail.email,
		mail.passwd,
		mail.host,
	)
	err := sendMailUsingTLS(
		fmt.Sprintf("%s:%d", mail.host, mail.port),
		auth,
		mail.email,
		[]string{mail.toEmail},
		[]byte(message),
	)
	if err != nil {
		log.Print(err)
	}
}

//return a smtp client
func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func sendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := dial(addr)
	if err != nil {
		log.Println("Create smpt client errorInfo:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
