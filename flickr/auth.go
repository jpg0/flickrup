package flickr

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	"bufio"
	"os"
	"strings"
	"gopkg.in/masci/flickr.v2"
)

func getToken(client *flickr.FlickrClient) (*flickr.OAuthToken, error) {
	filepath := os.ExpandEnv("${HOME}/.flickrup")

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		token, err := requestToken(client)

		if err != nil {
			return nil, err
		}

		bytes, err := yaml.Marshal(token)

		err = ioutil.WriteFile(filepath, bytes, 0644)

		if err != nil {
			return nil, err
		}

		return token, nil

	} else {
		//load from file

		file, err := os.Open(filepath)
		defer file.Close()

		if err != nil {
			return nil, err
		}

		tokens := new(flickr.OAuthToken)

		bytes, err := ioutil.ReadAll(file)

		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(bytes, tokens)

		if err != nil {
			return nil, err
		}

		return tokens, nil
	}
}

func requestToken(client *flickr.FlickrClient) (*flickr.OAuthToken, error){
	requestToken, err := flickr.GetRequestToken(client)

	if err != nil {
		return nil, err
	}

	url, err := flickr.GetAuthorizeUrl(client, requestToken)

	if err != nil {
		return nil, err
	}

	fmt.Println("Open this url in your process to complete the authentication process: ", url)
	fmt.Println("Copy here the number given when you complete the process.")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	return flickr.GetAccessToken(client, requestToken, strings.TrimSpace(text))

}
