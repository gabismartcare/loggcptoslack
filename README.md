# loggcptoslack

Simple and small implementation to send log entry from GCP logger to Slack. Log should be sent to GCP pub/sub

# How to use?

create a main file

``` go
func main() {

	log.SetFlags(log.Lshortfile)
	port := getOrDefault("PORT", "8080")

	slackListener := http2.NewSlackListener(usecase.NewSlackUsecase(service.NewSlackClient(getOrDefault("SLACK_USERNAME", "gcp"), getOrDefault("SLACK_CHANNEL", "#log"),  getOrDefault("SLACK_WEBHOOK", "https://webhook.slack.com"))))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), slackListener); err != nil {
		log.Fatal(err)
	}
}

func getOrDefault(env string, defaultvalue string) string {
    if p := os.Getenv(env); p != "" {
        return p
    }
    return defaultvalue
}
```