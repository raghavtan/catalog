package utils

import "github.com/motain/of-catalog/internal/modules/scorecard/dtos"

func DetectDrifts(
	stateList,
	configList []*dtos.Criterion,
	getUniqueKey func(*dtos.Criterion) string,
	getID func(*dtos.Criterion) string,
	setID func(*dtos.Criterion, string),
	isEqual func(*dtos.Criterion, *dtos.Criterion) bool,
) (created, updated, deleted, unchanged []*dtos.Criterion) {
	var createdList, updatedList, deletedList, unchangedList []*dtos.Criterion

	stateListMap := make(map[string]*dtos.Criterion)
	configListMap := make(map[string]*dtos.Criterion)

	for _, stateItem := range stateList {
		stateListMap[getUniqueKey(stateItem)] = stateItem
	}

	for _, configItem := range configList {
		configListMap[getUniqueKey(configItem)] = configItem
	}

	for name, stateItem := range stateListMap {
		configItem, found := configListMap[name]
		if !found {
			deletedList = append(deletedList, stateItem)
			continue
		}
		setID(configItem, getID(stateItem))
		if isEqual(stateItem, configItem) {
			unchangedList = append(unchangedList, configItem)
			continue
		}
		setID(configItem, getID(stateItem))
		updatedList = append(updatedList, configItem)
	}

	for name, configItem := range configListMap {
		if _, found := stateListMap[name]; !found {
			createdList = append(createdList, configItem)
		}
	}

	return createdList, updatedList, deletedList, unchangedList
}
