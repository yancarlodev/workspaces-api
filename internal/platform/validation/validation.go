package validation

type ValidationErrors map[string][]string

func (e ValidationErrors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e ValidationErrors) IsEmpty() bool {
	return len(e) == 0
}
