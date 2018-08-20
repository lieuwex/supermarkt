package utils

func UniqStrings(slice []string) []string {
	m := make(map[string]bool)
	var res []string

	for _, str := range slice {
		if !m[str] {
			m[str] = true
			res = append(res, str)
		}
	}

	return res
}
