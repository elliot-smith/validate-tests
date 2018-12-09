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


    err = validateTests(testFile, testCommand, directory, filesText)

    if ( err != nil) {
        // TODO fix the following error
        fmt.Println("The following error occurred while trying to validate the tests: ")
        os.Exit(1)
    }

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

    fileString := strings.Replace(string(filesText), `\n`, "\n", -1)

    return fileString, nil
}

func validateTests(testFile string, testCommand string, directory string, filesText string) (error) {
    fileNameAndDirectory := directory + "/" + testFile
    err := parseAndValidateTestFile(fileNameAndDirectory, testCommand, directory, "", filesText)

    if (err != nil) {
        return err
    }

    return nil
}

func parseAndValidateTestFile(fileNameAndDirectory string, testCommand string, directory string, parsedText string, remainingText string) (error) {
    newRemainingText, nextStatement := getNextStatement(remainingText)

    testPercent := (len(parsedText) * 100) / (len(parsedText) + len(remainingText))
    fmt.Println("Parsed through %d of file", testPercent)

    if(newRemainingText != "") {

        err := ioutil.WriteFile(fileNameAndDirectory, []byte(parsedText + newRemainingText), 0644)
        _, err = runTests(testCommand, directory)

        if ( err == nil ) {
            fmt.Println("Tests passed with the following deleted line %s", nextStatement)
        }

        parseAndValidateTestFile(fileNameAndDirectory, testCommand, directory, parsedText + nextStatement, newRemainingText)
    }

    return nil
}

func getNextStatement (remainingText string) (string, string) {
    return getNextStatementRecursive(remainingText, "")
}

func getNextStatementRecursive (remainingText string, nextStatement string) (string, string) {
    if (remainingText == "") {
        return remainingText, ""
    }

    isTerminatingCharacters, terminatingCharacters := isTerminatingCharacterSet(remainingText)

    if(isTerminatingCharacters) {
        return remainingText[len(terminatingCharacters):], nextStatement + terminatingCharacters
    }

    fmt.Println(nextStatement, terminatingCharacters)

    return getNextStatementRecursive(remainingText[len(terminatingCharacters):], nextStatement + terminatingCharacters)
}

func isTerminatingCharacterSet (remainingText string) (bool, string) {
    terminatingString := ""
    switch remainingText[:1] {
        case ";":
            terminatingString = ";"
        case "/":
            switch remainingText[1:2] {
                case "n":
                    terminatingString = "/n"
            }
    }

    if (terminatingString == "") {
       return false, remainingText[:1]
    }
    return true, terminatingString
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