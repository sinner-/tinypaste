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
    "errors"
)

var validPaste = regexp.MustCompile("[A-Z0-9]{40}")

func getFile(url string) (string, error) {
    file := url[1:]
    if validPaste.MatchString(file) || file == "" || file == "index.html" {
        return file, nil
    }
    return "", errors.New("Paste not found.")
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
    file, err := getFile(r.URL.Path)
    if err != nil {
        log.Print(err.Error())
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    http.ServeFile(w, r, fmt.Sprintf("%s", file))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    paste := r.FormValue("paste")

    hasher := sha1.New()
    io.WriteString(hasher, fmt.Sprintf("%s%s", time.Now(), paste))
    hash := fmt.Sprintf("%X", hasher.Sum(nil))

    f, err := os.Create(fmt.Sprintf("%s", hash))
    if err != nil {
        log.Print(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    defer f.Close()

    _, err = f.WriteString(paste)
    if err != nil {
        log.Print(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    f.Sync()

    http.Redirect(w, r, fmt.Sprintf("%s", hash), http.StatusFound)
}

func main() {
    http.HandleFunc("/", staticHandler)
    http.HandleFunc("/save", saveHandler)
    http.ListenAndServe(":8080", nil)
}
