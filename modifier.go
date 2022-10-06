package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	//"github.com/spf13/viper"
	"github.com/pelletier/go-toml"
)

type Profile struct {
	aws_access_key_id     string
	aws_secret_access_key string
	aws_session_token     string
}

type block struct {
	Profile Profile
}

func searchString(profile string) Profile {
	defer reportPanic()

	creds, _ := toml.LoadFile(credentialsPath)

	access := creds.Get(fmt.Sprintf("%v.aws_access_key_id", profile)).(string)
	secret := creds.Get(fmt.Sprintf("%v.aws_secret_access_key", profile)).(string)
	token := creds.Get(fmt.Sprintf("%v.aws_session_token", profile)).(string)

	fmt.Println(access)
	fmt.Println(secret)
	fmt.Println(token)

	//fileIO, err := os.OpenFile(credentialsPath, os.O_RDWR, 0666)
	//checkPanic(err)
	//defer fileIO.Close()
	//rawBytes, err := ioutil.ReadAll(fileIO)
	//check(err)
	//
	//document := rawBytes
	////fmt.Println(string(document))
	//
	////awsProfile := block{}
	//awsProfile := Profile{}
	//
	//errP := toml.Unmarshal(document, &awsProfile)
	//checkFatal(errP)
	//
	//fmt.Println(awsProfile.aws_access_key_id)
	return Profile{
		aws_access_key_id:     access,
		aws_secret_access_key: secret,
		aws_session_token:     token,
	}
}

//func searchString(profile string) {
//	viper.SetConfigName(credentialsFile) // name of config file (without extension)
//	viper.AddConfigPath(".")             // optionally look for config in the working directory
//	err := viper.ReadInConfig()          // Find and read the config file
//
//	if err != nil { // Handle errors reading the config file
//		panic(fmt.Errorf("Fatal error config file: %s \n", err))
//	}
//	fmt.Println(err)
//
//	//fmt.Println("Access Key", viper.GetString("profile")
//	//fmt.Println("database user", viper.GetString("database.user")
//}

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
		fmt.Sprintf("aws_access_key_id = \"%v\"", accessTempKey),
		fmt.Sprintf("aws_secret_access_key = \"%v\"", secretTempkey),
		fmt.Sprintf("aws_session_token = \"%v\"", tempToken),
	}

	output := strings.Join(tempCredentials, "\n")
	err := ioutil.WriteFile(credentialsPath, []byte(output), 0644)
	checkFatal(err)
	return
}
