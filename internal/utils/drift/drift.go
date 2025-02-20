package drift

func Detect[T any](
	stateList,
	configList []*T,
	getUniqueKey func(*T) string,
	getID func(*T) string,
	setID func(*T, string),
	isEqual func(*T, *T) bool,
) (created, updated, deleted, unchanged []*T) {
	var createdList, updatedList, deletedList, unchangedList []*T

	stateListMap := make(map[string]*T)
	configListMap := make(map[string]*T)

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
