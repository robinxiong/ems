package gomailer

import (
	"ems/mailer"
	"bytes"
	"github.com/jinzhu/configor"
	"gopkg.in/gomail.v2"
	"crypto/tls"
	"fmt"
	"io"
	"strings"
	"testing"
	"net/mail"
)

var Mailer *mailer.Mailer

var config = struct {
	SendRealEmail bool `env:"DEBUG"`
	Address string `env:"SMTP_Address"`
	Port int `env:"SMTP_Port"`
	User string `env:"SMTP_User"`
	Password string `env:SMTP_Password`
	DefaultTo string `env:"SMTP_To" default:"robin.xiong@live.com"`
	DefaultFrom string `env:"SMTP_From" default:"154371199@qq.com"`
}{}

var Box bytes.Buffer

func init(){
	configor.Load(&config)

	if config.SendRealEmail {
		d := gomail.NewDialer(config.Address, config.Port, config.User, config.Password)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		sender, err := d.Dial()
		if err != nil {
			panic(fmt.Sprintf("Got error %v when dail mail server: %#v", err, config))
		}

		Mailer = mailer.New(&mailer.Config{
			Sender: New(&Config{Sender: sender}),
		})
	} else {
		sender := gomail.SendFunc(func(from string, to []string, msg io.WriterTo) error {
			Box.WriteString(fmt.Sprintf("From: %v\n", from))
			Box.WriteString(fmt.Sprintf("To: %v\n", strings.Join(to, ", ")))
			_, err := msg.WriteTo(&Box)
			return err
		})

		Mailer = mailer.New(&mailer.Config{
			Sender:New(&Config{Sender: sender}),
		})
	}
}

func TestSendEmail(t *testing.T){
	Box = bytes.Buffer{}
	err := Mailer.Send(mailer.Email{
		TO:          []mail.Address{{Address: config.DefaultTo}},
		From:        &mail.Address{Address: config.DefaultFrom},
		Text:        "text email",
		HTML:        "html email <img src='cid:logo.png'/>",
		Attachments: []mailer.Attachment{{FileName: "gomail.go"}, {FileName: "../test/logo.png", Inline: true}},
	})
	if err != nil {
		t.Errorf("No error should raise when send email", err)
	}

	//fmt.Println(Box.String())
}
func TestSendEmailWithLayout(t *testing.T) {
	Box = bytes.Buffer{}

	err := Mailer.Send(
		mailer.Email{
			TO:      []mail.Address{{Address: config.DefaultTo}},
			From:    &mail.Address{Address: config.DefaultFrom},
			Subject: "hello",
		},
		mailer.Template{Name: "template", Layout: "application"},
	)

	if err != nil {
		t.Errorf("No error should raise when send email")
	}

	//fmt.Println(Box.String())
}

func TestSendEmailWithMissingLayout(t *testing.T) {
	Box = bytes.Buffer{}

	err := Mailer.Send(
		mailer.Email{
			TO:      []mail.Address{{Address: config.DefaultTo}},
			From:    &mail.Address{Address: config.DefaultFrom},
			Subject: "hello",
		},
		mailer.Template{Name: "template_with_missing_layout", Layout: "html_layout"},
	)

	if err != nil {
		t.Errorf("No error should raise when send email")
	}

	//fmt.Println(Box.String())
}
