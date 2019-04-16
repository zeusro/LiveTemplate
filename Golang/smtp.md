
## 使用 Go 语言通过 SMTP 发送邮件，go语言的版本为1.9


```go
package main

import (
    "fmt"
    "net/smtp"
    "strings"
)

func SendToMail(user, password, host, subject, body, mailtype, replyToAddress string, to, cc, bcc []string) error {
    hp := strings.Split(host, ":")
    auth := smtp.PlainAuth("", user, password, hp[0])
    var content_type string
    if mailtype == "html" {
        content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
    } else {
        content_type = "Content-Type: text/plain" + "; charset=UTF-8"
    }

    cc_address := strings.Join(cc, ";")
    bcc_address := strings.Join(bcc, ";")
    to_address := strings.Join(to, ";")
    msg := []byte("To: " + to_address + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\nReply-To: " + replyToAddress + "\r\nCc: " + cc_address + "\r\nBcc: " + bcc_address + "\r\n" + content_type + "\r\n\r\n" + body)

    send_to := MergeSlice(to, cc)
    send_to = MergeSlice(send_to, bcc)
    err := smtp.SendMail(host, auth, user, send_to, msg)
    return err
}

func main() {
    user := "控制台创建的发信地址"
    password := "控制台设置的SMTP密码"
    host := "smtpdm.aliyun.com:25"
    to := []string{"收件人地址","收件人地址1"}
    cc := []string{"抄送地址","抄送地址1"}
    bcc := []string{"密送地址","密送地址1"}

    subject := "test Golang to sendmail"
    mailtype :="html"
    replyToAddress:="***@xxx.com"

    body := `
        <html>
        <body>
        <h3>
        "Test send to email"
        </h3>
        </body>
        </html>
        `
    fmt.Println("send email")
    err := SendToMail(user, password, host, subject, body, mailtype, replyToAddress, to, cc, bcc)
    if err != nil {
        fmt.Println("Send mail error!")
        fmt.Println(err)
    } else {
        fmt.Println("Send mail success!")
    }

}

func MergeSlice(s1 []string, s2 []string) []string {
    slice := make([]string, len(s1)+len(s2))
    copy(slice, s1)
    copy(slice[len(s1):], s2)
    return slice
}

```

## 如果go语言的版本为1.9.2，出现错误提示:“unencrypted connection”，因为此版本需要加密认证，可采用LOGIN认证，则需要增加以下内容：


```go
type loginAuth struct {
    username, password string
}

func LoginAuth(username, password string) smtp.Auth {
    return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
    // return "LOGIN", []byte{}, nil
    return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
    if more {
        switch string(fromServer) {
        case "Username:":
            return []byte(a.username), nil
        case "Password:":
            return []byte(a.password), nil
        }
    }
    return nil, nil
}

```



auth := LoginAuth(user, password)

## ssl

```go
package main

import (
    "crypto/tls"
    "fmt"
    "log"
    "net"
    "net/smtp"
)

func main() {
    host := "smtpdm.aliyun.com"
    port := 465
    email := "xxx@xxx.com"
    password := "TExxxxxst"
    toEmail := "xxx@xxxxx.com"

    header := make(map[string]string)
    header["From"] = "test" + "<" + email + ">"
    header["To"] = toEmail
    header["Subject"] = "邮件标题"
    header["Content-Type"] = "text/html; charset=UTF-8"

    body := "我是一封测试电子邮件!"

    message := ""
    for k, v := range header {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body

    auth := smtp.PlainAuth(
        "",
        email,
        password,
        host,
    )

    err := SendMailUsingTLS(
        fmt.Sprintf("%s:%d", host, port),
        auth,
        email,
        []string{toEmail},
        []byte(message),
    )

    if err != nil {
        panic(err)
    } else {
        fmt.Println("Send mail success!")
    }
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
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
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
    to []string, msg []byte) (err error) {

    //create smtp client
    c, err := Dial(addr)
    if err != nil {
        log.Println("Create smpt client error:", err)
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

```

参考:
[SMTP 之 Go 调用示例](https://help.aliyun.com/document_detail/29457.html?spm=a2c4g.11186623.6.613.1bf92649pCkcJh#h2--go-1-9-2-unencrypted-connection-login-2)
