package controllers

import (
  "fmt"
  "strconv"

  "net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
  "github.com/wisedog/ladybug/models"
  "github.com/wisedog/ladybug/interfacer"
	"github.com/wisedog/ladybug/errors"
  log "gopkg.in/inconshreveable/log15.v2"
)

// makeMessage is private utility function for comparing and make
// messages for status changing
//func  makeMessage(historyUnit *[]models.HistoryTestCaseUnit){
func makeHistoryMessage(historyUnit *[]models.HistoryTestCaseUnit){

	if historyUnit == nil{
		log.Error("Nil historyunit!")
		return
	}
	
	var msg string
	
	// TODO L10N with using GetI18nMessage format string
	for i := 0; i < len(*historyUnit); i++{
		if (*historyUnit)[i].ChangeType == models.HISTORY_CHANGE_TYPE_CHANGED{
			if((*historyUnit)[i].From == 0 && (*historyUnit)[i].To == 0 ){
				msg = fmt.Sprintf(`"%s" is changed from "%s" to "%s".`, 
				(*historyUnit)[i].What, (*historyUnit)[i].FromStr, (*historyUnit)[i].ToStr)
			}else{
				msg = fmt.Sprintf(`"%s" is changed from %d to %d.`, 
				(*historyUnit)[i].What, (*historyUnit)[i].From, (*historyUnit)[i].To)
			}
		}else if(*historyUnit)[i].ChangeType == models.HISTORY_CHANGE_TYPE_SET{
			msg = fmt.Sprintf(`"%s" is set to "%s".`, 
				(*historyUnit)[i].What, (*historyUnit)[i].Set)
		}else if(*historyUnit)[i].ChangeType == models.HISTORY_CHANGE_TYPE_DIFF{
			msg = fmt.Sprintf(`"%s" is changed(diff).`, 
				(*historyUnit)[i].What)
		}else if(*historyUnit)[i].ChangeType == models.HISTORY_CHANGE_TYPE_NOTE{
			msg = fmt.Sprintf("%s added a note.", (*historyUnit)[i].What)
		}else{
			msg = ""
		}
		(*historyUnit)[i].Msg = msg
	}
}


// CaseView renders a page to show testcase's information
func CaseView(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }

  vars := mux.Vars(r)
  project := vars["projectName"]
  id := vars["id"]

	var tc models.TestCase
  if err := c.Db.Where("id = ?", id).Preload("Category").First(&tc); err.Error != nil{
    log.Error("testcases", "Error on Select-ID query in TestCases.Index", err.Error)
    return errors.HttpError{http.StatusInternalServerError, "DB error"}
  }
	
	tc.PriorityStr = getPriorityI18n(tc.Priority)
	
	var histories []models.History
	c.Db.Where("category = ?", models.HISTORY_TYPE_TC).Where("target_id = ?", tc.ID).Preload("User").Find(&histories)
	
	
	// making status changement messages
	for i := 0; i < len(histories); i++{
		var res []models.HistoryTestCaseUnit
		json.Unmarshal([]byte(histories[i].ChangesJson), &res)
		
		// make message
		makeHistoryMessage(&res)
		
		histories[i].Changes = res
		log.Debug("testcases", "res", histories[i].Changes)
	}
	
  items := struct {
    User *models.User
    ProjectName string
    TestCase  models.TestCase
    Histories  []models.History
    Active_idx  int
  }{
    User: user,
    ProjectName : project,
    TestCase : tc,
    Histories : histories,
    Active_idx : 2,
  }
  return Render(w, items, "views/base.tmpl", "views/testcases/caseindex.tmpl")
}


// renderAddEdit renders ADD and EDIT pages
func renderAddEdit(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request, isEdit bool) error{
  var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }

  vars := mux.Vars(r)
  projectName := vars["projectName"]
  id := vars["id"]

  // valid when rendering Add 
  section_id :=  r.URL.Query().Get("sectionid")
  if section_id == "" && isEdit == false {
    log.Error("testcase", "type" , "general", "msg", "invalid condition. Add rendering requires section id")
    return errors.HttpError{http.StatusBadRequest, "invalid condition. Add rendering requires section id"}
  }

  var category []models.Category
  c.Db.Find(&category)
  
  // TODO when validation failed, do something like below code. 
  // belows are revel's code
  /*if c.Validation.HasErrors(){
    c.Validation.Keep()
    c.FlashParams()
  
    return c.Render(project, id, category)
  }*/
  
  var testcase models.TestCase
  if isEdit{
    if err := c.Db.Where("id = ?", id).First(&testcase); err.Error != nil{
      log.Error("testcase", "type" , "database", "msg", err.Error)
      return errors.HttpError{http.StatusInternalServerError, "An Error while SELECT operation for TestCase.Edit"}
    }
  }else{
    var prj models.Project
    if err := c.Db.Where("name = ?", projectName).First(&prj); err.Error != nil{
      log.Error("testcase", "type" , "database", "msg", err.Error)
      return errors.HttpError{http.StatusInternalServerError, "An Error while SELECT operation for TestCase.Edit"}
    }

    testcase.SectionID, _ = strconv.Atoi(section_id)
    testcase.Prefix = prj.Prefix
  }
  
  
  //TODO change section here.

  items := struct {
    Id string
    User *models.User
    ProjectName string
    SectionID string
    TestCase models.TestCase
    Category []models.Category
    IsEdit bool
    Active_idx  int
  }{
    Id : id,
    User: user,
    ProjectName : projectName,
    SectionID : section_id,
    TestCase : testcase,
    Category : category,
    IsEdit : isEdit,
    Active_idx : 2,
  }

  return Render(w, items, "views/base.tmpl", "views/testcases/caseaddedit.tmpl")
}

// CaseEdit just renders a EDIT page. Add and Edit pages are integrated. see renderAddEdit
func CaseEdit(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
	return renderAddEdit(c, w, r, true)
}

// CaseAdd just renders a Add page. Add and Edit pages are integrated. see renderAddEdit
func CaseAdd(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{

  return renderAddEdit(c, w, r, false)
}

// handleSaveUpdate handles POST request from save and update
func handleSaveUpdate(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request, isUpdate bool) error{
  var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }
  
  var testcase models.TestCase
  vars := mux.Vars(r)
  projectName := vars["projectName"]
  id := vars["id"]

  if err := r.ParseForm(); err != nil {
    log.Error("Testcase", "Error ", err )
    return errors.HttpError{http.StatusInternalServerError, "ParseForm failed"}
  }

  //fmt.Printf("After ParseForm : %+v\n", r)

  decoder := schema.NewDecoder()
  
  if err := decoder.Decode(&testcase, r.PostForm); err != nil {
    log.Warn("Testcase", "Error", err, "msg", "Decode failed but go ahead")
  }

  // TODO : Validate input testcase
  //testcase.Validate()

  redirectionTarget := fmt.Sprintf("/project/%s/design", projectName)

  if isUpdate == false{

    var largest_seq_tc models.TestCase

    if err := c.Db.Where("prefix=?", testcase.Prefix).Order("seq desc").First(&largest_seq_tc); err.Error != nil {
      log.Error("An error while SELECT operation to find largest seq")
    } else {
      testcase.Seq = largest_seq_tc.Seq + 1
      log.Debug("testcase", "Largest number is : ", largest_seq_tc.Seq+1)
    }

    c.Db.NewRecord(testcase)

    if err := c.Db.Create(&testcase); err.Error != nil {
      log.Error("testcase", "type" , "database", "msg", err.Error)
      return errors.HttpError{http.StatusInternalServerError, "Insert operation failed in TestCase.Save"}
    }
    
    display_id := testcase.Prefix + "-" + strconv.Itoa(testcase.ID)
    testcase.DisplayID = display_id
    
    
    if err := c.Db.Save(&testcase); err.Error != nil{
      log.Error("testcase", "type" , "database", "msg", err.Error)
      return errors.HttpError{http.StatusInternalServerError, "Update operation failed in TestCase.Save"}
    }

    http.Redirect(w, r, redirectionTarget, http.StatusFound)
  }
  
  var existCase models.TestCase

  if err := c.Db.Where("id = ?", id).First(&existCase); err.Error != nil {
    log.Error("testcase", "type", "database", "msg", err.Error)
    return errors.HttpError{http.StatusInternalServerError, "An error while select exist testcase operation"}
  }
  note := r.FormValue("Note")
  findDiff(c, &existCase, &testcase, note, user)
  
  existCase.Title = testcase.Title
  existCase.Description = testcase.Description
  existCase.Seq = testcase.Seq
  existCase.Status = testcase.Status
  existCase.Prefix = testcase.Prefix
  existCase.Precondition = testcase.Precondition
  existCase.Steps = testcase.Steps
  existCase.Expected = testcase.Expected
  existCase.Priority = testcase.Priority

  if err := c.Db.Save(&existCase); err.Error != nil {
    log.Error("testcase", "type", "database", "msg", err.Error)
    return errors.HttpError{http.StatusInternalServerError, "An error while SAVE operation on TestCases.Update"}
  }
  
  //c.Flash.Success("Update Success!")
  
  http.Redirect(w, r, redirectionTarget, http.StatusFound)
  return nil 
}

// CaseSave handles POST request from CaseAdd. Save and Update handlers are integrated. see handleSaveUpdate
func CaseSave(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
  return handleSaveUpdate(c, w, r, false)
}

// CaseUpdate POST handler for Testcase Edit
func CaseUpdate(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
	return handleSaveUpdate(c, w, r, true)
}

// CaseDelete is a POST handler for DELETE request
func CaseDelete(c *interfacer.AppContext, w http.ResponseWriter, r *http.Request) error{
    var user *models.User
  if user = connected(c, r); user == nil{
    log.Debug("Not found login information")
    http.Redirect(w, r, "/", http.StatusFound)
  }

  if err := r.ParseForm(); err != nil {
    log.Error("Testcase", "Error ", err )
    return errors.HttpError{http.StatusInternalServerError, "ParseForm failed"}
  }

  id := r.FormValue("id")

  // Delete the testcase  permanently for sequence
  if err := c.Db.Unscoped().Where("id = ?", id).Delete(&models.TestCase{}); err.Error != nil{
    log.Error("testcase", "An error while delete testcase ", err.Error)
    return errors.HttpError{http.StatusInternalServerError, "testcase delete failed"} 
  }

  // the client do redirect or refresh. 
  return nil
}

// findDiff compares between two models.TestCase and create 
// HistoryTestCaseUnit to database. Used in Update
func  findDiff(c *interfacer.AppContext, existCase, newCase *models.TestCase, note string, user *models.User){
	
	var changes []models.HistoryTestCaseUnit
	his := models.History{Category : models.HISTORY_TYPE_TC,
			TargetID : existCase.ID, UserID : user.ID,
	}
	
	// check note
	if note != ""{
		// this means user add a note on this testcase
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_NOTE, What : user.Name,
		}
		
		his.Note = note
		changes = append(changes, unit)
	}
	
	// check title
	if existCase.Title != newCase.Title {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_CHANGED, What : "Title",
			FromStr : existCase.Title, ToStr : newCase.Title,
		}
		changes = append(changes, unit)
	}
	
	// check priority
	if existCase.Priority != newCase.Priority {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_CHANGED, 
			What : GetI18nMessage("testcase.priority"),
			FromStr : getPriorityI18n(existCase.Priority),
			ToStr : getPriorityI18n(newCase.Priority),
		}
		changes = append(changes, unit)
	}
	
	// check execution type
	if existCase.ExecutionType != newCase.ExecutionType {
		arr := [2]string{"Manual", "Automated"}
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_CHANGED, What : "Execution type",
			FromStr : arr[existCase.ExecutionType], ToStr : arr[newCase.ExecutionType],
		}
		changes = append(changes, unit)
	}
	
	// TODO check Status
	
	// check Description
	if existCase.Description != newCase.Description {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, 
			What : GetI18nMessage("description"),
			DiffID : 2,
		}
		changes = append(changes, unit)
	}
	
	// check Precondition
	if existCase.Precondition != newCase.Precondition {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, 
			What : GetI18nMessage("priority.precondition"),
			DiffID : 2,//TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}
	
	// check Estimated
	if existCase.Estimated != newCase.Estimated {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, What : "Estimated Time",
			DiffID : 2,	//TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}
	
	// check Steps
	if existCase.Steps != newCase.Steps {
    // FIXME A bug : those strings are same, but,,,, 
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, What : "Steps",
			DiffID : 2,	//TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}
	
	// check Expected
	if existCase.Expected != newCase.Expected {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_DIFF, 
			What : GetI18nMessage("priority.expected"),
			DiffID : 2,	//TODO should be implemnted DIFF
		}
		changes = append(changes, unit)
	}
	
	var existCategory models.Category
	var testcaseCategory models.Category
	
	c.Db.Where("id = ?", existCase.CategoryID).Find(&existCategory)
	c.Db.Where("id = ?", newCase.CategoryID).Find(&testcaseCategory)
	
	// check CategoryID
	if existCase.CategoryID != newCase.CategoryID {
		unit := models.HistoryTestCaseUnit{
			ChangeType : models.HISTORY_CHANGE_TYPE_CHANGED, 
			What : GetI18nMessage("priority.category"),
			FromStr : "",
			ToStr : "",
		}
		changes = append(changes, unit)
	}
	
	result, _ := json.Marshal(changes)
	his.ChangesJson = string(result)
	
	c.Db.NewRecord(his)
	c.Db.Create(&his)
}

