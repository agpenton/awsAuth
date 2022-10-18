/*
Copyright © 2022 Asdrubal Gonzalez Penton agpenton@gmail.com
*/

package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
)

func check(e error) {
	if e != nil {
		log.Println(e)
		debug.PrintStack()
	}
}

func checkPanic(e error) {
	if e != nil {
		panic(e)
		debug.PrintStack()
	}
}

func checkFatal(e error) {
	if e != nil {
		log.Fatal(e)
		debug.PrintStack()
	}
}

func reportPanic() {
	p := recover()
	if p == nil {
		return
	}
	err, ok := p.(error)
	if ok {
		fmt.Println(err)
	} else {
		panic(p)
	}
}

var homeDir, _ = os.UserHomeDir()
var awsDir = homeDir + "/.aws/"
var ssoCacheDir = awsDir + "sso/cache/"
var pwdDir = currentDir()
var credentialsFile = "credentials"
var credentialsPath = awsDir + credentialsFile
var profile string

func createSession(profile string) (*session.Session, error) {
	//sess, err := session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//	Profile:           profile,
	//})
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	}))

	//checkFatal(err)

	return sess, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a aws profile")
		return
	}

	profile := os.Args[1]

	log.Println("Login in to aws profile: ", profile)

	_, err := os.Stat(ssoCacheDir)
	if err != nil {
		println("os.Stat(): error folder name ", ssoCacheDir)
		println("and error is: ", err.Error())
		if os.IsNotExist(err) {
			ssoLogin(profile)
		}
	} else {
		if timeValidator().Before(time.Now().Local()) {
			log.Println("The credentials are Expired")
			ssoLogin(profile)
		} else {
			timeValidator().Before(time.Now().Local())
		}
	}

	log.Println("Login Success!!")

	sess, _ := createSession(profile)

	credentials, err := sess.Config.Credentials.Get()
	check(err)

	accessTempKey := credentials.AccessKeyID
	secretTempkey := credentials.SecretAccessKey
	tempToken := credentials.SessionToken

	os.Setenv("AWS_ACCESS_KEY_ID", accessTempKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretTempkey)
	os.Setenv("AWS_SESSION_TOKEN", tempToken)
	os.Setenv("AWS_PROFILE", profile)

	searchString(profile)

	if commandExists("direnv") == true {
		//envFile(profile, accessTempKey, secretTempkey, tempToken)
		envFile()
		log.Println("The temporary credentials were added to the .envrc file")
	} else {
		fmt.Println("----------------------------------")
		fmt.Println("export AWS_PROFILE=", profile)
		fmt.Println("Access Key: ", os.Getenv("AWS_ACCESS_KEY_ID"))
		fmt.Println("Secret Key: ", os.Getenv("AWS_SECRET_ACCESS_KEY"))
		fmt.Println("Session Token: ", os.Getenv("AWS_SESSION_TOKEN"))
		fmt.Println("----------------------------------")
	}

}
