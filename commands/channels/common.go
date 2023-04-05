package channels

import (
	"fmt"

	"github.com/theherk/viper"
)

func wsHost(cluster string) string {
	host := viper.GetString("wshost")
	if host == "" {
		host = fmt.Sprintf("ws-%s.pusher.com", cluster)
	}
	return host
}

func httpHost(cluster string) string {
	host := viper.GetString("httphost")
	if host == "" {
		host = fmt.Sprintf("sockjs-%s.pusher.com", cluster)
	}
	return host
}

func apiHost(cluster string) string {
	host := viper.GetString("apihost")
	if host == "" {
		host = fmt.Sprintf("api-%s.pusher.com", cluster)
	}
	return host
}
