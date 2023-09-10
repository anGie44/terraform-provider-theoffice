package theoffice

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientDefault(t *testing.T) {
	_, err := NewClient(&Config{})
	assert.NoError(t, err)
}

func TestClientConnections(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/season/1/format/connections", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		_, err := w.Write([]byte(`[{"episode":1,"episode_name":"Pilot","links":[{"source":"Michael","target":"Pam","value":3},{"source":"Phyllis","target":"Stanley","value":2}],"nodes":[{"id":"Pam"},{"id":"Michael"},{"id":"Phyllis"},{"id":"Stanley"}]}]`))
		assert.NoError(t, err)
	}))
	defer srv.Close()

	c, err := NewClient(&Config{
		Address: srv.URL,
	})
	assert.NoError(t, err)

	resp, err := c.GetConnections(context.Background(), 1)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resp.Connections))
	assert.Equal(t, 1, resp.Connections[0].Episode)
	assert.Equal(t, "Pilot", resp.Connections[0].EpisodeName)
	assert.Equal(t, 2, len(resp.Connections[0].Links))
	assert.Equal(t, "Michael", resp.Connections[0].Links[0].Source)
	assert.Equal(t, "Pam", resp.Connections[0].Links[0].Target)
	assert.Equal(t, 3, resp.Connections[0].Links[0].Value)
	assert.Equal(t, "Phyllis", resp.Connections[0].Links[1].Source)
	assert.Equal(t, "Stanley", resp.Connections[0].Links[1].Target)
	assert.Equal(t, 2, resp.Connections[0].Links[1].Value)
	assert.Equal(t, 4, len(resp.Connections[0].Nodes))
	assert.Equal(t, "Pam", resp.Connections[0].Nodes[0].ID)
	assert.Equal(t, "Michael", resp.Connections[0].Nodes[1].ID)
	assert.Equal(t, "Phyllis", resp.Connections[0].Nodes[2].ID)
	assert.Equal(t, "Stanley", resp.Connections[0].Nodes[3].ID)
}

func TestClientQuotes_season(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/season/1/format/quotes", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		_, err := w.Write([]byte(`[{"season": 1,"episode": 1,"scene": 1,"episode_name": "Diversity Day","character": "Jim","quote": "Really?"}]`))
		assert.NoError(t, err)
	}))
	defer srv.Close()

	c, err := NewClient(&Config{
		Address: srv.URL,
	})
	assert.NoError(t, err)

	resp, err := c.GetQuotes(context.Background(), 1, 0)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resp.Quotes))
	assert.Equal(t, 1, resp.Quotes[0].Season)
	assert.Equal(t, 1, resp.Quotes[0].Episode)
	assert.Equal(t, 1, resp.Quotes[0].Scene)
	assert.Equal(t, "Diversity Day", resp.Quotes[0].EpisodeName)
	assert.Equal(t, "Jim", resp.Quotes[0].Character)
	assert.Equal(t, "Really?", resp.Quotes[0].Quote)
}

func TestClientQuotes_episode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/season/1/episode/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		_, err := w.Write([]byte(`[{"season": 1,"episode": 1,"scene": 1,"episode_name": "Diversity Day","character": "Jim","quote": "Really?"}]`))
		assert.NoError(t, err)
	}))
	defer srv.Close()

	c, err := NewClient(&Config{
		Address: srv.URL,
	})
	assert.NoError(t, err)

	resp, err := c.GetQuotes(context.Background(), 1, 1)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resp.Quotes))
	assert.Equal(t, 1, len(resp.Quotes))
	assert.Equal(t, 1, resp.Quotes[0].Season)
	assert.Equal(t, 1, resp.Quotes[0].Episode)
	assert.Equal(t, 1, resp.Quotes[0].Scene)
	assert.Equal(t, "Diversity Day", resp.Quotes[0].EpisodeName)
	assert.Equal(t, "Jim", resp.Quotes[0].Character)
	assert.Equal(t, "Really?", resp.Quotes[0].Quote)
}

func TestClientQuotes_none(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/season/1/format/quotes", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		_, err := w.Write([]byte(`[]`))
		assert.NoError(t, err)
	}))
	defer srv.Close()

	c, err := NewClient(&Config{
		Address: srv.URL,
	})
	assert.NoError(t, err)

	resp, err := c.GetQuotes(context.Background(), 1, 0)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(resp.Quotes))
}
