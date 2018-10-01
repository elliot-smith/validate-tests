package main

import "fmt"
import "strings"
import "os/exec"
import "os"

func main() {
    testCommand := os.Args[1]
    directory := os.Args[2]

    _, err := runTests(testCommand, directory)

    if ( err != nil) {
        fmt.Println("Tests did not pass without deleted lines")
        os.Exit(1)
    }

    fmt.Println("Tests passed")
}

func runTests(testCommand string, directory string) (string, error) {
    splitTestCommand := strings.Split(testCommand, " ")
    cmd := exec.Command(splitTestCommand[0])
    cmd.Args = splitTestCommand
    cmd.Dir = directory
    out, error := cmd.Output()
    return string(out), error
}