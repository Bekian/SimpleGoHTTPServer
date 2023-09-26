package main

import (
	"encoding/json"
	"fmt"
	"io"

	//"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func countdown(w http.ResponseWriter, targetTime time.Time, userTimezone string) {
	for {
		currentTime := time.Now()
		if currentTime.After(targetTime) {
			fmt.Println("Countdown reached!")
			break
		}

		remainingTime := targetTime.Sub(currentTime)
		days := remainingTime.Hours() / 24
		hours := int(remainingTime.Hours()) % 24
		minutes := int(remainingTime.Minutes()) % 60

		fmt.Fprintf(w, "Time remaining: %02d days %02d hours %02d minutes\n", int(days), hours, minutes)

		time.Sleep(time.Minute) // Wait for 1 minute before updating the countdown
	}
}

type TimeZoneInfo struct {
	Timezone string `json:"timezone"`
}

func getUserTimezone() (string, error) {
	resp, err := http.Get("https://worldtimeapi.org/api/ip")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tzInfo TimeZoneInfo
	if err := json.Unmarshal(body, &tzInfo); err != nil {
		return "", err
	}

	return tzInfo.Timezone, nil
}

func calendar(w http.ResponseWriter, req *http.Request) {
	// Determine the user's timezone automatically
	userTimezone, err := getUserTimezone()
	if err != nil {
		fmt.Println("Error determining user's timezone:", err)
		return
	}

	// Define the target date and time (in UTC)
	targetTime, err := time.Parse("2006-01-02 15:04:05", "2023-12-31 23:59:59")
	if err != nil {
		fmt.Println("Error parsing target time:", err)
		return
	}

	// Convert the target time to the user's timezone
	location, err := time.LoadLocation(userTimezone)
	if err != nil {
		fmt.Println("Error loading timezone:", err)
		return
	}
	targetTime = targetTime.In(location)

	// Start the countdown
	fmt.Fprint(w, "Countdown to ", targetTime.Format("2006-01-02 15:04:05"), " in ", userTimezone)
	countdown(w, targetTime, userTimezone)

	// tmpl, err := template.ParseFiles("templates/base.html")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// err = tmpl.Execute(w, data)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
	log.Printf("Received request for /hello from %s\n", req.RemoteAddr)
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
	log.Printf("Received request for /headers from %s\n", req.RemoteAddr)
}

func base(w http.ResponseWriter, req *http.Request) {
	routes := []string{"/hello", "/headers", "/calendar"}
	fmt.Fprint(w, "<title>Home</title> <h1>Endpoints:</h1> <ul>")
	for _, url := range routes {
		fmt.Fprint(w, "<li>", url, "</li>")
	}
	fmt.Fprint(w, "</ul>")
}

func main() {

	///// NEW STUFF

	///// NEW STUFF

	// TODO: add christmas endpoint using datetime
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090" // Default port if PORT is not set
	}

	http.HandleFunc("/", base)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/calendar", calendar)

	log.Printf("Server is running on port %s\n", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}
