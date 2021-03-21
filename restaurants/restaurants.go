package restaurants

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/silveiralexf/cuisinesearch/logging"
	"github.com/silveiralexf/cuisinesearch/searching"
)

// Restaurant represents the different Restaurants by type of Cuisine extracted from CSV file
type Restaurant struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Rating      int    `json:"rating"`
	Distance    int    `json:"distance"`
	Price       int    `json:"price"`
	CuisineID   int    `json:"cuisine_id"`
	CuisineName string `json:"cuisine"`
	Rank        int    `json:"rank"`
}

// byRank implements a sorting interface based on the values from Rank field
type byRank []Restaurant

func (re byRank) Len() int           { return len(re) }
func (re byRank) Less(i, j int) bool { return re[i].Rank < re[j].Rank }
func (re byRank) Swap(i, j int)      { re[i], re[j] = re[j], re[i] }

// Pull return a list of restaurants provided by its given csv file
func (re *Restaurant) Pull(restaurantsFile, cuisinesFile string) ([]Restaurant, error) {
	cs, err := PullCuisinesList(cuisinesFile)
	if err != nil {
		return nil, err
	}

	rawCSVdata, err := readCSVFile(restaurantsFile, 5)
	if err != nil {
		return nil, err
	}

	rs := []Restaurant{}
	for i := 1; i < len(rawCSVdata); i++ {
		re.ID = i
		re.Name = cleanCommonTypos(rawCSVdata[i][0])

		numRating, err := strToInt(rawCSVdata[i][1])
		if err != nil {
			return nil, err
		}
		re.Rating = numRating

		numDistance, err := strToInt(rawCSVdata[i][2])
		if err != nil {
			return nil, err
		}
		re.Distance = numDistance

		numPrice, err := strToInt(rawCSVdata[i][3])
		if err != nil {
			return nil, err
		}
		re.Price = numPrice

		numCuisineID, err := strToInt(rawCSVdata[i][4])
		if err != nil {
			return nil, err
		}
		re.CuisineID = numCuisineID

		re.CuisineName = cleanCommonTypos(getCuisineName(cs, re.CuisineID))

		rs = append(rs, *re)
	}

	return rs, nil
}

// returns the full list of restaurants
func handleRestaurants(w http.ResponseWriter, r *http.Request) {
	re := Restaurant{}
	rs, err := re.Pull(restaurantsFilePath, cuisinesfilePath)
	if err != nil {
		fmt.Fprint(w, jsonResponseWrap(r, 500, err))
		logging.Error(r.Host, err, logging.CallerInfo())
		return
	}

	data, err := json.Marshal(rs)
	if err != nil {
		logging.Error(r.Host, err, logging.CallerInfo())
		fmt.Fprintln(w, err)
		return
	}

	logging.Info(r.Host, fmt.Sprintf("HTTP status %v on '%v'", http.StatusOK, r.RequestURI))
	fmt.Fprintln(w, string(data))
}

// returns the list of restaurants based on provided search terms
func handleRestaurantSearch(w http.ResponseWriter, r *http.Request) {
	searchTerms := getRestaurantSearchTerms(r)
	if len(searchTerms) < 1 {
		msg := fmt.Sprintf("HTTP status %v on '%v': mandatory parameters missed", http.StatusBadRequest, r.RequestURI)

		fmt.Fprint(w, jsonResponseWrap(r, http.StatusBadRequest, msg))
		logging.Error(r.Host, msg, logging.CallerInfo())
		return
	}

	re := Restaurant{}
	rs, err := re.Pull(restaurantsFilePath, cuisinesfilePath)
	if err != nil {
		fmt.Fprint(w, jsonResponseWrap(r, 500, err))
		logging.Error(r.Host, err, logging.CallerInfo())
		return
	}

	err = rankByTerms(rs, searchTerms)
	if err != nil {
		fmt.Fprint(w, jsonResponseWrap(r, 500, err))
		logging.Error(r.Host, err, logging.CallerInfo())
		return
	}

	topFive := getTopFive(rs)

	data, err := json.Marshal(topFive)
	if err != nil {
		logging.Error(r.Host, err, logging.CallerInfo())
		fmt.Fprintln(w, err)
		return
	}
	logging.Info(r.Host, fmt.Sprintf("HTTP status %v on '%v'", http.StatusOK, r.RequestURI))
	fmt.Fprintln(w, string(data))
}

func getRestaurantSearchTerms(r *http.Request) []byte {
	terms := make(map[string]interface{})
	validKeys := []string{"name", "cuisine", "distance", "price", "rating"}

	for _, term := range validKeys {
		searchValue := r.URL.Query().Get(term)
		if searchValue != "" {
			num, err := strToInt(searchValue)
			if err != nil {
				terms[term] = searchValue
			} else {
				terms[term] = num
			}
		}
	}

	b, err := json.Marshal(terms)
	if err != nil {
		return nil
	}
	return b
}

func (re *Restaurant) buildSearchCriteria(rs []Restaurant, searchTerms []byte) error {
	err := json.Unmarshal(searchTerms, &re)
	if err != nil {
		return err
	}
	return nil
}

func rankByTerms(rs []Restaurant, searchTerms []byte) error {
	re := Restaurant{}
	err := re.buildSearchCriteria(rs, searchTerms)
	if err != nil {
		return err
	}

	for i := range rs {
		rank := 0

		if re.Name != "" {
			rank = rank + searching.GetStringsDistance([]rune(re.Name), []rune(rs[i].Name))
		} else {
			rank = rank + 1
		}

		if re.CuisineName != "" {
			rank = rank + searching.GetStringsDistance([]rune(re.CuisineName), []rune(rs[i].CuisineName))
		} else {
			rank = rank + 1
		}

		if re.Rating != 0 {
			if re.Rating < rs[i].Rating {
				rank = rank + searching.GetIntegersDistance(re.Rating, rs[i].Rating)
			} else {
				rank = rank + 1
			}
		}

		if re.Distance != 0 {
			if re.Distance < rs[i].Distance {
				rank = rank + searching.GetIntegersDistance(re.Distance, rs[i].Distance)
			} else {
				rank = rank - searching.GetIntegersDistance(re.Distance, rs[i].Distance)
			}
		}

		if re.Price != 0 {
			if re.Price <= rs[i].Price {
				rank = rank + searching.GetIntegersDistance(re.Price, rs[i].Price)
			} else {
				rank = rank - 1
			}
		}

		rs[i].Rank = rank
	}
	sort.Sort(byRank(rs))
	return nil
}

func getTopFive(rs []Restaurant) []Restaurant {
	topFive := []Restaurant{}
	limit := len(rs)

	if limit < 4 {
		limit = len(rs)
	} else {
		limit = 4
	}
	for i := 0; i <= limit; i++ {
		topFive = append(topFive, rs[i])
	}
	return topFive
}
