package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func runServer(){
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		
		var output Output
		var df DeviceData
		err := json.NewDecoder(r.Body).Decode(&df)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		output,err = ProcessDevice(df.Auth,df.UserPasskeyHistory,df.DeviceFeatures)
		if err != nil {
			panic(err)
		}

		// Return the result as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	})

	fmt.Println("Server is running on port 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil{
		panic(err.Error())
	}

}