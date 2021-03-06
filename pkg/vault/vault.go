package vault

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// GetSecrets fills a map with the values of secrets pulled from Vault.
func GetSecrets(ctx context.Context, secretValues *map[string]map[string]string, secretNames []string) error {
	environment := os.Getenv("ENVIRONMENT")

	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.TraceLevel)
	}

	c := newVaultClient()
	c.ctx = ctx

	err := c.loadVaultEnvironment()
	if err != nil {
		return fmt.Errorf("Failed to load client environment: %v", err)
	}

	err = c.initClient()
	if err != nil {
		return fmt.Errorf("Failed to initialze client: %v", err)
	}

	if c.traceEnabled {
		var span *trace.Span
		c.ctx, span = trace.StartSpan(c.ctx, fmt.Sprintf("%s/GetSecrets", c.tracePrefix))
		defer span.End()
	}

	for _, secretName := range secretNames {
		log.Debug(fmt.Sprintf("secret= %s", secretNames))

		secret, err := c.getSecretFromVault(secretName)
		if err != nil {
			return fmt.Errorf("Error getting secret: %v", err)
		}

		(*secretValues)[secretName] = secret
	}

	return nil
}



// GetSecretVersions fills a map with the versions of secrets pulled from Vault.
func GetSecretVersions(ctx context.Context, secretVersions *map[string]int64, secretNames []string) error {
	environment := os.Getenv("ENVIRONMENT")

	if environment == "production" {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.TraceLevel)
	}

	c := newVaultClient()
	c.ctx = ctx

	err := c.loadVaultEnvironment()
	if err != nil {
		return fmt.Errorf("Failed to load client environment: %v", err)
	}

	err = c.initClient()
	if err != nil {
		return fmt.Errorf("Failed to initialze client: %v", err)
	}

	if c.traceEnabled {
		var span *trace.Span
		c.ctx, span = trace.StartSpan(c.ctx, fmt.Sprintf("%s/GetSecrets", c.tracePrefix))
		defer span.End()
	}

	for _, secretName := range secretNames {
		log.Debug(fmt.Sprintf("secret= %s", secretNames))

		secretVersion, err := c.getSecretVersionFromVault(secretName)
		if err != nil {
			return fmt.Errorf("Error getting secret version: %v", err)
		}

		(*secretVersions)[secretName] = secretVersion
	}

	return nil
}

func getEnv(varName, defaultVal string) string {

	if value, isPresent := os.LookupEnv(varName); isPresent {
		return value
	}

	return defaultVal
}

// getEncrEnvVar takes the name of an environment variable that's value begins
// with "berglas://", decrypts the value from a Google Storage Bucket with KMS,
// replaces the original environment variable value with the decrypted value,
// and returns the value as a string. If there's an error fetching the value, it
// will return an empty string along with the error message.
func getEncrEnvVar(ctx context.Context, n string) (string, error) {

	val := os.Getenv(n)
	if strings.HasPrefix(val, "berglas://") {
		if err := berglas.Replace(ctx, n); err != nil {
			return "", err
		}
	}

	return os.Getenv(n), nil
}
