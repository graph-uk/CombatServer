package arrayUtils

// check is a string presented in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ConcatStringSlice(slice []string) string {
	result := ""
	for _, curElement := range slice {
		result += curElement
	}
	return result
}
