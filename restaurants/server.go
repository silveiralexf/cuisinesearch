package restaurants

import (
	"fmt"
	"log"
	"net/http"

	"github.com/silveiralexf/cuisinesearch/logging"
)

// hard-coding values for convenience, in real scenario would read from flags/variables/etc
var (
	defaultPort         = ":8080"
	cuisinesfilePath    = "files/cuisines.csv"
	restaurantsFilePath = "files/restaurants.csv"
)

// ListenAndServe will will start listening on a given port and serving requests to the specified resources
func ListenAndServe() {
	http.HandleFunc("/", notFound)
	http.HandleFunc("/restaurants", handleRestaurants)
	http.HandleFunc("/restaurants/search", handleRestaurantSearch)
	http.HandleFunc("/favicon.ico", faviconHandler)

	logging.Info("Startup", fmt.Sprintf("listening and serving on port %v", defaultPort))
	log.Fatal(http.ListenAndServe(defaultPort, nil))
}

// default response when a invalid resource is provided
func notFound(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("endpoint '%v' could not be found", r.RequestURI)

	fmt.Fprint(w, jsonResponseWrap(r, 404, msg))
	logging.Error(r.Host, msg, logging.CallerInfo())

}

// avoids printing useless favicon missing message at the console
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

// conveniently wrap all common fields together for the JSON response along with information on the available resources
func jsonResponseWrap(r *http.Request, statusCode int, msg interface{}) string {
	availableEndpoints := []string{
		"restaurants",
		"restaurants/search",
	}

	spec := searchSpec{}
	spec.Fields = availableEndpoints

	msgJson := fmt.Sprintf(`{"method": "%v", "status_code": "%v", "host": "%v", "error": "%v", "filterspec": "%v"}`,
		r.Method, statusCode, r.Host, msg, spec)
	return msgJson
}

// searchSpec helps providing list of available resources when users receives an error response
type searchSpec struct {
	Fields []string `json:"searchspec"`
}
