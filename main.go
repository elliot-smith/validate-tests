package main

import (
    "fmt"
    "strings"
    "os/exec"
    "os"
    "io/ioutil"
)

func main() {
    testCommand := os.Args[1]
    directory := os.Args[2]
    testFile := os.Args[3]

    // Validate that the tests can run successfully
    _, err := runTests(testCommand, directory)
    if ( err != nil) {
        fmt.Println("Tests did not pass without deleted lines")
        os.Exit(1)
    }

    dat, err := ioutil.ReadFile(directory + "/" + testFile)

    fmt.Println(string(dat))

    fmt.Println("Tests passed")
    fmt.Println(directory + testFile)
}

func runTests(testCommand string, directory string) (string, error) {
    splitTestCommand := strings.Split(testCommand, " ")
    cmd := exec.Command(splitTestCommand[0])
    cmd.Args = splitTestCommand
    cmd.Dir = directory
    out, error := cmd.Output()
    return string(out), error
}