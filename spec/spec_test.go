package spec

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	. "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type SpecSuite struct{}

var _ = Suite(&SpecSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *SpecSuite) TestParsing(c *C) {
	spec, err := Read("../testdata/test1.spec")

	c.Assert(err, NotNil)
	c.Assert(spec, IsNil)

	spec, err = Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	c.Assert(spec.GetLine(-1), DeepEquals, Line{-1, ""})
	c.Assert(spec.GetLine(99), DeepEquals, Line{-1, ""})
	c.Assert(spec.GetLine(34), DeepEquals, Line{34, "%{__make} %{?_smp_mflags}"})
}

func (s *SpecSuite) TestSections(c *C) {
	spec, err := Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	c.Assert(spec.HasSection(SECTION_BUILD), Equals, true)
	c.Assert(spec.HasSection(SECTION_CHECK), Equals, false)

	sections := spec.GetSections()
	c.Assert(sections, HasLen, 13)
	sections = spec.GetSections(SECTION_BUILD)
	c.Assert(sections, HasLen, 1)
	c.Assert(sections[0].Data, HasLen, 2)
	sections = spec.GetSections(SECTION_SETUP)
	c.Assert(sections[0].Name, Equals, "setup")
	c.Assert(sections[0].Args, DeepEquals, []string{"-qn", "%{name}-%{version}"})
}

func (s *SpecSuite) TestHeaders(c *C) {
	spec, err := Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	headers := spec.GetHeaders()
	c.Assert(headers, HasLen, 2)
	c.Assert(headers[0].Package, Equals, "")
	c.Assert(headers[0].Subpackage, Equals, false)
	c.Assert(headers[0].Data, HasLen, 11)
	c.Assert(headers[1].Package, Equals, "magic")
	c.Assert(headers[1].Subpackage, Equals, true)
	c.Assert(headers[1].Data, HasLen, 4)

	pkgName, subPkg := parsePackageName("%package magic")
	c.Assert(pkgName, Equals, "magic")
	c.Assert(subPkg, Equals, true)
	pkgName, subPkg = parsePackageName("%package -n magic")
	c.Assert(pkgName, Equals, "magic")
	c.Assert(subPkg, Equals, false)
}