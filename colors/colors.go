package colors

func Reverse(str []string) []string {
	vs := make([]string, len(str))
	n := copy(vs, str)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		vs[i], vs[j] = vs[j], vs[i]
	}
	return vs
}
