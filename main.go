package main

import (
    "fmt"
    "strings"
    "os/exec"
    "os"
    "io/ioutil"
    "path/filepath"
)

func main() {
    testCommand := os.Args[1]
    directory := os.Args[2]
    testFilePattern := os.Args[3]

    matches, error := filepath.Glob(directory + "/" + testFilePattern)
    if ( error != nil || len(matches) == 0 ) {
        fmt.Println("An error occurred trying to find any files with the pattern", testFilePattern)
        os.Exit(1)
    }

    fmt.Println(matches)

    testFile := matches[0]

    // Validate that the tests can run successfully
    _, err := runTests(testCommand, directory)
    if ( err != nil) {
        fmt.Println("Tests did not pass without deleted lines")
        os.Exit(1)
    }

    filesText, err := readAndBackupFile(testFile)

    fmt.Println(filesText)


    err = validateTests(testFile, testCommand, directory, filesText)

    if ( err != nil) {
        // TODO fix the following error
        fmt.Println("The following error occurred while trying to validate the tests: ")
        os.Exit(1)
    }

    err = restoreSystem(testFile)

    if ( err != nil) {
        fmt.Println("Could not restore the system to it's former state")
        os.Exit(1)
    }
}

func readAndBackupFile(fileNameAndDirectory string) (string, error) {
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
    fmt.Println(nextStatement)

    testPercent := (len(parsedText) * 100) / (len(parsedText) + len(remainingText))
    // TODO Identify why %v or %#v doesn't work below...
    fmt.Println("Parsed through %#v of file", testPercent)

    if(newRemainingText != "") {

        err := ioutil.WriteFile(fileNameAndDirectory, []byte(parsedText + newRemainingText), 0644)
        _, err = runTests(testCommand, directory)

        if ( err == nil ) {
            fmt.Println("Tests passed with the following deleted line %#v", nextStatement)
        }

        parseAndValidateTestFile(fileNameAndDirectory, testCommand, directory, parsedText + nextStatement, newRemainingText)
    }

    return nil
}

func getNextStatement (remainingText string) (string, string) {
    isNextStatement := false
    isValidStatement := false
    for index, character := range remainingText {
        switch character {
            case ';':
            case '\n':
                isNextStatement = true
                break
            case ' ':
            case '\r':
                break
            default:
                isValidStatement = true
        }

        if (isNextStatement == true && isValidStatement == true) {
            return remainingText[index+1:], remainingText[:index+1]
        }

        isNextStatement = false
    }

    return "", remainingText
}

func restoreSystem(fileNameAndDirectory string) (error) {
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
