package social

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Google people api Response
// {
// 	"resourceName": "",
// 	"etag": "",
// 	"names": [
// 	  {
// 		"metadata": {
// 		  "primary": true,
// 		  "source": {
// 			"type": "PROFILE",
// 			"id": ""
// 		  }
// 		},
// 		"displayName": "",
// 		"familyName": "",
// 		"givenName": "",
// 		"displayNameLastFirst": "",
// 		"unstructuredName": ""
// 	  }
// 	],
// 	"photos": [
// 	  {
// 		"metadata": {
// 		  "primary": true,
// 		  "source": {
// 			"type": "PROFILE",
// 			"id": ""
// 		  }
// 		},
// 		"url": "",
// 		"default": true
// 	  }
// 	],
// 	"emailAddresses": [
// 	  {
// 		"metadata": {
// 		  "primary": true,
// 		  "verified": true,
// 		  "source": {
// 			"type": "ACCOUNT",
// 			"id": ""
// 		  },
// 		  "sourcePrimary": true
// 		},
// 		"value": "email@gmail.com"
// 	  }
// 	]
// }

type googlePeopleApiResponse struct {
	ResourceName string `json:"resourceName"`
	Etag         string `json:"etag"`
	Names        []struct {
		Metadata struct {
			Primary bool `json:"primary"`
			Source  struct {
				Type string `json:"type"`
				Id   string `json:"id"`
			} `json:"source"`
		} `json:"metadata"`
		DisplayName          string `json:"displayname"`
		FamilyName           string `json:"familyName"`
		GivenName            string `json:"givenName"`
		DisplayNameLastFirst string `json:"displayNameLastFirst"`
		UnstructuredName     string `json:"unstructuredName"`
	} `json:"names"`
	Photos []struct {
		MetaData struct {
			Primary bool `json:"primary"`
			Source  struct {
				Type string `json:"type"`
				Id   string `json:"id"`
			} `json:"struct"`
		} `json:"metadata"`
		Url     string `json:"url"`
		Default bool   `json:"default"`
	} `json:"photos"`
	EmailAddresses []struct {
		Metadata struct {
			Primary  bool `json:"primary"`
			Verified bool `json:"verified"`
			Source   struct {
				Type string `json:"type"`
				Id   string `json:"id"`
			} `json:"struct"`
			SourcePrimary bool `json:"sourcePrimary"`
		} `json:"metadata"`
		Value string `json:"value"`
	} `json:"emailAddresses"`
}

type googlePeopleProfile struct {
	SocialID    string
	Email       string
	PhotoUrl    string
	DisplayName string
}

type googleApi struct {
	accessToken string
}

func NewGoogleApi(actk string) *googleApi {
	return &googleApi{
		accessToken: actk,
	}
}

func (g *googleApi) Do() (*googlePeopleProfile, error) {
	req, err := g.createRequest(g.accessToken)
	if err != nil {
		return nil, err
	}

	res, err := g.excuteRequest(req)
	if err != nil {
		return nil, err
	}

	gapiRes, err := g.mapTogooglePeopleApiResponse(res)
	if err != nil {
		return nil, err
	}

	return g.mapToGooglePeopleProfile(gapiRes), nil
}

func (g *googleApi) createRequest(accessToken string) (*http.Request, error) {
	url := "https://people.googleapis.com/v1/people/me?personFields=names,emailAddresses,photos"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "social GoogleApi createRequest")
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	return req, nil
}

func (g *googleApi) excuteRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "social GoogleApi excuteRequest")
	}
	return res, nil
}

func (g *googleApi) mapTogooglePeopleApiResponse(res *http.Response) (*googlePeopleApiResponse, error) {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "social GoogleApi mapTogooglePeopleApiResponse res.Body read")
	}

	var gapiRes googlePeopleApiResponse
	if err := json.Unmarshal(body, &gapiRes); err != nil {
		return nil, errors.Wrap(err, "social GoogleApi mapTogooglePeopleApiResponse Unmarshal")
	}
	return &gapiRes, nil
}

func (g *googleApi) mapToGooglePeopleProfile(gapiRes *googlePeopleApiResponse) *googlePeopleProfile {
	var socialId string
	var email string
	var photoUrl string
	var displayName string

	// must not be null
	if gapiRes.ResourceName != "" {
		replaceResourceName := strings.ReplaceAll(gapiRes.ResourceName, "people/", "")
		if replaceResourceName != "" {
			socialId = replaceResourceName
		} else {
			socialId = uuid.NewString()
		}
	} else {
		socialId = uuid.NewString()
	}

	if len(gapiRes.EmailAddresses) > 0 {
		value := gapiRes.EmailAddresses[0].Value
		if value != "" {
			email = value
		}
	}

	if len(gapiRes.Photos) > 0 {
		url := gapiRes.Photos[0].Url
		if len(url) != 0 {
			photoUrl = url
		}
	}

	if len(gapiRes.Names) > 0 {
		name := strings.Split(gapiRes.Names[0].DisplayName, " ")[0]
		if name != "" {
			displayName = name
		}
	}

	return &googlePeopleProfile{
		SocialID:    socialId,
		Email:       email,
		PhotoUrl:    photoUrl,
		DisplayName: displayName,
	}
}
