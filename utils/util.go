package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var basePath string = os.Getenv("APP_DIR")

// ValidateEnvs is the function that validates the environment variables needed for script execution.
func ValidateEnvs() error {
	envs := []string{"AWS_REGION", "AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY", "API_HOST", "API_SWAGGER_JSON_PATH",
		"GATEWAY_NAME", "GATEWAY_STAGE", "GATEWAY_BASE_PATH", "GATEWAY_AUTHORIZER_ENABLE"}

	isError := false

	for _, item := range envs {
		err := validateEnv(item)
		if err == true {
			isError = true
		}
	}

	if isError {
		return fmt.Errorf("Environment variables are not filled in correctly")
	}

	return nil
}

func validateEnv(env string) bool {
	if os.Getenv(env) != "" {
		return false
	}

	fmt.Printf("Environment variable %s not filled\n", env)
	return true
}

// ReturnError is the function that returns error in execution
func ReturnError(err error) {

	if err != nil {
		log.Println(err)
	}

	log.Println("Process terminated with error")
	os.Exit(1)
}

// ReturnSuccess is the function that returns success in execution
func ReturnSuccess() {

	log.Println("Finishing the run")
	os.Exit(0)
}

// Clone is the function and install the package on node
func Clone(host, branch, path string) error {
	cmd := exec.Command("git", "clone", host, "-b", branch, branch)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Dir = fmt.Sprintf("%s/%s", basePath, path)
	err := cmd.Run()
	fmt.Println(output.String())

	if err != nil {
		log.Printf("Error ao clonar a branch %s %v", branch, err)
		return err
	}

	log.Printf("Sucesso ao clonar a branch %s", branch)

	return nil
}
