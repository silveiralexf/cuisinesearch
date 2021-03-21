package restaurants

import (
	"sync"
)

var mutex = &sync.Mutex{}

// PullCuisinesList return a list of cuisines provided by a given csv file
func PullCuisinesList(file string) (map[int]string, error) {
	rawCSVdata, err := readCSVFile(file, 2)
	if err != nil {
		return nil, err
	}

	cs := make(map[int]string)
	for i := 1; i < len(rawCSVdata); i++ {

		num, err := strToInt(rawCSVdata[i][0])
		if err != nil {
			return nil, err
		}

		mutex.Lock()
		cs[num] = rawCSVdata[i][1]
		mutex.Unlock()
	}
	return cs, nil
}

func getCuisineName(cs map[int]string, id int) string {
	if cs[id] == "" {
		return "unknown"
	}

	return cs[id]
}
