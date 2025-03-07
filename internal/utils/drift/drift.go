package drift

func Detect[T any](
	stateMap, configMap map[string]*T,
	fromStateToConfig func(state *T, conf *T),
	isEqual func(*T, *T) bool,
) (created, updated, deleted, unchanged map[string]*T) {
	createdList := make(map[string]*T)
	updatedList := make(map[string]*T)
	deletedList := make(map[string]*T)
	unchangedList := make(map[string]*T)

	processStateItems(stateMap, configMap, fromStateToConfig, isEqual, updatedList, deletedList, unchangedList)
	processConfigItems(stateMap, configMap, createdList)

	return createdList, updatedList, deletedList, unchangedList
}

func processStateItems[T any](
	stateMap, configMap map[string]*T,
	fromStateToConfig func(state *T, conf *T),
	isEqual func(*T, *T) bool,
	updatedList, deletedList, unchangedList map[string]*T,
) {
	for key, stateItem := range stateMap {
		configItem, found := configMap[key]
		if !found {
			deletedList[key] = stateItem
			continue
		}

		fromStateToConfig(stateItem, configItem)
		if isEqual(stateItem, configItem) {
			unchangedList[key] = configItem
			continue
		}

		updatedList[key] = configItem
	}
}

func processConfigItems[T any](stateMap, configMap map[string]*T, createdList map[string]*T) {
	for key, configItem := range configMap {
		if _, found := stateMap[key]; !found {
			createdList[key] = configItem
		}
	}
}
