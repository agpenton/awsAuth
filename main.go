/*
Copyright © 2022 Asdrubal Gonzalez Penton agpenton@gmail.com
*/

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var homeDir, _ = os.UserHomeDir()
var awsDir = homeDir + "/.aws/"

func main() {
	profile := os.Getenv("AWS_PROFILE")
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

	tempCredentials := "tempCredentials"

	f, err := os.Create(awsDir + tempCredentials)
	check(err)
	defer f.Close()

	aak := fmt.Sprintf("[%v]\naws_access_key_id = %v\n", profile, accessTempKey)
	f.WriteString(aak)
	f.Sync()
	w := bufio.NewWriter(f)
	asak := fmt.Sprintf("aws_secret_access_key = %v\naws_session_token = %v", secretTempkey, tempToken)
	w.WriteString(asak)
	w.Flush()

	fmt.Println(profile)

}
