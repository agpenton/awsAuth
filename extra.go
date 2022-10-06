package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func currentDir() string {
	path, err := os.Getwd()
	check(err)
	return path
}

func loadEnvrc() {
	err := godotenv.Load(".envrc")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	awsProfile := os.Getenv("AWS_PROFILE")
	//secretKey := os.Getenv("SECRET_KEY")

	log.Printf("The Profile is: %v", awsProfile)
}

func envrcVars() string {
	var envrc = []string{
		fmt.Sprintf("export AWS_PROFILE=\"%v\"", os.Getenv("AWS_PROFILE")),
		fmt.Sprintf("export AWS_ACCESS_KEY_ID=\"%v\"", os.Getenv("AWS_ACCESS_KEY_ID")),
		fmt.Sprintf("export AWS_SECRET_ACCESS_KEY=\"%v\"", os.Getenv("AWS_SECRET_ACCESS_KEY")),
		fmt.Sprintf("export AWS_SESSION_TOKEN=\"%v\"", os.Getenv("AWS_SESSION_TOKEN")),
	}
	output := strings.Join(envrc, "\n")

	return output
}

func envFile(profile string, accessTempKey string, secretTempkey string, tempToken string) {
	filename := ".envrc"
	var _, err = os.Stat(filename)

	output := envrcVars()
	//f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if _, err := f.Write([]byte(output)); err != nil {
	//	log.Fatal(err)
	//}
	//if err := f.Close(); err != nil {
	//	log.Fatal(err)
	//}
	if os.IsNotExist(err) {
		log.Printf("Creating the file %v", filename)
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkFatal(err)
		if _, err := f.Write([]byte(output)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("File already exists!", filename)
		loadEnvrc()
		modifyEnvrc()
		return
	}
	//if os.IsNotExist(err) {
	//	file, err := os.Create(pwdDir + "/" + filename)
	//	check(err)
	//	defer file.Close()
	//
	//	aak := fmt.Sprintf("[%v]\naws_access_key_id = %v\n", profile, accessTempKey)
	//	file.WriteString(aak)
	//	file.Sync()
	//	w := bufio.NewWriter(file)
	//	asak := fmt.Sprintf("aws_secret_access_key = %v\naws_session_token = %v", secretTempkey, tempToken)
	//	w.WriteString(asak)
	//	w.Flush()
	//} else {
	//	log.Println("File already exists!", filename)
	//	modifyEnvrc()
	//	return
	//}

	log.Println("File created successfully", pwdDir+"/"+filename)
}

func modifyEnvrc() {
	filename := ".envrc"
	file := pwdDir + "/" + filename
	output := envrcVars()
	//err := godotenv.Load(".envrc")
	//checkFatal(err)
	err := ioutil.WriteFile(file, []byte(output), 0644)
	checkFatal(err)
}
