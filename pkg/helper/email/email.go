package email

import (
	"fmt"
	"gopher/internal/core"

	"gopkg.in/gomail.v2"
)

// ResetPasswordNotification for sending link to reset password
func ResetPasswordNotification(engine *core.Engine, link string, toEmail string) (err error) {

	var title, body string
	title = "Forgot your password? It happens to the best of us."
	body = "To reset your password, click the link below."

	host := engine.Environments.Email.Host
	port := engine.Environments.Email.Port
	email := engine.Environments.Email.Username
	password := engine.Environments.Email.Password

	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Reset Password")
	m.SetBody("text/html",
		fmt.Sprintf("<html><body><h4>%v</h4></body>"+
			"<pre>%v</pre><br>"+
			"<a href=\"%v\">%v</a>"+
			"</html>", title, body, link, link))

	d := gomail.NewPlainDialer(host, port, email, password)

	if err := d.DialAndSend(m); err != nil {
		engine.ServerLog.CheckError(err, "E1000166", "proccess email notification is failed")
		return err
	}

	return nil
}
