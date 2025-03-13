package consensus

func GetSlotLeader(epoch int, seed int, registryKeys [][]byte, stakeData map[string]int) []byte {
	// For now, sequentially choose the slot leader from registryKeys, parameters are, stakedata
	stakeProbs := GetStakes(stakeData)
	cumulativeProb := 0.0

	for _, registry := range registryKeys {
		cumulativeProb += stakeProbs[string(registry)]
		if seed <= int(cumulativeProb) {
			return registry
		}
	}

	return registryKeys[len(registryKeys) - 1]
}

func GetStakes(stakeData map[string]int) map[string]float64 {
	sum := 0.0
	for _, numDomains := range stakeData {
		sum += float64(numDomains)
	}

	stakeProbs := make(map[string]float64)
	for registry, numDomains := range stakeData {
		stakeProbs[registry] = float64(numDomains) / sum
	}

	return stakeProbs
}
