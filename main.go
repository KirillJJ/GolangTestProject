/*query example: localhost:8080/api/v1/distance/?from=samara&to=perm */
package main

import (
	"net/http"
	"encoding/json"
	"log"
	"io/ioutil"

	"github.com/julienschmidt/httprouter"
)

const ( 
	apiKey = "AIzaSyAJUvHqy4Yn3QN278toq0Wfg0GXfiLmUBo"
 	keyFrom = "from"
 	keyTo = "to"
 	okStatus = "OK"
 	googleApiLink = "https://maps.googleapis.com/maps/api/distancematrix/json?"
)

type Result struct {
	Distance string `json:"distance"`
	Result string `json:"result"`
}

type GoogleMapsResponse struct {
	_ []string
	_ []string
	Rows []struct {
		Elements []struct {
			Distance struct {
				Text string `json:"text"`
				_ int
			} `json:"distance"`
			_ struct {
				_ string
				_ int
			}
			Status string `json:"status"`
		}`json:"elements"`
	}`json:"rows"`
	Status string `json:"status"`
}

func getDistance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	from := r.URL.Query().Get(keyFrom)
	to := r.URL.Query().Get(keyTo)
	log.Println(from)
	log.Println(to)

	if from == "" || to == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := http.Client{}

	resp, err := client.Get( googleApiLink +
		"origins=" + from +
		"&destinations=" + to +
		"&key" + apiKey)
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(resp)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	google := GoogleMapsResponse{}
	err = json.Unmarshal(body, &google)
	log.Println(google)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if status := google.Status; status != okStatus {
		log.Println(status)
		writeResponse(w, Result{"", status})
		return
	}

	if status := google.Rows[0].Elements[0].Status; status != okStatus {
		log.Println(status)
		writeResponse(w, Result{"", status})
		return
	}

	result := Result{google.Rows[0].Elements[0].Distance.Text, okStatus}
	writeResponse(w, result)
}

func writeResponse(w http.ResponseWriter, body Result) {
	res, err := json.Marshal(body)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func main() {

	router := httprouter.New()
	router.GET("/api/v1/distance/",getDistance)

	http.ListenAndServe(":8080", router)
}
