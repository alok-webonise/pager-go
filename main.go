package pager-go

import (
	"bytes"
	"encoding/json"
	"net/http"
	"log"

)




type AlertNotificationRequest struct {
	Payload     AlertPayload `json:"payload"`
	RoutingKey  string       `json:"routing_key"`
	EventAction string       `json:"event_action"`
	Client      string
	DedupKey    string `json:"dedup_key"`
}

type AlertPayload struct {
	Summary  string `json:"summary"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
	ErrorSource string `json:"error_source"`
	RoutingKey string `json:"routing_key"`
}

func (a *AlertPayload) Error() string {
	a.Source = "api-server :  ENV : "+ a.Source
	a.Severity = "error"
	createPagerDutyAlert(a)
	return a.Summary
}

func createPagerDutyAlert(payload *AlertPayload) {
	alertRequestBody := &AlertNotificationRequest{
		Payload:     *payload,
		RoutingKey:  payload.RoutingKey,
		EventAction: "trigger",
		DedupKey:    "DedupKeyToSingulariseTheIncident",
	}

	reqByte, err := json.Marshal(alertRequestBody)
  if err != nil {
    	log.Println(err)
      return 
  }
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://events.pagerduty.com/v2/enqueue", bytes.NewBuffer(reqByte))
	if err != nil{
		log.Println(err)
    return 
	}
	resp, err := client.Do(req)
	if err != nil{
		log.Println(err)
	}
	resp.Body.Close()
}

