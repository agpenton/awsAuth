package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	StartUrl    string    `json:"startUrl"`
	Region      string    `json:"region"`
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

//func dirCreation() {
//	directoryName := awsDir + "sso/cache/"
//	_, err := os.Stat(directoryName)
//	if err != nil {
//		println("os.Stat(): error folder name ", directoryName)
//		println("and error is: ", err.Error())
//		if os.IsNotExist(err) {
//			ssoLogin(profile)
//		}
//	} else {
//		println("os.Stat(): No Error for folderName : ", directoryName)
//	}
//
//}

// Validating the expiration date of the session.
func timeValidator() time.Time {
	var expirationDate time.Time
	defer reportPanic()

	f, err := os.Open(ssoCacheDir)
	if err != nil {
		fmt.Println(err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range files {
		if v.Name() != "botocore-client-id-eu-central-1.json" {
			// Open our jsonFile
			jsonFile, err := os.Open(ssoCacheDir + "/" + v.Name())

			// if we os.Open returns an error then handle it
			if err != nil {
				fmt.Println(err)
			}

			defer jsonFile.Close()
			log.Println("Successfully Opened ", v.Name())

			byteValue, _ := ioutil.ReadAll(jsonFile)
			var config Config

			err = json.Unmarshal(byteValue, &config)
			check(err)

			expirationDate = config.ExpiresAt
		}

	}

	log.Println(expirationDate.Local())

	return expirationDate.Local()
}

// login command for the aws sso if the session is expired.
func ssoLogin(profile string) string {
	app := "aws"

	arg0 := "sso"
	arg1 := "login"
	arg2 := "--profile"
	arg3 := profile

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	// Print the output
	log.Println(string(stdout))

	return string(stdout)
}
