package main

import (
    "fmt"
    "net/http"
    "log"
    "io"
    "os"
    "crypto/sha1"
    "time"
    "regexp"
)

var pasteDir = "pastes"
var validPaste = regexp.MustCompile("[A-Z0-9]{40}")

func getFile(url string) string {
    file := url[1:]
    if !validPaste.MatchString(file) {
        return "index.html"
    }
    return fmt.Sprintf("%s/%s", pasteDir, file)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
    file := getFile(r.URL.Path)
    http.ServeFile(w, r, fmt.Sprintf("%s", file))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    paste := r.FormValue("paste")

    hasher := sha1.New()
    io.WriteString(hasher, fmt.Sprintf("%s%s", time.Now(), paste))
    hash := fmt.Sprintf("%X", hasher.Sum(nil))

    f, err := os.Create(fmt.Sprintf("%s/%s", pasteDir, hash))
    if err != nil {
        log.Print(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer f.Close()

    _, err = f.WriteString(paste)
    if err != nil {
        log.Print(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    f.Sync()

    http.Redirect(w, r, fmt.Sprintf("%s", hash), http.StatusFound)
}

func main() {
    if _, err := os.Stat(pasteDir); err != nil {
        fmt.Println(fmt.Sprintf("Paste directory '%s' is not accessible or doesn't exist.", pasteDir))
        return
    }

    http.HandleFunc("/", staticHandler)
    http.HandleFunc("/save", saveHandler)
    http.ListenAndServe(":8080", nil)
}
