package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Enc  bool   `json:"enc"`
}

var channels []Channel

func main() {
	err := loadCSV("channels.csv")
	if err != nil {
		log.Fatalf("Failed to load CSV: %v", err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/search", searchHandler)

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func loadCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, record := range records {
		// first row if id is not a number
		if i == 0 && isNotNumber(record[1]) {
			continue
		}
		if len(record) < 8 {
			continue
		}
		channels = append(channels, Channel{
			ID:   record[1],
			Name: record[2],
			Enc:  record[7] == "C",
		})
	}
	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	var builder strings.Builder
	for _, channel := range channels {
		badge := ""
		if channel.Enc {
			badge = `<span class="text-xs text-white bg-gray-500 px-2 py-1 rounded ml-2">C</span>`
		}
		builder.WriteString(fmt.Sprintf(`
			<div class="card border p-4 rounded shadow">
				<h3 class="text-lg font-bold">%s%s</h3>
				<p class="text-gray-600">%s</p>
			</div>
		`, channel.Name, badge, channel.ID))
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, builder.String())
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	if query == "" {
		listHandler(w, r)
		return
	}

	var builder strings.Builder
	for _, channel := range channels {
		if strings.Contains(strings.ToLower(channel.Name), query) {
			badge := ""
			if channel.Enc {
				badge = `<span class="text-xs text-white bg-gray-500 px-2 py-1 rounded ml-2">C</span>`
			}
			builder.WriteString(fmt.Sprintf(`
				<div class="card border p-4 rounded shadow">
					<h3 class="text-lg font-bold">%s%s</h3>
					<p class="text-gray-600">%s</p>
				</div>
			`, channel.Name, badge, channel.ID))
		}
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, builder.String())
}

func isNotNumber(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err != nil
}
