package repositories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gitlab.com/timtoobias-projects/timtoobias-core/entities"
	"gitlab.com/timtoobias-projects/timtoobias-datas/configuration"
	"gitlab.com/timtoobias-projects/timtoobias-datas/models"
)

// TwitchRepository is a StreamingRepository implementation
type TwitchRepository struct {
	CM       *configuration.CredentialsManager
	Client   *http.Client
	cache    *entities.StreamingStatus
	lastSync time.Time
}

// GetStreamingStatusByID return the cached streaming status
func (repository *TwitchRepository) GetStreamingStatusByID(ID string) (*entities.StreamingStatus, error) {

	if time.Since(repository.lastSync).Seconds() < 30 {

		return repository.cache, nil
	}

	err := repository.updateStreamingStatus(ID)

	if err != nil {
		return nil, err
	}

	return repository.cache, nil
}

// GetLastSync return the last time where datas was syncronized with Twitch
func (repository *TwitchRepository) GetLastSync() *time.Time {

	return &repository.lastSync
}

func (repository *TwitchRepository) getTwitchBearerToken(credentials *configuration.TwitchCredentials) (string, error) {

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"https://id.twitch.tv/oauth2/token?client_id=%v&client_secret=%v&grant_type=client_credentials",
			credentials.Client,
			credentials.Secret,
		),
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := repository.Client.Do(req)

	if err != nil {
		return "", nil
	}

	var data map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &data)

	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("%v", data["access_token"]), nil

}

func (repository *TwitchRepository) updateStreamingStatus(ID string) error {

	credentials := repository.CM.GetCredentials().Twitch

	bearer, err := repository.getTwitchBearerToken(credentials)

	if err != nil {
		log.Printf("L'authentification Twitch à échouée: %v\n", err)

		return err
	}

	log.Println("Actualisation des données Twitch en cours...")
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_id="+ID, nil)

	if err != nil {
		log.Println("La requête à l'API distante à échouée")

		return err
	}

	req.Header.Add("Authorization", "Bearer "+bearer)
	req.Header.Add("Client-ID", credentials.Client)

	httpResponse, err := repository.Client.Do(req)

	if httpResponse.StatusCode != http.StatusOK {

		log.Printf("API status %v, echec de la mise à jour des données\n", httpResponse.StatusCode)

		return err
	}

	body, _ := ioutil.ReadAll(httpResponse.Body)

	response := &models.TwitchResponse{}

	_ = json.Unmarshal(body, response)

	if len(response.Datas) != 0 {

		repository.cache = &entities.StreamingStatus{
			Title: response.Datas[0].Title,
			Date:  response.Datas[0].StartedAt,
		}
	} else {
		repository.cache = nil
	}

	repository.lastSync = time.Now()

	log.Println("Les données Twitch ont été mise à jour")

	return nil
}
