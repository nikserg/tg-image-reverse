package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func CheckTineye(picUrl string) (bool, string) {

	response, err := http.Get(picUrl)
	if err != nil {
		return true, err.Error()
	}
	defer response.Body.Close()
	imageContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return true, err.Error()
	}
	imageContentString := string(imageContent)
	fmt.Println(imageContentString)

	url := "https://tineye.com/result_json/"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("/C:/Users/n.zarubin/Downloads/photo_2022-03-31_15-17-46.jpg")
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("image", filepath.Base("/C:/Users/n.zarubin/Downloads/photo_2022-03-31_15-17-46.jpg"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return true, errFile1.Error()
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return true, err.Error()
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return true, err.Error()
	}
	req.Header.Add("Cookie", "order=desc; sort=score; tineye=VPYcq9fTXgtT6BYyjHKN8HN8Xt0_Xr3Kov0zabivVJjtkX9JmKb6UnmouC4mcjcfeywQbTNIgHmStqEM3QXAav_VkmiDu793uAt33Bn3uU4Cdh4k9vl2Rtw856PoBZMXWsLBO7jANUISYufbSzC47Eanl6k0Cs9QAbUIJLsm8U63QUCrLGdWLhrQSztDPDXB2HADkd6G5VAya6OUPSmEdwtR4duBsbQM4HwPshSUNtTf3qP0")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return true, err.Error()
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return true, err.Error()
	}
	fmt.Println(string(body))

	return false, string(body)

}
