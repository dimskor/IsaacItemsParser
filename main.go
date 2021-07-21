package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "IsaacItemsParser/parser"
    "strings"
)

func main() {
    fmt.Println("Starting...")

    items:= parser.Parse()

    fmt.Printf("Parsed total: %d\n", len(items))

    jsonData, _ := json.MarshalIndent(items, "", "    ")
    _ = ioutil.WriteFile("items.json", jsonData, 0644)

    fmt.Println("Done!")
}
