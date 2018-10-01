package main

import "fmt"
import "os/exec"
import "os"

func main() {
    _, err := run_tests()

    if ( err != nil) {
        fmt.Println("Tests did not pass without deleted lines")
        os.Exit(1)
    }

    fmt.Println("Tests passed")
}

func run_tests() (string, error) {
    cmd := exec.Command("npm", "test");
    cmd.Dir = "F:/Developer/workspace/example-typescript-nyc-mocha-coverage"
    out, error := cmd.Output()
    return string(out), error
}