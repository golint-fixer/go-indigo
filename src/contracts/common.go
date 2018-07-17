package contracts

func (variable *ContractVariable) greater() bool {
	if variable.Input.(float64) > variable.Output.(float64) {
		return true
	}
	return false
}

func (variable *ContractVariable) greaterEqual() bool {
	if variable.Input.(float64) >= variable.Output.(float64) {
		return true
	}
	return false
}

func (variable *ContractVariable) lessthanEqual() bool {
	if variable.Input.(float64) <= variable.Output.(float64) {
		return true
	}
	return false
}

func (variable *ContractVariable) lessthan() bool {
	if variable.Input.(float64) < variable.Output.(float64) {
		return true
	}
	return false
}

func (variable *ContractVariable) equalEqual() bool {
	if variable.Input == variable.Output {
		return true
	}
	return false
}
