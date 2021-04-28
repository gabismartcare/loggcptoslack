package pub_sub_http

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/logging/v2"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gabismartcare/loggcptoslack/domain/model"
	"github.com/gabismartcare/loggcptoslack/usecase"
)

type SlackListener struct {
	slackUseCase usecase.SlackUsecase
}

func NewSlackListener(slackUsecase usecase.SlackUsecase) SlackListener {
	return SlackListener{slackUseCase: slackUsecase}
}

type PubSubMessage struct {
	Message Base64Message `json:"message"`
}
type Base64Message struct {
	Data []byte `json:"data"`
}

func (s SlackListener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	logentry, err := asLogEntry(body)
	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte("Cannot get message " + string(body) + " " + err.Error()))
		return
	}

	payload, err := extractPayload(logentry)
	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte("Cannot get payload " + err.Error()))
		return
	}

	if shouldIgnore(logentry) {
		return
	}

	service := ""
	if logentry.Resource != nil && len(logentry.Resource.Labels) > 0 {
		if job, ok := logentry.Resource.Labels["job_id"]; ok {
			service = job
		} else if method, ok := logentry.Resource.Labels["method"]; ok {
			service = method
		} else if conf, ok := logentry.Resource.Labels["configuration_name"]; ok {
			service = conf
		} else {
			service = logentry.LogName
		}
	}
	icon := ""

	switch logentry.Severity {
	case "FATAL":
		icon = ":skull:"
	case "ERROR":
		icon = ":red_circle:"
	case "WARN":
		icon = ":large_yellow_circle:"
	}
	if err := s.slackUseCase.SendMessageToSlack(model.SimpleSlackRequest{
		Text:      fmt.Sprintf("%s : GCP service %s %s, ", logentry.Severity, service, payload),
		IconEmoji: icon,
	}); err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
}

func shouldIgnore(logentry logging.LogEntry) bool {
	return logentry.HttpRequest != nil && strings.HasSuffix(logentry.LogName, "endpoints_log") && logentry.HttpRequest.Status < 404
}

func extractPayload(logentry logging.LogEntry) (payload string, err error) {
	if logentry.HttpRequest != nil && logentry.HttpRequest.Status != 0 {
		request := logentry.HttpRequest
		payload = fmt.Sprintf("Error when calling \n`%s %s`\n (Status %d) (User-Agent : %s) ", request.RequestMethod, request.RequestUrl, request.Status, request.UserAgent)
	} else {
		payload = logentry.TextPayload
		if payload == "" {
			msg, err := logentry.JsonPayload.MarshalJSON()
			if err != nil {
				return payload, err
			}
			if len(msg) > 0 {
				var data map[string]interface{}
				err := json.Unmarshal(msg, &data)
				if err != nil {
					return string(msg), err
				}
				if s, ok := data["jobName"]; ok {
					payload = s.(string)
				}
				if s, ok := data["url"]; ok {
					payload += " " + s.(string)
				}
				if s, ok := data["log_message"]; ok {
					payload = s.(string)
				}
				if s, ok := data["grpc_status_code"]; ok {
					payload += " " + s.(string)
				}
				if payload == "" {
					return " (" + string(msg) + ")", nil
				}
				return payload, nil
			}
		}
		if payload == "" {
			payload, err = extractRawMessage(logentry.ProtoPayload)
			if err != nil {
				return payload, err
			}
		}
	}
	return payload, nil
}

func asLogEntry(body []byte) (logging.LogEntry, error) {
	pubSubMessage := PubSubMessage{}
	logentry := logging.LogEntry{}
	if err := json.Unmarshal(body, &pubSubMessage); err != nil {
		return logentry, err
	}
	if err := json.Unmarshal(pubSubMessage.Message.Data, &logentry); err != nil {
		return logentry, err
	}
	return logentry, nil
}

func extractRawMessage(logentry googleapi.RawMessage) (string, error) {
	data, err := logentry.MarshalJSON()
	if err != nil {
		return "", err
	}
	if len(data) > 0 {
		return string(data), nil
	}
	return "", nil
}
