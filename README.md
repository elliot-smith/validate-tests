# validate-tests
A tool that checks whether you code impact your test output



## Commands
validate-tests.exe "INSERT_TEST_COMMAND" "INSERT_TEST_DIRECTORY_HERE" "INSERT_TEST_FILE_HERE"
go build && validate-tests.exe "npm run test" "F:/Developer/workspace/example-typescript-nyc-mocha-coverage" "calculator/ts/src/add.ts"