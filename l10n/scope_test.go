package l10n

import (
	"testing"

	"github.com/jinzhu/gorm"
	. "gopkg.in/check.v1"
)

//Hook up gocheck into the go test runner
func Test(t *testing.T) {
	TestingT(t)
}

type ScopeSuite struct {
	Scope   *gorm.Scope
	ScopeCN *gorm.Scope
}

var _ = Suite(&ScopeSuite{})

func (s *ScopeSuite) SetUpSuite(c *C) {
	product := Product{}
	scope := dbGlobal.NewScope(&product)
	scopeCN := dbCN.NewScope(&product)
	s.Scope = scope
	s.ScopeCN = scopeCN
}

func (s *ScopeSuite) TestIsLocalizable(c *C) {

	result := IsLocalizable(s.Scope)
	c.Assert(result, Equals, true)
}

func (s *ScopeSuite) TestGetLocale(c *C) {
	globalLocale, isLocale := getLocale(s.Scope)
	c.Assert(isLocale, Equals, false)
	c.Assert(Global, Equals, globalLocale)
	zh, isLocale := getLocale(s.ScopeCN)

	c.Assert(isLocale, Equals, true)
	c.Assert(zh, Equals, "zh")
}

func (s *ScopeSuite) TestIsLocaleCreateable(c *C) {
	ok := isLocaleCreatable(s.Scope)
	c.Assert(ok, Equals, false)
}

func (s *ScopeSuite) BenchmarkIsLocalizable(c *C) {
	for i := 0; i < c.N; i++ {
		// Logic to benchmark
		IsLocalizable(s.Scope)
	}
}
func (s *ScopeSuite) BenchmarkSetLocale(c *C) {
	for i := 0; i < c.N; i++ {
		// Logic to benchmark
		setLocale(s.Scope, "zh")
	}
}
