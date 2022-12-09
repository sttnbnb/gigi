package bosyu

func getRoleName(source string) (name string) {
	name = source

	sourceRune := []rune(source)
	if len(sourceRune) > 9 {
		name = string(sourceRune[:9]) + "..."
	}

	return
}
