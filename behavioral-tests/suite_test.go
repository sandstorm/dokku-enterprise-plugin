package main

import "github.com/DATA-DOG/godog"

func FeatureContext(s *godog.Suite) {
	s.Step(`^I remove the file "([^"]*)"$`, iRemoveTheFile)
	s.Step(`^I have an empty node\.js application$`, iHaveAnEmptyNodejsApplication)
	s.Step(`^I create the file "([^"]*)" with the following contents:$`, iCreateTheFileWithTheFollowingContents)
	s.Step(`^I deploy the application as "([^"]*)"$`, iDeployTheApplicationAs)
	s.Step(`^I call the URL "([^"]*)" of the "([^"]*)" application$`, iCallTheURLOfTheApplication)
	s.Step(`^the response should contain "([^"]*)"$`, theResponseShouldContain)
	s.Step(`^the response should not contain "([^"]*)"$`, theResponseShouldNotContain)
	s.Step(`^the event log is empty$`, theEventLogIsEmpty)
	s.Step(`^I expect (\d+) event log entry$`, iExpectEventLogEntry)
	s.Step(`^the API delivery http server is disabled$`, theAPIDeliveryHttpServerIsDisabled)
	s.Step(`^the configuration is:$`, theConfigurationIs)
}
