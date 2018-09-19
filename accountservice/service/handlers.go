package service

import(
	"github.com/dasetova/goblog/accountservice/dbclient"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
	"fmt"
	"net"
	"github.com/dasetova/goblog/accountservice/model"
	"io/ioutil"
)

var DBClient dbclient.IBoltClient
var isHealthy = true //To make unhealthy service an probe the docker restart
var client = &http.Client{}

func init() {
	var transport http.RoundTripper = &http.Transport{
			DisableKeepAlives: true,
	}
	client.Transport = transport
}

func GetAccount(w http.ResponseWriter, r *http.Request){
	var accountId = mux.Vars(r)["accountId"]
	account, err := DBClient.QueryAccount(accountId)
	account.ServedBy = getIP()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	quote, err := getQuote()
	if err == nil {
			account.Quote = quote
	}

	data, _ := json.Marshal(account)
	writeJsonResponse(w, http.StatusOK, data)
}

func getQuote() (model.Quote, error) {
	req, _ := http.NewRequest("GET", "http://quotes-service:8080/api/quote?strength=4", nil)
	resp, err := client.Do(req)

	if err == nil && resp.StatusCode == 200 {
			quote := model.Quote{}
			bytes, _ := ioutil.ReadAll(resp.Body)
			json.Unmarshal(bytes, &quote)
			return quote, nil
	} else {
			return model.Quote{}, fmt.Errorf("Some error")
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request){
	dbUp := DBClient.Check()
	if dbUp && isHealthy {
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJsonResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Database unaccessible"})
		writeJsonResponse(w, http.StatusServiceUnavailable, data)
	}
}

func writeJsonResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

func SetHealthyState(w http.ResponseWriter, r *http.Request){
	var state, err = strconv.ParseBool(mux.Vars(r)["state"])

	if err != nil {
		fmt.Println("Invalid request to SetHealthyState, allowed values are true or false")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isHealthy = state
	w.WriteHeader(http.StatusOK)
}

func getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("Unable to determine local IP address (non loopback). Exiting.")
}

