package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client for interacting with AnvilService
type AnvilClient struct {
	Address string
}

// Publish a module artifact to the anvil repo
func (c *AnvilClient) PublishModule(b io.Reader, fileName string) (string, error) {
	url := fmt.Sprintf("http://%s/admin/module/%s", c.Address, fileName)

	req, err := http.NewRequest("PUT", url, b)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", processErrors(body, resp.StatusCode)
	}
	var adminMod *AdminModule
	err = json.Unmarshal(body, &adminMod)
	return adminMod.FileUri, err
}

// query releases  filtered by specified module
func (c *AnvilClient) GetReleaseByModule(user string, module string) (*Response, error) {
	var entity *Response

	url := fmt.Sprintf("http://%s/v3/releases?module=%s-%s", c.Address, user, module)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return entity, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return entity, nil
	}
	if resp.StatusCode != http.StatusOK {
		return entity, processErrors(body, resp.StatusCode)
	}
	err = json.Unmarshal(body, &entity)
	return entity, err
}

func processErrors(body []byte, code int) error {
	var e *ErrorResponse

	if err := json.Unmarshal(body, &e); err != nil {
		return err
	}
	return fmt.Errorf("%d: %s", code, strings.Join(e.Errors, "; "))

}
