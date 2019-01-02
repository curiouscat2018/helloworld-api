package vault

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const demoSecretUrl = "https://helloworld-vault.vault.azure.net/secrets/demosecret/ff32ea6957d04e529407cc839eef2cf8?api-version=7.0"

var integrationTest = flag.Bool("int", false, "run integration test")

func TestMockVault_GetSecret(t *testing.T) {
	vault := NewMockVault()
	secret, err := vault.GetSecret("https://mock-url.com")
	assert.Nil(t, err)
	assert.NotEmpty(t, secret)
	log.Printf("secret value: %v", secret)
}

func TestAzureVault_GetSecret(t *testing.T) {
	if !*integrationTest {
		return
	}

	vault := NewAzureVault()
	secret, err := vault.GetSecret(demoSecretUrl)
	assert.Nil(t, err)
	assert.NotEmpty(t, secret)
	log.Printf("secret value: %v", secret)
}
