package email

import "gopkg.in/gomail.v2"

func SendMsg(msg string, email string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@trace.ac")
	m.SetHeader("To", email)
	// m.SetAddressHeader("Cc", "xxx@163.com", "Dan")
	m.SetHeader("Subject", "监控消息通知")
	m.SetBody("text/html", msg)
	// m.Attach("/home/Alex/lolcat.jpg")

	// d := gomail.NewDialer("smtp.88.com", 25, "leo876@88.com", "QN8iQKnHFIdq6yF4")
	// d := gomail.NewDialer("smtp.gmail.com", 465, "i@leo.moe", "twhzfayonncaldks")
	d := gomail.NewDialer("smtp.zeptomail.com.cn", 587, "emailapikey", "eiwqDPgM7DkIKgdHkn1hKum82+NpBeuc96LO4xclVfhEHubPGntIUgc0oFa+clN8ei+uT71saqF1ncLxtAR295kZb3lZuqCq+CaF7ISNMHtAL/6LeVmIwBtQgQM2b6sOWqcH8EkxBpnqNw==")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
