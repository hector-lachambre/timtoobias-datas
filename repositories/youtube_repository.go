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

type YoutubeRepository struct {
	Cache  VideosCache
	CM     *configuration.CredentialsManager
	Client *http.Client
}

type SyncVideos struct {
	Video    *entities.Video
	LastSync time.Time
}

type VideosCache map[string]*SyncVideos

func (repository *YoutubeRepository) GetLastVideoByChannelID(ID string) (*entities.Video, error) {

	if repository.Cache == nil {
		repository.Cache = make(VideosCache)
	}

	vs, found := repository.Cache[ID]

	if found == true && time.Since(vs.LastSync).Seconds() < 60*2 {
		return vs.Video, nil
	}

	video, err := repository.getLastVideoOnYoutubeByChannelID(ID)

	if err != nil {
		return nil, err
	}

	repository.Cache[ID] = &SyncVideos{
		Video:    video,
		LastSync: time.Now(),
	}

	log.Println("Les données Youtube ont été mise à jour")

	return video, nil
}

func (repository *YoutubeRepository) getLastVideoOnYoutubeByChannelID(ID string) (*entities.Video, error) {

	credentials := repository.CM.GetCredentials()

	log.Println("Actualisation des données Youtube en cours...")

	req, err := http.NewRequest(
		"GET",
		"https://www.googleapis.com/youtube/v3/search?key="+
			credentials.Youtube.Key+
			"&channelId="+
			ID+
			"&part=snippet,id&order=date&maxResults=1",
		nil,
	)

	if err != nil {
		log.Fatal("La requête à l'API Youtube à échouée")
	}

	resp, err := repository.Client.Do(req)

	if resp.StatusCode != http.StatusOK {

		log.Printf("API Youtube status %v, echec de la mise à jour des données\n", resp.StatusCode)

		return nil,
			fmt.Errorf("API Youtube status %v, echec de la mise à jour des données", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var structuredResponse *models.YoutubeResponse

	_ = json.Unmarshal(body, &structuredResponse)

	video := &entities.Video{
		ID:          structuredResponse.Datas[0].Id.Id,
		Title:       structuredResponse.Datas[0].Snippet.Title,
		Description: structuredResponse.Datas[0].Snippet.Description,
		Date:        structuredResponse.Datas[0].Snippet.PublishedAt,
		Thumbnail:   structuredResponse.Datas[0].Snippet.Thumbnails.Default.Url,
	}

	return video, nil
}
