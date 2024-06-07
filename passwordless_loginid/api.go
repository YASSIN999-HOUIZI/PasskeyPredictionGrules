package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var df DeviceFact
		err := json.NewDecoder(r.Body).Decode(&df)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var devices []*DeviceFact
		devices = append(devices, &df)

		err = ProcessDevices(devices, rules)
		if err != nil {
			panic(err)
		}

		// Return the result as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(df.Output)
	})

	fmt.Println("Server is running on port 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil{
		panic(err.Error())
	}

}