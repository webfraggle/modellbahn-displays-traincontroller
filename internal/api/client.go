package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	endpoint string
	http     *http.Client
}

type TrainSlot struct {
	Nr    int
	Train TrainPayload
}

type TrainPayload struct {
	VonNach    string `json:"vonnach"`
	Nr         string `json:"nr"`
	Zeit       string `json:"zeit"`
	Via        string `json:"via"`
	Abw        string `json:"abw"`
	Hinweis    string `json:"hinweis"`
	Fusszeile  string `json:"fusszeile"`
	Abschnitte string `json:"abschnitte"`
	Reihung    string `json:"reihung"`
	Path       string `json:"path"`
}

func NewClient(endpoint string, timeoutMs int) *Client {
	return &Client{
		endpoint: endpoint,
		http:     &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond},
	}
}

func (c *Client) SkipNext(path string) error {
	return c.get(fmt.Sprintf("%s/skipNext?path=%s", c.endpoint, path))
}

func (c *Client) SkipPrev(path string) error {
	return c.get(fmt.Sprintf("%s/skipPrev?path=%s", c.endpoint, path))
}

func (c *Client) SetTime(path, timeStr string) error {
	return c.postForm(c.endpoint+"/setTime", url.Values{
		"path": {path},
		"time": {timeStr},
	})
}

func (c *Client) ShowImage(path, filename string) error {
	return c.get(fmt.Sprintf("%s/showImage?path=%s&filename=%s", c.endpoint, path, filename))
}

func (c *Client) SetTrains(trains []TrainSlot) error {
	for _, t := range trains {
		u := fmt.Sprintf("%s/zug%d", c.endpoint, t.Nr)
		if err := c.postJSON(u, t.Train); err != nil {
			return fmt.Errorf("zug%d: %w", t.Nr, err)
		}
	}
	return nil
}

// Ping checks whether the endpoint is reachable.
func (c *Client) Ping() error {
	return c.get(c.endpoint)
}

// ParseTrain parses a pipe-separated train string into a TrainSlot.
// Format: "TrainID|Time|Destination|Via|Delay|SpecialInfo"
func ParseTrain(nr int, s, path string) TrainSlot {
	parts := splitParts(s, 6)
	return TrainSlot{
		Nr: nr,
		Train: TrainPayload{
			Nr:      parts[0],
			Zeit:    parts[1],
			VonNach: repairUTF8(parts[2]),
			Via:     repairUTF8(parts[3]),
			Abw:     parts[4],
			Hinweis: repairUTF8(parts[5]),
			Path:    path,
		},
	}
}

func splitParts(s string, n int) []string {
	parts := make([]string, n)
	raw := []byte(s)
	i := 0
	for field := 0; field < n; field++ {
		start := i
		for i < len(raw) && raw[i] != '|' {
			i++
		}
		parts[field] = string(raw[start:i])
		if i < len(raw) {
			i++ // skip '|'
		}
	}
	return parts
}

func (c *Client) get(url string) error {
	resp, err := c.http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
}

func (c *Client) postForm(rawURL string, values url.Values) error {
	resp, err := c.http.PostForm(rawURL, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
}

func (c *Client) postJSON(rawURL string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := c.http.Post(rawURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
}
