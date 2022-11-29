package email

import "gopkg.in/gomail.v2"

func SendMsg(header string, msg string, email string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "i@leo.moe")
	m.SetHeader("To", "shen_jia_le@qq.com")
	// m.SetAddressHeader("Cc", "xxx@163.com", "Dan")
	m.SetHeader("Subject", header)
	m.SetBody("text/html", msg)
	// m.Attach("/home/Alex/lolcat.jpg")

	// d := gomail.NewDialer("smtp.88.com", 25, "leo876@88.com", "QN8iQKnHFIdq6yF4")
	d := gomail.NewDialer("smtp.gmail.com", 465, "i@leo.moe", "twhzfayonncaldks")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
