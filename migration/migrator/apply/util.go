package apply

import (
	ai "github.com/amirkode/go-mongr8/migration/translator/mongodb/api_interpreter"
)

func getLatestMigrationID() string {
	// TODO: implement something to get latest id from db
	return ""
}

func filterSubActionApi(apis []ai.SubActionApi) []ai.SubActionApi {
	res := []ai.SubActionApi{}
	latestMigrationID := getLatestMigrationID()
	for _, api := range apis {
		if api.MigrationID > latestMigrationID {
			res = append(res, api)
		}
	}

	return res
}