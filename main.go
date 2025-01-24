package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

const (
	seleniumPath = "chromedriver"
	port         = 8080
	testURL      = "https://chulo-solutions.github.io/qa-internship/"
)

var ValidTestData = map[string]string{
	"username":   "biraj123",
	"password":   "Password1*",
	"creditCard": "4111666622227777",
	"telephone":  "(123) 456-7890",
}

// TestCase represents a negative test case scenario
type TestCase struct {
	Field         string
	Input         string
	ExpectedError string
}

// InvalidTestCases
var InvalidTestCases = []TestCase{
	// Invalid username but other fields valid
	{
		"username",
		"user",
		"Username must be alphanumeric and between 5 to 15 characters.",
	},
	// Invalid password but other fields valid
	{
		"password",
		"Password123",
		"Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.",
	},
	// Invalid credit card but other fields valid
	{
		"creditCard",
		"1234567890123456",
		"Enter a valid credit card number.",
	},
	// Invalid telephone but other fields valid
	{
		"telephone",
		"1234567890",
		"Telephone number must follow the format (XXX) XXX-XXXX.",
	},
}

func main() {
	// Starting selenium webdriver
	service, err := selenium.NewChromeDriverService(seleniumPath, port)
	if err != nil {
		log.Fatalf("Error starting ChromeDriver: %v", err)
	}
	defer service.Stop()

	// Connect to WebDriver instance
	caps := selenium.Capabilities{"browserName": "chrome"}
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Error connecting to WebDriver: %v", err)
	}
	defer driver.Quit()

	// Get the registration form of above URL
	if err := driver.Get(testURL); err != nil {
		log.Fatalf("Error opening page: %v", err)
	}

	time.Sleep(2 * time.Second) // Initial wait for page load

	// Fill form with valid data
	if err := FillForm(driver, ValidTestData); err != nil {
		log.Fatalf("Error filling form: %v", err)
	}

	// Verify success message (if applicable)
	if success := VerifySuccessMessage(driver); !success {
		log.Println("Test Failed: Success message verification failed.")
	}

	// Run negative test cases sequentially
	for _, testCase := range InvalidTestCases {
		RunNegativeTest(driver, testCase)

		time.Sleep(2 * time.Second)
	}
}

func FillForm(driver selenium.WebDriver, formData map[string]string) error {
	for name, value := range formData {
		field, err := driver.FindElement(selenium.ByName, name)
		if err != nil {
			return fmt.Errorf("could not find field %s: %v", name, err)
		}
		field.Clear()
		field.SendKeys(value)
	}

	submitBtn, err := driver.FindElement(selenium.ByXPATH, "//button[text()='Submit']")
	if err != nil {
		return fmt.Errorf("could not find submit button: %v", err)
	}
	submitBtn.Click()
	// 2 secs delay for submission process
	time.Sleep(2 * time.Second)
	return nil
}

// VerifySuccessMessage checks if the success alert is displayed
func VerifySuccessMessage(driver selenium.WebDriver) bool {
	alertText, err := driver.AlertText()
	if err != nil {
		log.Println("Test Failed: Success alert not found.")
		return false
	}

	driver.AcceptAlert()

	if alertText == "Form submitted successfully!" {
		log.Println("Test Passed: Form submitted successfully!")
		return true
	} else {
		log.Println("Test Failed: Unexpected success alert message.")
		return false
	}
}

// RunNegativeTest
func RunNegativeTest(driver selenium.WebDriver, testCase TestCase) {

	fields := map[string]string{
		"username":   "biraj123",
		"password":   "Password1*",
		"creditCard": "4111666622227777",
		"telephone":  "(123) 456-7890",
	}

	for id, value := range fields {
		element, err := driver.FindElement(selenium.ByID, id)
		if err != nil {
			log.Printf("Error finding field %s: %v", id, err)
			return
		}
		element.Clear()
		if id == testCase.Field {
			element.SendKeys(testCase.Input) // Insert invalid value
		} else {
			element.SendKeys(value) // Insert valid values
		}
	}

	// Submit the form
	submitBtn, err := driver.FindElement(selenium.ByXPATH, "//button[text()='Submit']")
	if err != nil {
		log.Println("Error finding submit button:", err)
		return
	}

	err = submitBtn.Click()
	if err != nil {
		log.Println("Error clicking submit button:", err)
		return
	}

	// Handle alert popup
	time.Sleep(1 * time.Second) // Allow alert to appear
	alertText, err := driver.AlertText()
	if err != nil {
		log.Printf("Test Failed: No alert for field '%s'. Error: %v", testCase.Field, err)
		return
	}

	driver.AcceptAlert()

	// verification for error message matches the alert message
	if alertText == testCase.ExpectedError {
		log.Printf("Test Passed: Correct alert for field '%s'.", testCase.Field)
	} else {
		log.Printf("Test Failed: Incorrect alert for field '%s'. Expected: %s, Got: %s",
			testCase.Field, testCase.ExpectedError, alertText)
	}
}
