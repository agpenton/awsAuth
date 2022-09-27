/*
Copyright © 2022 Asdrubal Gonzalez Penton agpenton@gmail.com
*/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
)

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}
func checkPanic(e error) {
	if e != nil {
		panic(e)
	}
}

func checkFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

var homeDir, _ = os.UserHomeDir()
var awsDir = homeDir + "/.aws/"
var pwdDir = currentDir()
var credentialsFile = "credentials"
var credentialsPath = awsDir + credentialsFile

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a aws profile")
		return
	}

	profile := os.Args[1]
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	})
	check(err)

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
		envFile(profile, accessTempKey, secretTempkey, tempToken)
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
