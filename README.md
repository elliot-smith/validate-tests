# validate-tests
A tool that checks whether you code impact your test output

## Reasoning behind the validate tests 
Code coverage tools are amazing at identifying spots that you haven't tested in some manner very quickly but all have one flaw. This is that the output/side effects/behaviour of your code aren't actually matched to your tests explicitly. What this means is that you can have code that is shown as covered but realistically it had zero impact on your test! 
(insert example here) 
Example: Not testing all effects of the function
```
let globalValue = 5
let secondGlobalValue = 7

function newGlobalValues(newValue1, newValue2) {
    globalValue = newValue1
    secondGlobalValue = newValue2
} 

function testGlobalValues() {
    newGlobalValues(1, 2)
    expeect(globalValue).to.equal(1)
} 
```

## Commands
validate-tests.exe "INSERT_TEST_COMMAND" "INSERT_TEST_DIRECTORY_HERE" "INSERT_TEST_FILE_HERE" "TEXT_TO_REPLACE_FOR_TEST_ISOLATION" "REPLACE_TEXT_FOR_TEST_ISOLATION" "TEST_FILE_ENDING"
go build && validate-tests.exe "npm run test" "F:/Developer/workspace/example-typescript-nyc-mocha-coverage" "src/**/add.ts" "describe(" "describe.only(" ".test.ts"


Argument Name | Description
--- | ---
"INSERT_TEST_COMMAND" | The command that tests your code. This could be `npm run test` or `mvn clean package` for example
"INSERT_TEST_DIRECTORY_HERE" | The directory that your project sits in. This can either be relative to the validate-test file on your machine or the absolute path
"INSERT_TEST_FILE_HERE" | The regex that will be used to determine all files that you would like tested.
"TEXT_TO_REPLACE_FOR_TEST_ISOLATION" | A string in tests files that can be modified to isolate tests to that file
"REPLACE_TEXT_FOR_TEST_ISOLATION" | A string to put in the test file (over the one above) to isolate the tests
"TEST_FILE_ENDING" | The ending to all of your test files (i.e. "spec.js", "test.js").