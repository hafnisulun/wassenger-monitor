package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Wassenger struct{}

type session struct {
	Status     string `json:"status"`
	Operative  string `json:"operative"`
	Uptime     string `json:"uptime"`
	LastSyncAt string `json:"lastSyncAt"`
	AppVersion string `json:"appVersion"`
	Error      string `json:"error"`
	Phone      string `json:"phone"`
}

type device struct {
	Id          string  `json:"id"`
	Phone       string  `json:"phone"`
	Alias       string  `json:"alias"`
	Description string  `json:"description"`
	Wid         string  `json:"wid"`
	Status      string  `json:"status"`
	Session     session `json:"session"`
	Info        string  `json:"info"`
	CreatedAt   string  `json:"createdAt"`
	Webhooks    string  `json:"webhooks"`
	Profile     string  `json:"profile"`
}

type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type SlackAttachment struct {
	Fallback string       `json:"fallback"`
	Pretext  string       `json:"pretext"`
	Color    string       `json:"color"`
	Fields   []slackField `json:"fields"`
}

type slackRequestBody struct {
	Attachments []SlackAttachment `json:"attachments"`
}

func (r Wassenger) Monitor(status *string) {
	log.Printf("Monitor start, status: %s\n", *status)

	client := &http.Client{}
	url := os.Getenv("WASSENGER_BASE_URL") + "/devices"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println("[Error] Create request to Wassenger failed, err:", err)
		return
	}

	req.Header.Add("Token", os.Getenv("WASSENGER_TOKEN"))
	res, err := client.Do(req)

	if err != nil {
		log.Println("[Error] Request to Wassenger failed, err:", err)
		return
	}

	defer res.Body.Close()

	var devices []device
	json.NewDecoder(res.Body).Decode(&devices)

	log.Println("Wassenger response status code:", res.StatusCode)
	log.Println("Wassenger response body:", devices)

	slackAttachment := new(SlackAttachment)
	slackAttachment.Color = "good"

	if len(devices) == 0 {
		slackAttachment.Color = "danger"
		slackAttachment.Fallback = "No devices"
		slackAttachment.Pretext = "No devices"
		*status = "no devices"
	} else {
		device := devices[0]

		if device.Session.Status != "online" {
			slackAttachment.Color = "danger"
		} else if *status == "online" {
			log.Println("All is good, monitor finish")
			return
		}

		*status = device.Session.Status

		slackAttachment.Fallback = "Device is " + device.Session.Status
		slackAttachment.Pretext = "Device is " + device.Session.Status

		slackAttachment.Fields = []slackField{
			{
				Title: "Account",
				Value: os.Getenv("WASSENGER_ACCOUNT"),
				Short: true,
			},
			{
				Title: "Phone",
				Value: device.Phone,
				Short: true,
			},
			{
				Title: "Alias",
				Value: device.Alias,
				Short: true,
			},
			{
				Title: "Status",
				Value: device.Status,
				Short: true,
			},
			{
				Title: "Session status",
				Value: device.Session.Status,
				Short: true,
			},
			{
				Title: "Session last sync at",
				Value: device.Session.LastSyncAt,
				Short: true,
			},
		}
	}

	slackRequestBody := &slackRequestBody{
		Attachments: []SlackAttachment{
			*slackAttachment,
		},
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(slackRequestBody)

	log.Println("Slack request body:", payloadBuf.String())

	res, err = http.Post(os.Getenv("SLACK_WEBHOOK_URL"), "application/json", payloadBuf)

	if err != nil {
		log.Println("[Error] Request to Slack failed, err:", err)
		return
	}

	defer res.Body.Close()

	resBodyBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println("[Error] Read response from Slack failed, err:", err)
		return
	}

	resBodyString := string(resBodyBytes)

	log.Printf("Slack response code: %d\n", res.StatusCode)
	log.Println("Slack response body:", resBodyString)

	log.Println("Monitor finish")
}
