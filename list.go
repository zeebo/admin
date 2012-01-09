package admin

import (
	"launchpad.net/mgo"
	"net/url"
	"strconv"
	"strings"
)

func grabInt(v url.Values, key string, def int) int {
	val := v.Get(key)
	n, err := strconv.ParseInt(val, 10, 0)
	if err != nil || n < 0 {
		return def
	}
	return int(n)
}

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
	query := c.Find(nil).Sort(sort)

	//pagination
	page, numpage := grabInt(v, "page", 0), grabInt(v, "numpage", 20)
	if page == 0 {
		page = 1
	}
	query = query.Skip(numpage * (page - 1)).Limit(numpage)

	return query.Iter()
}
