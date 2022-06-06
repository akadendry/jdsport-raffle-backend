package controllers

import (
	// "fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"gopkg.in/gomail.v2"
)

func sendMail(to []string, cc []string, subject, message string, attach string) error {
  // fmt.Println("test : ",to)
    mailer := gomail.NewMessage()
    mailer.SetHeader("From", os.Getenv("CONFIG_SENDER_NAME"))
    mailer.SetHeader("To", strings.Join(to, ","))
    // if(strings.Join(cc, ",") == ""){
    //     mailer.SetAddressHeader("Cc", strings.Join(cc, ","), "")
    // }
    mailer.SetHeader("Subject", subject)
    mailer.SetBody("text/html", message)
    // if(attach != ""){
    //    mailer.Attach(attach)
    // }
   

	CONFIG_SMTP_PORT, _ := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))

    dialer := gomail.NewDialer(
        os.Getenv("CONFIG_SMTP_HOST"),
        CONFIG_SMTP_PORT,
        os.Getenv("CONFIG_AUTH_EMAIL"),
        os.Getenv("CONFIG_AUTH_PASSWORD"),
    )

    err := dialer.DialAndSend(mailer)
    if err != nil {
        log.Fatal(err.Error())
    }
    if err != nil {
        return err
    }

    // log.Println("Mail sent!")
    return nil
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

func StringRandom(length int) string {
  return StringWithCharset(length, charset)
}