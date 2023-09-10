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
