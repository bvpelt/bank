/*
 * Simple API overview
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 2.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}

func HelloUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	type User struct {
		Name string `json:"name"`
	}

	var user string

	keys := strings.Split(r.URL.Path, "/")
	if len(keys) < 1 {
		log.Println("Url parameter missing")
		return
	} else {
		user = keys[len(keys)-1]
		log.Printf("Length: %d Keys: %v ", len(keys), keys)
		json.NewEncoder(w).Encode("Hello: " + user)
	}

}
