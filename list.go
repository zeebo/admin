package admin

import (
	"launchpad.net/mgo"
	"net/url"
	"strconv"
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
	page, numpage := grabInt(v, "page", 0), grabInt(v, "numpage", 20)

	if page == 0 {
		page = 1
	}

	query := c.Find(nil).Skip(numpage * (page - 1)).Limit(numpage)
	return query.Iter()
}
