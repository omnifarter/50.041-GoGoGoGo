package helpers

import nodes "gogogogo/nodes"

func GetLatestDatabaseEntryValue(data []nodes.DatabaseEntry) int {
	latestClock := -1
	var value int

	for _, v := range data {
		if v.Clock > latestClock {
			latestClock = v.Clock
			value = v.Value
		}
	}

	return value
}
