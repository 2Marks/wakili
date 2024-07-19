package proxy

import (
	"fmt"
	"os"
)

func PurgeCachedDataHandler(baseUrl string) {
	err := purgeCache(baseUrl)

	if err != nil {
		fmt.Printf("error occured while purging data. err: %s \n", err.Error())
	} else {
		fmt.Printf("data purged successfully \n")
	}

	os.Exit(0)
}
