package vault

import (
	"log"
	"testing"

	"github.com/curiouscat2018/helloworld-api/config"
	"github.com/stretchr/testify/assert"
)

func TestMockVault_GetSecret(t *testing.T) {
	vault := NewMockVault()
	secret, err := vault.GetSecret("https://mock-url.com")
	assert.Nil(t, err)
	assert.NotEmpty(t, secret)
	log.Printf("secret value: %v", secret)
}

func TestAzureVault_GetSecret(t *testing.T) {
	if testing.Short() {
		return
	}

	vault := NewAzureVault()
	secret, err := vault.GetSecret(config.DemosecretVaultIdentifier)
	assert.Nil(t, err)
	assert.NotEmpty(t, secret)
	log.Printf("secret value: %v", secret)
}
