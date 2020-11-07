package configuration

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"gitlab.com/timtoobias-projects/timtoobias-datas/errors"
	"gopkg.in/ini.v1"
)

// CredentialsManager manage configuration loading
type CredentialsManager struct {
	Internal *Credentials
}

// GetCredentials load configuration for the first call and return it
func (cm *CredentialsManager) GetCredentials() *Credentials {
	if cm.Internal != nil {
		return cm.Internal
	}

	credentials, err := loadCredentials("config.ini")

	if err != nil {
		log.Fatalf("Une erreur est survenue lors du chargement de la configuration: %v", err)
		return nil
	}

	cm.Internal = credentials

	return credentials
}

// Credentials needed by repositories
type Credentials struct {
	Twitch  *TwitchCredentials
	Youtube *YoutubeCredentials
}

// TwitchCredentials contains twich keys to request API
type TwitchCredentials struct {
	Client string
	Secret string
}

// YoutubeCredentials contains youtube key to request API
type YoutubeCredentials struct {
	Key string
}

func loadCredentials(filename string) (*Credentials, error) {

	rootPath, err := os.Executable()

	if err != nil {

		return nil, &errors.CannotLoadConfigurationError{
			Message: "Can not get information about the runtime environment",
		}
	}

	configPath := filepath.Join(path.Dir(rootPath), filename)

	cfg, err := ini.Load(configPath)

	if err != nil {

		return nil, &errors.CannotLoadConfigurationError{
			Message: fmt.Sprintf("Fail to read configuration file: %v", err),
		}
	}

	return &Credentials{
			Twitch: &TwitchCredentials{
				Client: cfg.Section("Twitch").Key("client").String(),
				Secret: cfg.Section("Twitch").Key("secret").String(),
			},
			Youtube: &YoutubeCredentials{
				Key: cfg.Section("Youtube").Key("key").String(),
			},
		},
		nil
}
