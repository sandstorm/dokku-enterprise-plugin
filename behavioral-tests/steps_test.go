package main

import "github.com/DATA-DOG/godog"

func FeatureContext(s *godog.Suite) {

	// steps_configuration_test
	s.Step(`^the configuration is:$`, theConfigurationIs)
	s.Step(`^the cloud configuration is:$`, theCloudConfigurationIs)

	// steps_dokku_test
	s.Step(`^I have an empty Dockerfile application$`, iHaveAnEmptyDockerfileApplication)
	s.Step(`^I create the file "([^"]*)" with the following contents:$`, iCreateTheFileWithTheFollowingContents)
	s.Step(`^I remove the file "([^"]*)"$`, iRemoveTheFile)
	s.Step(`^I deploy the application as "([^"]*)"$`, iDeployTheApplicationAs)
	s.Step(`^I call dokku "([^"]*)"$`, iCallDokku)
	s.Step(`^I call dokku "([^"]*)" with payload:$`, iCallDokkuWithPayload)
	s.Step(`^I call the URL "([^"]*)" of the "([^"]*)" application$`, iCallTheURLOfTheApplication)
	s.Step(`^the response should contain "([^"]*)"$`, theResponseShouldContain)
	s.Step(`^the response should not contain "([^"]*)"$`, theResponseShouldNotContain)
	s.Step(`^I get back a JSON object with the following structure:$`, iGetBackAJSONObjectWithTheFollowingStructure)
	s.Step(`^I get back a message "([^"]*)"$`, iGetBackAMessage)
	s.Step(`^I get back a table with content:$`, iGetBackATableWithContent)

	// steps_eventLog_test
	s.Step(`^the event log is empty$`, theEventLogIsEmpty)
	s.Step(`^I expect (\d+) event log entr(?:y|ies) on disk$`, iExpectEventLogEntry)

	// steps_files_test
	s.Step(`^an empty folder "([^"]*)" exists$`, anEmptyFolderExists)
	s.Step(`^I expect a file "([^"]*)" in folder "([^"]*)"$`, iExpectAFileInFolder)

	// steps_mockHttpServer_test
	s.Step(`^the API delivery http server is disabled$`, theAPIDeliveryHttpServerIsDisabled)
	s.Step(`^the API delivery http server is available at port (\d+) for at most (\d+) seconds and (\d+) request$`, theAPIDeliveryHttpServerIsAvailableAt)
	s.Step(`^the API delivery http server always responds with status code (\d+)$`, theAPIDeliveryHttpServerAlwaysRespondsWithStatusCode)
	s.Step(`^the API delivery http server received request (\d+) with the following JSON at "([^"]*)":$`, theAPIDeliveryHttpServerReceivedTheFollowingJSONAtEvent)


	// todo: sort them into the respective category
	s.Step(`^a file "([^"]*)" is created with contents:$`, aFileIsCreatedWithContents)
	s.Step(`^I execute the following SQL statements on database "([^"]*)":$`, iExecuteTheFollowingSQLStatementsOnDatabase)
	s.Step(`^I expect a file "([^"]*)" with contents:$`, iExpectAFileWithContents)
	s.Step(`^the SQL statement "([^"]*)" on database "([^"]*)" must return "([^"]*)"$`, theSQLStatementOnDatabaseMustReturn)
}
