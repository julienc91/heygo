package tools

func InArray(a []string, e string) bool {
    for _, x := range a {
        if x == e {
            return true
        }
    }
    return false
}
