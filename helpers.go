package main

func StringArrayContains(alpha []string, str string) bool {

	// iterate using the for loop
	for i := 0; i < len(alpha); i++ {
		// check
		if alpha[i] == str {
			// return true
			return true
		}
	}
	return false
}
