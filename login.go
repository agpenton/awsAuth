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

func timeValidator() time.Time {
	var expirationDate time.Time
	defer reportPanic()

	f, err := os.Open(awsDir + "sso/cache/")
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
			jsonFile, err := os.Open(v.Name())

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

			expirationDate := config.ExpiresAt
			log.Println("The expiration date is: ", expirationDate)
			return expirationDate
		}

	}

	log.Println(expirationDate)

	return expirationDate
}

func ssoLogin(profile string) {
	app := "aws"

	arg0 := "sso"
	arg1 := "login"
	arg2 := "--profile"
	arg3 := profile

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	log.Println(string(stdout))
}
