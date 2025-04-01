package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/viper"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func InitViper(runEnv string) {
	if runEnv == "local" {
		viper.AddConfigPath("configs")
		// viper.SetConfigName("config")
		viper.SetConfigName("config.local")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("ReadInConfig error: cannot read in viper config:%s", err)
		}

	} else {
		secret, err := getSecret()
		if err != nil {
			log.Fatalf("GetSecret error: cannot get secret:%s", err)
		}

		viper.SetConfigType("yaml")

		var yamlSecret = []byte(secret)
		if err := viper.ReadConfig(bytes.NewBuffer(yamlSecret)); err != nil {
			log.Fatalf("ReadConfig error: cannot read in viper config:%s", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func getSecret() (string, error) {
	projectID := os.Getenv("PROJECT_ID")
	secretID := os.Getenv("SECRET_ID")

	secretVersion := os.Getenv("SECRET_VERSION")
	if secretVersion == "" {
		secretVersion = "latest"
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
		return "", err
	}
	defer client.Close()

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{Name: fmt.Sprintf("projects/%s/secrets/%s", projectID, secretID)})
	if err != nil {
		log.Fatalf("failed to create secret: %v", err)
		return "", err
	}

	getSecretVersionReq := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secret.Name + fmt.Sprintf("/versions/%s", secretVersion),
	}

	result, err := client.AccessSecretVersion(ctx, getSecretVersionReq)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
		return "", err
	}

	return string(result.Payload.Data), nil
}

func GetSecret(projectID, secretID, secretVersion string) (string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
		return "", err
	}

	defer client.Close()

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{Name: fmt.Sprintf("projects/%s/secrets/%s", projectID, secretID)})
	if err != nil {
		log.Fatalf("failed to create secret: %v", err)
		return "", err
	}

	getSecretVersionReq := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secret.Name + fmt.Sprintf("/versions/%s", secretVersion),
	}

	result, err := client.AccessSecretVersion(ctx, getSecretVersionReq)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
		return "", err
	}

	return string(result.Payload.Data), nil
}

func ViperGetString(path string) string {
	return viper.GetString(path)
}

func ViperGetInt(path string) int {
	return viper.GetInt(path)
}

func ViperGetFloat(path string) float64 {
	return viper.GetFloat64(path)
}

func CheckLanguage(text string) bool {
	for _, char := range text {
		if !unicode.IsOneOf([]*unicode.RangeTable{unicode.Thai, unicode.Latin, unicode.Space, unicode.Number}, char) {
			return false
		}
	}
	return true
}
