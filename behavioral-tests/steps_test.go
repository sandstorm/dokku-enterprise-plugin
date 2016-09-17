package main

import "github.com/DATA-DOG/godog"

func FeatureContext(s *godog.Suite) {
	s.Step(`^I remove the file "([^"]*)"$`, iRemoveTheFile)
	s.Step(`^I have an empty Dockerfile application$`, iHaveAnEmptyDockerfileApplication)
	s.Step(`^I create the file "([^"]*)" with the following contents:$`, iCreateTheFileWithTheFollowingContents)
	s.Step(`^I deploy the application as "([^"]*)"$`, iDeployTheApplicationAs)
	s.Step(`^I call the URL "([^"]*)" of the "([^"]*)" application$`, iCallTheURLOfTheApplication)
	s.Step(`^the response should contain "([^"]*)"$`, theResponseShouldContain)
	s.Step(`^the response should not contain "([^"]*)"$`, theResponseShouldNotContain)
	s.Step(`^the event log is empty$`, theEventLogIsEmpty)
	s.Step(`^I expect (\d+) event log entr(?:y|ies) on disk$`, iExpectEventLogEntry)
	s.Step(`^the API delivery http server is disabled$`, theAPIDeliveryHttpServerIsDisabled)
	s.Step(`^the configuration is:$`, theConfigurationIs)

	s.Step(`^the API delivery http server is available at port (\d+) for at most (\d+) seconds and (\d+) request$`, theAPIDeliveryHttpServerIsAvailableAt)
	s.Step(`^the API delivery http server always responds with status code (\d+)$`, theAPIDeliveryHttpServerAlwaysRespondsWithStatusCode)
	s.Step(`^I call dokku "([^"]*)"$`, iCallDokku)
	s.Step(`^the API delivery http server received request (\d+) with the following JSON at "([^"]*)":$`, theAPIDeliveryHttpServerReceivedTheFollowingJSONAtEvent)
}
