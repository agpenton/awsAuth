package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

// Function to check if command exist
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Function to get the current directory.
func currentDir() string {
	path, err := os.Getwd()
	check(err)
	return path
}

// Loading the data from .envrc file.
func loadEnvrc() {
	err := godotenv.Load(".envrc")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	awsProfile := os.Getenv("AWS_PROFILE")

	log.Printf("The Profile is: %v", awsProfile)
}

// Export the variables to the envrc file.
func envrcVars() string {
	var envrc = []string{
		fmt.Sprintf("export AWS_PROFILE=\"%v\"", os.Getenv("AWS_PROFILE")),
		fmt.Sprintf("export AWS_ACCESS_KEY_ID=\"%v\"", os.Getenv("AWS_ACCESS_KEY_ID")),
		fmt.Sprintf("export AWS_SECRET_ACCESS_KEY=\"%v\"", os.Getenv("AWS_SECRET_ACCESS_KEY")),
		fmt.Sprintf("export AWS_SESSION_TOKEN=\"%v\"", os.Getenv("AWS_SESSION_TOKEN")),
		fmt.Sprintf("export AWS_REGION=\"%v\"", os.Getenv("AWS_REGION")),
	}
	output := strings.Join(envrc, "\n")

	return output
}

// Writing the data inside the .envrc file
func envFile() {
	filename := ".envrc"
	var _, err = os.Stat(filename)

	output := envrcVars()
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
		log.Printf("The file %s already exists!\n", filename)
		loadEnvrc()
		modifyEnvrc()
		Shellout("direnv allow")
		return
	}

	log.Println("File created successfully", pwdDir+"/"+filename)
}

// Modify the file if exist.
func modifyEnvrc() {
	filename := ".envrc"
	file := pwdDir + "/" + filename
	output := envrcVars()
	err := ioutil.WriteFile(file, []byte(output), 0644)
	checkFatal(err)
}

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
