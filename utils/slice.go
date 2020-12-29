package utils

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func Unpack(s []string, vars ...*string) {
	for i, str := range s {
		if vars[i] == nil {
			continue
		}
		*vars[i] = str
	}
}
