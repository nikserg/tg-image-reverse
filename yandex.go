package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func CheckYandex(picUrl string) (bool, string) {
	searchUrl := "https://yandex.ru/images/search?rpt=imageview&url=" + url.QueryEscape(picUrl)


	response, err := http.Get(searchUrl)
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
	if !strings.Contains(bodyString, "Таких же изображений не найдено") {

		resultMessage = ""
		r := regexp.MustCompile(`<div class="CbirSites-ItemTitle"><a href="(.+?)" target="_blank" class="Link Link_theme_normal">(.+?)</a>`)
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
