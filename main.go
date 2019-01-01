package main

import (
    "fmt"
    "strings"
    "os/exec"
    "os"
    "io/ioutil"
    "github.com/bmatcuk/doublestar"
	"math/rand"
	"time"
)

func main() {
    testCommand := os.Args[1]
    directory := os.Args[2]
    testFilePattern := os.Args[3]
    testFileOriginalText := os.Args[4]
    testFileReplaceText := os.Args[5]
    endTestExtension := os.Args[6]

    matches, error := doublestar.Glob(directory + "/" + testFilePattern)
    if ( error != nil || len(matches) == 0 ) {
        fmt.Println("An error occurred trying to find any files with the pattern", testFilePattern)
        os.Exit(1)
    }

    fmt.Println("The matches are: ", matches)


    randomSeed := rand.NewSource(time.Now().UnixNano())
    randomGenerator := rand.New(randomSeed)

    testFile := matches[randomGenerator.Intn(len(matches))]

    validateCurrentCode(testCommand, directory)

    backupAndUpdateTestFile(testFile, endTestExtension, testFileOriginalText, testFileReplaceText)

    filesText := readAndBackupFile(testFile)

    fmt.Println(filesText)


    err := validateTests(testFile, testCommand, directory, filesText)

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

func validateCurrentCode(testCommand string, directory string) {
    _, err := runTests(testCommand, directory)
    if ( err != nil) {
        fmt.Println("Tests did not pass without deleted lines")
        os.Exit(1)
    }
}

func backupAndUpdateTestFile(fileNameAndDirectory string, endTestExtension string, testFileOriginalText string, testFileReplaceText string) () {
    testFileAndDirectory := fileNameAndDirectory[:strings.LastIndex(fileNameAndDirectory, ".")] + endTestExtension

    filesText := readAndBackupFile(testFileAndDirectory)

    filesText = strings.Replace(filesText, testFileOriginalText, testFileReplaceText, -1)

    writeFileErr := ioutil.WriteFile(testFileAndDirectory, []byte(filesText), 0444)

    if (writeFileErr != nil) {
        fmt.Println("Unable to create the backup file for file" + fileNameAndDirectory)
        os.Exit(1)
    }
}

func readAndBackupFile(fileNameAndDirectory string) (string) {
    filesText, err := ioutil.ReadFile(fileNameAndDirectory)

    if (err != nil) {
        fmt.Println("Unable to read the file" + fileNameAndDirectory)
        os.Exit(1)
    }

    backupTestFileName := fileNameAndDirectory + ".backup"
    writeFileErr := ioutil.WriteFile(backupTestFileName, filesText, 0444)

    if (writeFileErr != nil) {
        fmt.Println("Unable to create the backup file for file" + fileNameAndDirectory)
        os.Exit(1)
    }

    fileString := strings.Replace(string(filesText), `\n`, "\n", -1)

    return fileString
}

func validateTests(fileNameAndDirectory string, testCommand string, directory string, filesText string) (error) {
    err := parseAndValidateTestFile(fileNameAndDirectory, testCommand, directory, "", filesText)

    if (err != nil) {
        return err
    }

    return nil
}

func parseAndValidateTestFile(fileNameAndDirectory string, testCommand string, directory string, parsedText string, remainingText string) (error) {
    newRemainingText, nextStatement := getNextStatement(remainingText)

    testPercent := (len(parsedText) * 100) / (len(parsedText) + len(remainingText))
    // TODO Identify why %v or %#v doesn't work below...
    fmt.Println("Parsed through %#v of file", testPercent)


    err := ioutil.WriteFile(fileNameAndDirectory, []byte(parsedText + newRemainingText), 0444)

    if ( err != nil ) {
        fmt.Println("Error overwriting the file", err)
    } else {
        _, err = runTests(testCommand, directory)

        if ( err == nil ) {
            fmt.Println("Tests passed with the following deleted line %#v", nextStatement)
        }
    }

    if(newRemainingText != "") {
        return parseAndValidateTestFile(fileNameAndDirectory, testCommand, directory, parsedText + nextStatement, newRemainingText)
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
