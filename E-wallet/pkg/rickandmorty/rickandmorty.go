package rickandmorty

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Rate struct {
	log    *logrus.Entry
	xrHost string
	apiKey string
}

type Resp struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	AirDate string    `json:"air_date"` 
	Episode string    `json:"episode"`
	Characters []string `json:"characters"`
	URL     string    `json:"url"`
	Created time.Time `json:"created"`
}

func NewRaMRate(log *logrus.Entry, xrHost string, apikey string) *Rate{
	return &Rate{
		log: log.WithField("transport","rickandmorty"),
		xrHost: xrHost,
		apiKey: apikey,
	}
}

func (e *Rate) GetAllEpisodes() ([]Resp,error){
	
	url := "https://rickandmortyapi.com/api/episode"

	client := http.Client{}
	req,err := http.NewRequest(http.MethodGet, url,nil)
    var resp []Resp
	if err != nil {
		return resp, fmt.Errorf("RaM api internal server error: %w", err)
	}
	
	res, err := client.Do(req)
	if err != nil {
		return  resp, fmt.Errorf("morti and ricky api server is error: %w", err) 
	}

	if res.Body != nil{
		defer res.Body.Close()
	}

	
	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return  resp,fmt.Errorf("erisodts not found")
		    
	case http.StatusForbidden:
		return resp, fmt.Errorf("invalid amount")
		  
	default:
		body, err := io.ReadAll(res.Body)
		if err != nil {
		 return   resp, fmt.Errorf("unexpected error")   
		}  
 
		return  resp, fmt.Errorf("unexpected status code", res.StatusCode, string(body))
	}


	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return resp, fmt.Errorf("error decoding responce")
	}

	return resp, nil
}

func (e *Rate) GetEpisode(id int) (Resp, error){
	idStr := fmt.Sprintf("%v", id)

	url := "https://rickandmortyapi.com/api/episode/" + idStr

	client := http.Client{}
	req,err := http.NewRequest(http.MethodGet, url,nil)
	var resp Resp
	if err != nil {
		return resp, fmt.Errorf("RaM api internal server error: %w", err)
	}
	
	res, err := client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("morti and ricky api server is error: %w", err) 
	}

	if res.Body != nil{
		defer res.Body.Close()
	}

	
	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return resp, fmt.Errorf("erisodts not found")
		    
	case http.StatusForbidden:
		return resp, fmt.Errorf("invalid amount")
		  
	default:
		body, err := io.ReadAll(res.Body)
		if err != nil {
		 return resp, fmt.Errorf("unexpected error")   
		}  
 
		return resp, fmt.Errorf("unexpected status code", res.StatusCode, string(body))
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return resp, fmt.Errorf("error decoding responce")
	}
	
	return resp, nil
}