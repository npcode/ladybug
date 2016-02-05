package models

import (
	"time"
	
	"github.com/revel/revel"
)

const (
	PRIORITY_HIGHEST = 1 + iota
	PRIORITY_HIGH
	PRIORITY_MEDIUM
	PRIORITY_LOW
	PRIORITY_LOWEST
	)

const (
	TC_STATUS_ACTIVATE = 1 + iota
	TC_STATUS_INACTIVE
	TC_STATUS_DRAFT
	)



// TestCase model represents a test case
type TestCase struct {
	ID            int
	
	// Prefix is a unique characters that each project has
	// The prefix is used when Human-friendly testcase name is created
	// for example, if the prefix is "TC" then a testcase name is "TC-1"
	// It is {{Prefix}}-{{ID}}, but will be {{Prefix}}-{{Seq}} soon
	Prefix        string `sql:"not null"`
	
	// DisplayID represents test case name for displaying. 
	// The format is looks like "TC-200"
	DisplayID		string	`sql:"not null"`
	
	// Seq is a unique sequence number in the project. 
	// It is not disabled now, but it will be used to make user-friendly test case name
	Seq           int    `sql:"not null"`
	
	// Title of the test case
	Title         string `sql:"size:400"`
	Section       Section
	SectionID     int    `sql:"index"`
	ExecutionType int    //Manual, Automated
	ExTypeStr			string
	Status        int    // Activation, Deactivation, closed
	Description   string `sql:"size:2000"` // description of the issue
	Precondition  string `sql:"size:1000"`
	Priority      int    // 1 to 5. 1 is highest priority
	PriorityStr		string	`sql:"-"`
	Estimated     int    // unit : min(s)
	RelatedReq    int    // TODO
	Version       int
	Steps         string `sql:"size:1000"`
	Expected      string `sql:"size:1000"`
	Project		  Project	
	ProjectID	  int
	Category		Category
	CategoryID		int
	
	// need many-to-many releationship between specification and testcase. 
	// Specifications       []Specification `gorm:"many2many:spec_cases;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Validation will check validation of the input form
func (testcase *TestCase) Validate(v *revel.Validation) {
	v.Required(testcase.Title)
	v.Check(testcase.Title,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{400},
	)
}
