package main

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func CheckYandex(picUrl string) (bool, string) {
	searchUrl := "https://yandex.ru/images/search?rpt=imageview&url=" + url.QueryEscape(picUrl)

	tlsConfig := &tls.Config{
		MaxVersion: tls.VersionTLS12,
	}
	client := http.Client{
		Transport: &http.Transport{
			DialTLS: func(network, addr string) (net.Conn, error) {
				return tls.Dial(network, addr, tlsConfig)
			},
		},
	}
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		panic(err)
	}
	req.Header = http.Header{
		"Host":                      []string{"yandex.ru"},
		"User-Agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36"},
		"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Accept-Language":           []string{"ru-RU,ru;q=0.9,en-GB;q=0.8,en;q=0.7,en-US;q=0.6"},
		"Cache-Control":             []string{"no-cache"},
		"Connection":                []string{"keep-alive"},
		"device-memory":             []string{"8"},
		"downlink":                  []string{"4.25"},
		"dpr":                       []string{"1"},
		"ect":                       []string{"4.25"},
		"Pragma":                    []string{"no-cache"},
		"rtt":                       []string{"50"},
		"Upgrade-Insecure-Requests": []string{"1"},
		"viewport-width":            []string{"1920"},
	}

	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	var resultMessage string
	if strings.Contains(bodyString, "нашему сервису временно запрещён") {
		return false, "Яндекс заблокировал бота, адрес для проверки: " + searchUrl
	} else if !strings.Contains(bodyString, "Таких же изображений не найдено") {

		resultMessage = ""
		r := regexp.MustCompile(`<div class="CbirSites-ItemTitle"><a href="(.+?)" target="_blank" class="Link [^"]+">(.+?)</a>`)
		for index, match := range r.FindAllStringSubmatch(bodyString, -1) {
			if index > 5 {
				break
			}
			resultMessage += match[1] + " " + match[2] + "\n"
		}
		return false, resultMessage

	} else {

		return true, ""
	}

}
