package admin

import (
	"launchpad.net/mgo"
	"net/url"
	"strconv"
	"strings"
)

//grabInt is a simple handler for grabbing an integer out of a url.Values
//with a default value if there are any errors
func grabInt(v url.Values, key string, def int) int {
	val := v.Get(key)
	n, err := strconv.ParseInt(val, 10, 0)
	if err != nil || n < 0 {
		return def
	}
	return int(n)
}

//listParse takes a collection and some query values and generates an iterator
//for the objects that should be returned on that page
func listParse(c mgo.Collection, v url.Values) *mgo.Iter {
	//parse out sorting
	sort := map[string]int{}
	for key, _ := range v {
		val := v.Get(key)
		if len(key) > 5 && strings.HasPrefix(key, "sort_") {
			field := key[5:]
			switch strings.ToLower(val) {
			case "asc":
				sort[field] = 1
			case "desc":
				sort[field] = -1
			}
		}
	}
	//set up the query with the correct sort order
	query := c.Find(nil).Sort(sort)

	//pagination
	page, numpage := grabInt(v, "page", 0), grabInt(v, "numpage", 20)
	//set page 0 to be page 1
	if page == 0 {
		page = 1
	}
	//pages are 1 indexed.
	query = query.Skip(numpage * (page - 1)).Limit(numpage)

	return query.Iter()
}
