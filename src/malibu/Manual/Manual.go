package Manual

import "fmt"

func PrintManual() {
	fmt.Println(`NAME
       combat - list tests in current directory

SYNOPSIS
       combat [ACTION]... [OPTIONS]... [PARAMETERS_FOR_TESTS]...

DESCRIPTION
       List information about tests conained in current directory.
       Behaviour by default - the same as "run" action without options.

ACTIONS
       help
              show this manual

       run
              run tests, return count of failed tests in exit-code

       list
              show list of tests, ordered by name

       tags
              show list of tests, ordered by tags

       params
              show list of tests, ordered by parameters

       cases
              show list of tests that will be run on "run" action.
              If some of required test parameters are not provided - needed parameters will be shown.


OPTIONS
       -name=TEST_NAME_REGEX
              select tests by name (regular expression). Tests with matched names will list or run, depends on action.
              Example: -name=anonymous.*

       -tag=TEST_TAG
              select tests by tag (able to listed with comma). Tests with matched tags will list or run, depends on action
              Example: -tag=NotForLive,paymentTest,longTest


PARAMETERS FOR TESTS
       There is parameters that will be provided to tests before run it.
       You can find out information about parameters of your tests by running combat with any list action (list,
       params, tags) "combat list".
       All test parameter's names is a sequence of characters without spaces. Only letters and numbers are accepted.
       All test parameter's values is a sequence of characters without spaces. Special symbols and numbers are accepted.
       There is two types of test parameters:

       stringParam
              Has no default value. It is must to be defined explicitly to run test.
              Example: -hostname=http://ProjectUAT.com

       enumParam
              String parameter that only accept values from the list. For example: Location(uk, us, fr,
              ru); Resolution(desktop, mobile)
              It is not must to be defined to run test.
              If it is defined explicitly - test will be run with each of all provided values.
              If it is not defined - test will be run with each of all accepted values.
              If the test has two or more enum parameters - test will be run with all combinations of parameters.
              For example (uk;desktop uk;mobile ru;desktop ru;mobile etc...)
              You can find out cases by "combat cases" command.

EXIT STATUS
       0                    all tests are passed, or list action executed successfully.
       error count          some tests fails, or some tests does not match the test format.


EXAMPLES (You are able run all following commands in "CombatSelfTesting" folder)
       combat
              run all tests in current directory, using all accepted cases

       combat list
              show list of tests with parameters and tags, ordered by name of test

       combat list -tag=NotForLive,adminPanelTest -name=.*Currency.*
              Show the list of tests, selected by following tags, and name contans "Currency"

       combat cases -tag=NotForLive,adminPanelTest -name=.*Currency.* -Locale=uk,us
              Show all cases for tests selected by name and tag, with following locales (uk,us)

REPORTING BUGS
       Report combat behaviour bugs to <alexander.eliseev@graph.uk>`)
}
