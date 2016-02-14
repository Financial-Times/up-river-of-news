package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/Financial-Times/go-fthealth"
)

// Healthcheck offers methods to measure application health.
type Healthcheck struct {
	client http.Client
	config AppConfig
}

func (h *Healthcheck) checkHealth() func(w http.ResponseWriter, r *http.Request) {
	return fthealth.HandlerParallel("Dependent services healthcheck", "Checks if all the dependent services are reachable and healthy.", h.messageQueueProxyReachable())
}

func (h *Healthcheck) gtg(writer http.ResponseWriter, req *http.Request) {
	healthChecks := []func() error{h.checkAggregateMessageQueueProxiesReachable}

	for _, hCheck := range healthChecks {
		if err := hCheck(); err != nil {
			writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}
}

func (h *Healthcheck) messageQueueProxyReachable() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "The UPP River of News channel will not be updated with news of new published content",
		Name:             "MessageQueueProxyReachable",
		PanicGuide:       "TODO add link",
		Severity:         1,
		TechnicalSummary: "Message queue proxy is not reachable/healthy",
		Checker:          h.checkAggregateMessageQueueProxiesReachable,
	}

}

func (h *Healthcheck) checkAggregateMessageQueueProxiesReachable() error {

	addresses := h.config.QueueConf.Addrs
	errMsg := ""
	for i := 0; i < len(addresses); i++ {
		error := h.checkMessageQueueProxyReachable(addresses[i])
		if error == nil {
			return nil
		}
		errMsg = errMsg + fmt.Sprintf("For %s there is an error %v \n", addresses[i], error.Error())
	}

	return errors.New(errMsg)

}

func (h *Healthcheck) checkMessageQueueProxyReachable(address string) error {
	req, err := http.NewRequest("GET", address+"/topics", nil)
	if err != nil {
		log.Warnf("Could not connect to proxy: %v", err.Error())
		return err
	}

	if len(h.config.QueueConf.AuthorizationKey) > 0 {
		req.Header.Add("Authorization", h.config.QueueConf.AuthorizationKey)
	}

	if len(h.config.QueueConf.Queue) > 0 {
		req.Host = h.config.QueueConf.Queue
	}

	resp, err := h.client.Do(req)
	if err != nil {
		log.Warnf("Could not connect to proxy: %v", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Proxy returned status: %d", resp.StatusCode)
		return errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	return checkIfTopicIsPresent(body, h.config.QueueConf.Topic)

}

func checkIfTopicIsPresent(body []byte, searchedTopic string) error {
	var topics []string

	err := json.Unmarshal(body, &topics)
	if err != nil {
		return fmt.Errorf("Error occured and topic could not be found. %v", err.Error())
	}

	for _, topic := range topics {
		if topic == searchedTopic {
			return nil
		}
	}

	return errors.New("Topic was not found")
}
