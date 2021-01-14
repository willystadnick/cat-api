package main

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"
)

func breedsSearch(name string) []Cat {
    url := "https://api.thecatapi.com/v1/breeds/search?q=" + name
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        panic(err)
    }
    req.Header.Add("x-api-key", os.Getenv("THECATAPI_KEY"))

    cli := &http.Client{}
    res, err := cli.Do(req)
    if err != nil {
        panic(err)
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err)
    }

    var cats []Cat
    json.Unmarshal(body, &cats)

    return cats
}
