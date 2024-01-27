package tauclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/modaniru/tau-client/entities"
)


type tauClient struct{
	client *http.Client
	ref string
}

func NewTauClient(ref string) (*tauClient, error){
	tClient := &tauClient{client: http.DefaultClient, ref: ref}
	err := tClient.Ping()
	if err != nil{
		return nil, err
	}
	return tClient, nil
}

func (t *tauClient) Ping() error{
	response, err := t.client.Get(fmt.Sprintf("%s/ping", t.ref))
	if err != nil{
		return err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil{
		return err
	}

	if response.StatusCode != 200{
		return errors.New(string(bytes))
	}

	return nil
}

type accessToken struct{
	Token string `json:"token"`
}

func (t *tauClient) SignIn(twitchToken string) (*entities.Token, error){
	accessToken := accessToken{Token: twitchToken}
	request, err := json.Marshal(accessToken)
	if err != nil{
		return nil, err
	}

	response, err := t.client.Post(fmt.Sprintf("%s/sign-in", t.ref), "json", bytes.NewReader(request))

	if err != nil{
		return nil, err
	}

	resp, err := io.ReadAll(response.Body)
	if err != nil{
		return nil, err
	}

	//TODO crete custom errors struct
	if response.StatusCode != 200{
		return nil, errors.New(string(resp))
	}

	token := entities.Token{}
	err = json.Unmarshal(resp, &token)
	if err != nil{
		return nil, err
	}

	return &token, nil
}

func (t *tauClient) GetUser(token *entities.Token) (*entities.User, error){
	req, err := http.NewRequest(http.MethodGet,fmt.Sprintf("%s/api/user", t.ref), nil)
	if err != nil{
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer " + token.Jwt)
	response, err := t.client.Do(req)
	if err != nil{
		return nil, err
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil{
		return nil, err
	}
//TODO custom errors
	if response.StatusCode != 200{
		return nil, errors.New(string(bytes))
	}

	user := entities.User{}
	err = json.Unmarshal(bytes, &user)
	if err != nil{
		return nil, err
	}

	return &user, nil
}