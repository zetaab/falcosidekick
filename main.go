package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Issif/falcosidekick/outputs"
	"github.com/Issif/falcosidekick/types"
)

// Globale variables
var slackClient, datadogClient, alertmanagerClient, elasticsearchClient *outputs.Client
var config *types.Configuration
var stats *types.Statistics

func init() {
	config = getConfig()
	stats = getInitStats()

	enabledOutputsText := "[INFO]  : Enabled Outputs : "
	if config.Slack.WebhookURL != "" {
		var err error
		slackClient, err = outputs.NewClient("Slack", config.Slack.WebhookURL, config, stats)
		if err != nil {
			config.Slack.WebhookURL = ""
		} else {
			enabledOutputsText += "Slack "
		}
	}
	if config.Datadog.APIKey != "" {
		var err error
		datadogClient, err = outputs.NewClient("Datadog", outputs.DatadogURL+"?apikey="+config.Datadog.APIKey, config, stats)
		if err != nil {
			config.Datadog.APIKey = ""
		} else {
			enabledOutputsText += "Datadog "
		}
	}
	if config.Alertmanager.HostPort != "" {
		var err error
		alertmanagerClient, err = outputs.NewClient("AlertManager", config.Alertmanager.HostPort+outputs.AlertmanagerURI, config, stats)
		if err != nil {
			config.Alertmanager.HostPort = ""
		} else {
			enabledOutputsText += "AlertManager "
		}
	}
	if config.Elasticsearch.HostPort != "" {
		var err error
		elasticsearchClient, err = outputs.NewClient("Elasticsearch", config.Elasticsearch.HostPort+"/"+config.Elasticsearch.Index+"/"+config.Elasticsearch.Type, config, stats)
		if err != nil {
			config.Elasticsearch.HostPort = ""
		} else {
			enabledOutputsText += "Elasticsearch "
		}
	}

	log.Printf("%v\n", enabledOutputsText)
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/test", testHandler)

	log.Printf("[INFO]  : Falco Sidekick is up and listening on port %v\n", config.ListenPort)
	if config.Debug {
		log.Printf("[INFO]  : Debug mode : %v\n", config.Debug)
	}
	if err := http.ListenAndServe(":"+strconv.Itoa(config.ListenPort), nil); err != nil {
		log.Fatalf("[ERROR] : %v\n", err.Error())
	}
}
