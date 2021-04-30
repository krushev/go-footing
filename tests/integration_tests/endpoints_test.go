package integration_tests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/krushev/go-footing/controllers"
	"github.com/krushev/go-footing/db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"testing"
)

type EndpointsTestSuite struct {
	*suite.Suite
	app *fiber.App
	Conn *gorm.DB
}

func (s *EndpointsTestSuite) SetupSuite() {
	viper.Set("app.port", "3000")
	viper.Set("app.log", true)
	viper.Set("app.cors", true)
	viper.Set("db.drv", "sqlite")
	viper.Set("db.name", "file::memory:?cache=shared") // footing-test.db (path to a file / memory instead)
	viper.Set("db.log", "silent") // silent, error, warn, info

	viper.Set("jwt.timeout", 15)
	viper.Set("jwt.maxRefresh", 45)
	viper.Set("jwt.loginUrl", "http://localhost:" + viper.GetString("app.port") + "/api/login")

	s.Conn = db.NewConnection().GetDB()
	s.app = controllers.NewApp(true)
}

func (s *EndpointsTestSuite) SetupTest() {
	// runs once before every individual test runs. you can do things like clear common
	// state so that everything is "clean" for the test to run. be careful if you're
	// doing things to an external database, though. by default, tests run in parallel
	// so if all your tests share the same database and tables, they'll be overwriting
	// each other.
	//
	// you can fix that by either using unique DB/table names or running the tests without
	// parallelism (go test -p 1 ./...)
	zap.S().Debugf("SetupTest")
}

func (s *EndpointsTestSuite) TearDownTest() {
	// runs once after every individual test runs. you can clean up after your tests here.
	// a common thing here is to shut down an HTTP server that you left running
	zap.S().Debugf("TearDownTest")
}

func (s *EndpointsTestSuite) TearDownSuite() {
	// runs once after the entire test suite runs. a common thing to do here is shut down
	// a database you created before the tests ran
	zap.S().Debugf("TearDownSuite")
}

// BeforeTest enable hook for cleaning database
func (s *EndpointsTestSuite) BeforeTest(suiteName, testName string) {
	zap.S().Debugf("Before test %s from suite %s", suiteName, testName)
}

// AfterTest trigger the hook
func (s *EndpointsTestSuite) AfterTest(suiteName, testName string) {
	zap.S().Debugf("After test %s from suite %s", suiteName, testName)
}

func TestEndpointsTestSuite(t *testing.T) {
	suite.Run(t, &EndpointsTestSuite{
		Suite: new(suite.Suite),
		app:   nil,
		Conn:  nil,
	})
}
