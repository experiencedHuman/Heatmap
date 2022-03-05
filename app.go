package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
)


func main() {
    // open file
    f, err := os.Open("download.csv")
    if err != nil {
        log.Fatal(err)
    }

    // remember to close the file at the end of the program
    defer f.Close()

    // read csv values using csv.Reader
    csvReader := csv.NewReader(f)
    for i := 0; i < 20; i++ {
        rec, err := csvReader.Read()

        if err != nil {
            log.Fatal(err)
        }
        // do something with read line
        fmt.Printf("%+v\n", rec)
    }
}
