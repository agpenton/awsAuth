package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

// Function to search for the credentials in the config file and convert them to variables.
func searchString(profile string) Profile {
	defer reportPanic()

	// Loading the data from the toml file.
	creds, _ := toml.LoadFile(credentialsPath)

	access := creds.Get(fmt.Sprintf("%v.aws_access_key_id", profile)).(string)
	secret := creds.Get(fmt.Sprintf("%v.aws_secret_access_key", profile)).(string)
	token := creds.Get(fmt.Sprintf("%v.aws_session_token", profile)).(string)

	// Return the values from the function.
	return Profile{
		aws_access_key_id:     access,
		aws_secret_access_key: secret,
		aws_session_token:     token,
	}
}

func credentialsFileCreation(profile string, accessTempKey string, secretTempkey string, tempToken string) {
	//log.Println("the path: " + awsDir)
	var _, err = os.Stat(credentialsPath)

	if os.IsNotExist(err) {
		file, err := os.Create(credentialsPath)
		check(err)
		defer file.Close()

		aak := fmt.Sprintf("[%v]\naws_access_key_id = %v\n", profile, accessTempKey)
		_, err = file.WriteString(aak)
		check(err)
		err = file.Sync()
		check(err)
		w := bufio.NewWriter(file)
		asak := fmt.Sprintf("aws_secret_access_key = %v\naws_session_token = %v", secretTempkey, tempToken)
		_, err = w.WriteString(asak)
		check(err)
		err = w.Flush()
		check(err)
	} else {
		log.Printf("The file %v already exists!\n", credentialsFile)
		log.Println("modifying the values")
		modifyCredentials(profile, accessTempKey, secretTempkey, tempToken)
		log.Println("done")
		return
	}

	log.Println("File created successfully", credentialsPath)
}

// Modify the credentials in the file if they exist.
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
