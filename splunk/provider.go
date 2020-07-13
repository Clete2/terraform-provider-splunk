package splunk

import (
	"crypto/tls"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"net/http"
	"time"
)

type SplunkProvider struct {
	Client *SplunkClient
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:        providerSchema(),
		DataSourcesMap:providerDataSources(),
		ResourcesMap:  providerResources(),
		ConfigureFunc: providerConfigure,
	}
}

func providerDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_URL", "https://localhost:8089"),
			Description: "Splunk instance URL",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_USERNAME", "admin"),
			Description: "Splunk instance admin username",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_PASSWORD", "changeme"),
			Description: "Splunk instance password",
		},
		"insecure_skip_verify": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("SPLUNK_INSECURE_SKIP_VERIFY", true),
			Description: "insecure skip verification flag",
		},
	}
}

// Returns a map of splunk resources for configuration
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"splunk_global_http_event_collector": globalHttpEventCollector(),
		"splunk_input_http_event_collector": inputHttpEventCollector(),
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise splunk client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client, err := NewClient(
		d.Get("url").(string),
		d.Get("username").(string),
		d.Get("password").(string),
		httpClient(&http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: d.Get("insecure_skip_verify").(bool)},
		},
	}))

	if err != nil {
		return client, err
	}

	provider := &SplunkProvider{
		Client: client,
	}

	return provider, nil
}
