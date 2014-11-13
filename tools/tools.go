package tools

import (
    "os"
)


func InArray(a []string, e string) bool {
    for _, x := range a {
        if x == e {
            return true
        }
    }
    return false
}


func CheckFilePath(path string) bool {
    
    _, err := os.Stat(path);
    return !os.IsNotExist(err)
}
