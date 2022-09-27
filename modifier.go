package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/pelletier/go-toml"
	//"github.com/BurntSushi/toml"
)

type Profile struct {
	aws_access_key_id     string
	aws_secret_access_key string
	aws_session_token     string
}

type block struct {
	Profile Profile
}

func searchString(profile string) {

	fileIO, err := os.OpenFile(credentialsPath, os.O_RDWR, 0666)
	checkPanic(err)
	defer fileIO.Close()
	log.Println(reflect.TypeOf(fileIO))
	rawBytes, err := ioutil.ReadAll(fileIO)
	log.Println(reflect.TypeOf(rawBytes))
	check(err)

	document := rawBytes

	awsProfile := block{}

	errP := toml.Unmarshal(document, &awsProfile)
	checkFatal(errP)

	fmt.Println(reflect.TypeOf(awsProfile.Profile.aws_secret_access_key))
}

func credentialsFileCreation(profile string, accessTempKey string, secretTempkey string, tempToken string) {
	log.Println("the path: " + awsDir)
	var _, err = os.Stat(credentialsPath)

	if os.IsNotExist(err) {
		file, err := os.Create(credentialsPath)
		check(err)
		defer file.Close()

		aak := fmt.Sprintf("[%v]\naws_access_key_id = %v\n", profile, accessTempKey)
		file.WriteString(aak)
		file.Sync()
		w := bufio.NewWriter(file)
		asak := fmt.Sprintf("aws_secret_access_key = %v\naws_session_token = %v", secretTempkey, tempToken)
		w.WriteString(asak)
		w.Flush()
	} else {
		log.Println("File already exists!", credentialsFile)
		log.Println("modifying the values")
		modifyCredentials(profile, accessTempKey, secretTempkey, tempToken)
		log.Println("done")
		return
	}

	log.Println("File created successfully", credentialsPath)
}

func modifyCredentials(profile string, accessTempKey string, secretTempkey string, tempToken string) {

	var tempCredentials = []string{
		fmt.Sprintf("[%v]", profile),
		fmt.Sprintf("aws_access_key_id = %v", accessTempKey),
		fmt.Sprintf("aws_secret_access_key = %v", secretTempkey),
		fmt.Sprintf("aws_session_token = %v", tempToken),
	}

	output := strings.Join(tempCredentials, "\n")
	err := ioutil.WriteFile(credentialsPath, []byte(output), 0644)
	checkFatal(err)
	return
}
