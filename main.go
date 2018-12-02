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

    filesText, err := readAndBackupFile(directory, testFile)

    fmt.Println(filesText)

    fmt.Println("Tests passed")
    fmt.Println(directory + testFile)

    err = restoreSystem(directory, testFile)

    if ( err != nil) {
        fmt.Println("Could not restore the system to it's former state")
        os.Exit(1)
    }
}

func readAndBackupFile(directory string, testFile string) (string, error) {
    fileNameAndDirectory := directory + "/" + testFile
    filesText, err := ioutil.ReadFile(fileNameAndDirectory)

    if (err != nil) {
        fmt.Println("Unable to read the file" + fileNameAndDirectory)
        return "", fmt.Errorf("Unable to read the file")
    }

    backupTestFileName := fileNameAndDirectory + ".backup"
    writeFileErr := ioutil.WriteFile(backupTestFileName, filesText, 0644)

    if (writeFileErr != nil) {
        fmt.Println("Unable to create the backup file for file" + fileNameAndDirectory)
        return "", fmt.Errorf("Unable to create the backup file")
    }

    return string(filesText), nil
}

func restoreSystem(directory string, testFile string) (error) {
    fileNameAndDirectory := directory + "/" + testFile
    filesText, err := ioutil.ReadFile(fileNameAndDirectory + ".backup")

    if (err != nil) {
        fmt.Println("Unable to read the backup file" + fileNameAndDirectory + ".backup")
        return fmt.Errorf("Unable to read the file")
    }

    err = ioutil.WriteFile(fileNameAndDirectory, filesText, 0644)

    if (err != nil) {
        fmt.Println("Unable to overwrite the original file" + fileNameAndDirectory)
        return fmt.Errorf("Unable to create the backup file")
    }

    err = os.Remove(fileNameAndDirectory + ".backup")

    if (err != nil) {
        fmt.Println("Unable to delete the backup file" + fileNameAndDirectory + ".backup")
        return fmt.Errorf("Unable to create the delete file")
    }

    return nil
}

func runTests(testCommand string, directory string) (string, error) {
    splitTestCommand := strings.Split(testCommand, " ")
    cmd := exec.Command(splitTestCommand[0])
    cmd.Args = splitTestCommand
    cmd.Dir = directory
    out, error := cmd.Output()
    return string(out), error
}