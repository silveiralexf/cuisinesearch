[![Go Report Card](https://goreportcard.com/badge/github.com/silveiralexf/cuisinesearch)](https://goreportcard.com/report/github.com/silveiralexf/cuisinesearch)

# Find the best-matched restaurants challenge

You have data about local restaurants located near your company, which you can find in the **restaurants.csv** file. You would like to develop a basic search function that allows your colleagues to search those restaurants to help them find where they would like to have lunch. The search is based on five criteria: **Restaurant Name, Customer Rating(1 star ~ 5 stars), Distance(1 mile ~ 10 miles), Price(how much one person will spend on average, $10 ~ $50), Cuisine(Chinese, American, Thai, etc.).** The requirements are listed below.

1. The function should allow users to provide up to five parameters based on the criteria listed above. *You can assume each parameter can contain only one value.*
2. If parameter values are invalid, return an error message.
3. The function should return up to five matches based on the provided criteria. If no matches are found, return an empty list. If less than 5 matches are found, return them all. If more than 5 matches are found, return the best 5 matches. The returned results should be sorted according to the rules explained below. Every record in the search results should at least contain the restaurant name.
4. "Best match" is defined as below:
   - A Restaurant Name match is defined as an exact or partial String match with what users provided. For example, "Mcd" would match "Mcdonald's".
   - A Customer Rating match is defined as a Customer Rating equal to or more than what users have asked for. For example, "3" would match all the 3 stars restaurants plus all the 4 stars and 5 stars restaurants.
   - A Distance match is defined as a Distance equal to or less than what users have asked for. For example, "2" would match any distance that is equal to or less than 2 miles from your company.
   - A Price match is defined as a Price equal to or less than what users have asked for. For example, "15" would match any price that is equal to or less than $15 per person.
   - A Cuisine match is defined as an exact or partial String match with what users provided. For example, "Chi" would match "Chinese". You can find all the possible Cuisines in the **cuisines.csv** file. *You can assume each restaurant offers only one cuisine.*
   - The five parameters are holding an "AND" relationship. For example, if users provide Name = "Mcdonald's" and Distance = 2, you should find all "Mcdonald's" within 2 miles.
   - When multiple matches are found, you should sort them as described below.
     - Sort the restaurants by Distance first.
     - After the above process, if two matches are still equal, then the restaurant with a higher customer rating wins.
     - After the above process, if two matches are still equal, then the restaurant with a lower price wins.
     - After the above process, if two matches are still equal, then you can randomly decide the order.
     - Example: if the input is Customer Rating = 3 and Price = 15. Mcdonald's is 4 stars with an average spend = $10, and it is 1 mile away. And KFC is 3 stars with an average spend = $8, and it is 1 mile away. Then we should consider Mcdonald's as a better match than KFC. (They both matches the search criteria -> we compare distance -> we get a tie -> we then compare customer rating -> Mcdonald's wins)
5. The final submitted work should include a README file. No UI is required in this assessment, but you may implement one if you would like. **The steps to run and test your program should be clearly introduced in the README file.** If you have made any additional **Assumptions** besides what we have listed above while working on this assessment, please document them so that we can better understand your solution.


## Solution

### 1. Assumptions

1. Goal of the challenge is to evaluate logic and approach of the problem, in a real scenario would begin by first by using a different type of persistence than CSV files.
2. Structure of files and column order won't change.
3. Customers rating will always range from 1 to 5.
4. Hard-coding values such as file path and default port should not be important, in a real scenario would read from different sources such as flags, environment varialbes, etc.


## 2. Setup

1. Clone the repository into your GOPATH:

```shell
git clone git@github.com:silveiralexf/cuisinesearch
```

For instructions on how setup your GO environment check the links below:

- [Install](https://golang.org/doc/install)
- [GOPATH Setup](https://golang.org/doc/gopath_code)

2. Enter to the repository base directory and build as instructed below:

```shell
go build -trimpath .
```

A new binary should be available at the same directory as shown below:

```shell
$ ls -ltr cuisinesearch
-rwxrwxr-x 1 silveiralex silveiralex 6874833 Mar 20 23:50 cuisinesearch

```

## 3. Resources

### Get Full List of Restaurants

Retrieves full list of restaurants in `application/json` format.

```
GET /restaurants
```

### Search Top 5 Restaurants by Query

Retrieves top 5 restaurants by search criteria in `application/json` format.

```
GET /restaurants/search?<field>=<search_term>
```
#### Parameters

| Type  | Name     | Schema  | Example                               |
|-------|----------|---------|---------------------------------------|
| Query | name     | string  | `/restaurants/search?name=Palacex`    |
| Query | cuisine  | string  | `/restaurants/search?cuisine=Palacex` |
| Query | distance | integer | `/restaurants/search?distance=6`      |
| Query | price    | integer | `/restaurants/search?price=15`        |
| Query | rating   | integer | `/restaurants/search?rating=2`        |


## 4. Usage

1. Execute the compiled binary so that it can start listening and serving on port 8080 as shown below:

```shell
$ ./cuisinesearch
2021-03-20 23:51:07 [INFO] [Startup] listening and serving on port :8080
```

2. Access the server from your browser, or from command line utility such as `curl` on the endpoints previously described, as the example below:

```shell
$ curl -s "localhost:8080/restaurants/search?cuisine=Mexican&price=50&rating=3" | jq
[
  {
    "id": 159,
    "name": "Yummyscape",
    "rating": 1,
    "distance": 3,
    "price": 35,
    "cuisine_id": 13,
    "cuisine": "Mexican",
    "rank": 1
  },
  {
    "id": 114,
    "name": "Grillarc",
    "rating": 2,
    "distance": 3,
    "price": 25,
    "cuisine_id": 13,
    "cuisine": "Mexican",
    "rank": 1
  },
  {
    "id": 41,
    "name": "Piece Chow",
    "rating": 4,
    "distance": 9,
    "price": 10,
    "cuisine_id": 13,
    "cuisine": "Mexican",
    "rank": 1
  },
  {
    "id": 157,
    "name": "Divine Yummy",
    "rating": 1,
    "distance": 10,
    "price": 25,
    "cuisine_id": 13,
    "cuisine": "Mexican",
    "rank": 1
  },
  {
    "id": 55,
    "name": "Tablebes",
    "rating": 4,
    "distance": 2,
    "price": 40,
    "cuisine_id": 13,
    "cuisine": "Mexican",
    "rank": 1
  }
]
```

**Note**: Piping the command output to a JSON formatter such as [JQ](https://stedolan.github.io/jq/download/) for a formatted response is optional, else the results will be returned as a single line JSON.

3. On the console output, or on the logs created at the same directory, you should receive similar responses for the example requests:

```shell
2021-03-21 00:16:28 [INFO] [Startup] listening and serving on port :8080
2021-03-21 00:16:29 [INFO] [localhost:8080] HTTP status 200 on '/restaurants/search?name=CutsDeliciouns'
2021-03-21 00:16:51 [ERROR] [localhost:8080] endpoint '/restaurants/invalid' could not be found [server.go.restaurants.notFound:33]
2021-03-21 00:17:45 [INFO] [localhost:8080] HTTP status 200 on '/restaurants/search?distance=5'
2021-03-21 00:17:58 [INFO] [localhost:8080] HTTP status 200 on '/restaurants/search?distance=5&name=Yummylia'
```

4. In case CSV files cannot be found at `files` directory, following should be returned:

```shell
$ curl -s "localhost:8080/restaurants/search?cuisine=Mexican&price=50&rating=3" | jq
{
  "method": "GET",
  "status_code": "500",
  "host": "localhost:8080",
  "error": "failed to open csv file 'files/cuisines.csv': [open files/cuisines.csv: no such file or directory]",
  "filterspec": "{[restaurants restaurants/search]}"
}
```

5. In case of invalid endpoint is informed, a similar response is expected:

```shell
$ curl -s "localhost:8080/invalid" | jq
{
  "method": "GET",
  "status_code": "404",
  "host": "localhost:8080",
  "error": "endpoint '/invalid' could not be found",
  "filterspec": "{[restaurants restaurants/search]}"
}
```