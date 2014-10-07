package extensions

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

const TokenID = "123"

func ServiceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{
			TokenID: TokenID,
		},
		Endpoint: th.Endpoint(),
	}
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/extensions", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `
{
	"extensions": [
		{
			"updated": "2013-01-20T00:00:00-00:00",
			"name": "Neutron Service Type Management",
			"links": [],
			"namespace": "http://docs.openstack.org/ext/neutron/service-type/api/v1.0",
			"alias": "service-type",
			"description": "API for retrieving service providers for Neutron advanced services"
		}
	]
}
			`)
	})

	count := 0

	List(ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractExtensions(page)
		if err != nil {
			t.Errorf("Failed to extract extensions: %v", err)
		}

		expected := []Extension{
			Extension{
				Updated:     "2013-01-20T00:00:00-00:00",
				Name:        "Neutron Service Type Management",
				Links:       []interface{}{},
				Namespace:   "http://docs.openstack.org/ext/neutron/service-type/api/v1.0",
				Alias:       "service-type",
				Description: "API for retrieving service providers for Neutron advanced services",
			},
		}

		th.AssertDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/extensions/agent", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
	"extension": {
		"updated": "2013-02-03T10:00:00-00:00",
		"name": "agent",
		"links": [],
		"namespace": "http://docs.openstack.org/ext/agent/api/v2.0",
		"alias": "agent",
		"description": "The agent management extension."
	}
}
		`)

		ext, err := Get(ServiceClient(), "agent").Extract()
		th.AssertNoErr(t, err)

		th.AssertEquals(t, ext.Updated, "2013-02-03T10:00:00-00:00")
		th.AssertEquals(t, ext.Name, "agent")
		th.AssertEquals(t, ext.Namespace, "http://docs.openstack.org/ext/agent/api/v2.0")
		th.AssertEquals(t, ext.Alias, "agent")
		th.AssertEquals(t, ext.Description, "The agent management extension.")
	})
}
