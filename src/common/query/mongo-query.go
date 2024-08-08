package query

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func SearchParser(search string) bson.M {
	searchSplit := strings.Split(search, ";")
	searchParsed := bson.M{}
	for i, s := range searchSplit {
		searchSplit := strings.Split(s, ":")
		// fmt.Println(123, searchSplit)
		if len(searchSplit) > 0 {
			if (i >0) {
				
			} else {
				searchParsed[string(search[0])] = bson.D{{Key: searchSplit[1], Value: searchSplit[2]}, {Key: "$options", Value: "i"}}
			}
			// return bson.M{searchSplit[0]: bson.D{{Key: searchSplit[1], Value: searchSplit[2]}, {Key: "$options", Value: "i"}}}
		} else {
			fmt.Printf("Invalid search query: %s\n", s)
			return searchParsed
		}
	}
	return searchParsed
}
