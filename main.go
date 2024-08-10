package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	urlBigIP := os.Getenv("URL_BIG_IP")

	http.HandleFunc("/del", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.RawQuery

		// ถ้า URL_BIG_IP ไม่มีการตั้งค่า ให้ใช้ IP:PORT ของต้นทางที่เรียกมา
		if urlBigIP == "" {
			urlBigIP = r.RemoteAddr
		}

		targetUrl := fmt.Sprintf("http://%s/del?%s", urlBigIP, queryParams)

		// เพิ่มการบันทึก targetUrl
		log.Printf("Calling target URL: %s", targetUrl)

		resp, err := http.Get(targetUrl)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching from target URL: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})

	log.Printf("Server running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}