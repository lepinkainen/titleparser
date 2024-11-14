//go:build mage
// +build mage

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	// Autoimport .env and add it to environment
	"github.com/joho/godotenv"
)

var Default = BuildLocal

const FUNCNAME = "titleparser"
const BINARYNAME = "bootstrap" // for provided.al2 runtime

func init() {
	// load .env into environment for test setup
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Vet() error {
	fmt.Println("Vet...")
	return sh.RunV("go", "vet", "./...")
}

func Test() error {
	fmt.Println("Testing...")
	mg.Deps(Vet)
	//return sh.RunV("go", "test", "-race", "-cover", "-v", "./...")
	return sh.RunV("go", "test", "-cover", "-v", "./...")
}

func Lint() error {
	fmt.Println("Linting...")
	return sh.RunV("golangci-lint", "run", "./...")
}

func Build() error {
	mg.Deps(Vet, Test, Lint)
	err := sh.RunV("go", "build", "-tags", "lambda.norpc", "-o", "build/"+BINARYNAME)
	if err != nil {
		return err
	}

	os.Chdir("build/")
	// List of Files to Zip
	files := []string{BINARYNAME}
	output := FUNCNAME + ".zip"

	if err := zipFiles(output, files); err != nil {
		return err
	}
	fmt.Println("Zipped File:", output)

	return nil
}

// Build executable file for testing and uploading
func BuildLocal() error {
	mg.Deps(Test)
	fmt.Println("Building..")
	return sh.RunV("go", "build", "-o", FUNCNAME)
}

func Publish() error {
	mg.Deps(Test, Lint, Build)

	// TODO: This doesn't work as-is, needs to use arm64 arch

	//return sh.RunV("aws lambda update-function-code", "--publish", "--function-name", FUNCNAME, "--zip-file", "fileb://build/"+FUNCNAME+".zip")
}
