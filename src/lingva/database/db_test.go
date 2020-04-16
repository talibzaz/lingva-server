package database

import (
	"testing"
	"encoding/json"
	"os"
	"fmt"
)

func TestLogin(t *testing.T) {
	var testLogin = []struct{
		username string
		password string
		expected bool
	}{
		{"miky", "qwertyu", true},
		{"john", "qwerty", false},
		{"ishfaq@gmail.com", "qwertyu", true},
		{"ishfaq@gmail.co", "qwertyu", false},
		{"ishfaq@gmail.com", "Qwertyu", false},
		{"aQuib2", "qwertyU", false},
	}

	for _, ts := range testLogin {
		res, _ := Login(ts.username, ts.password)

		if res.Status != ts.expected {
			t.Errorf("Expected result for %s and %s is %t, but got %t", ts.username, ts.password, ts.expected, res)
		}

	}

	//res, err := Login("miky", "qwertyu")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//e := json.NewEncoder(os.Stdout)
	//e.SetIndent(" ", " ")
	//e.Encode(&res)
}

func TestAdminLogin(t *testing.T) {

	res, err := AdminLogin("admi", "admi")
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestGetStudentList(t *testing.T) {
	res, _ := GetStudentList("2435867f-3e65-49bf-a59f-59625d6471ae")

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestGetStudentTestReportsByStudentID(t *testing.T) {
	res, err := GetStudentTestReportsByStudentID("308e5294-3f67-4ffc-ac58-eb1e1cb6fe17")
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestAllNotifications(t *testing.T) {
	res, err := AllNotifications()
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestNotificationsStudentList(t *testing.T) {
	res, err := NotificationsStudentList("2018-12-07 10:47:22.680329808 +0530 IST m=+105.144449568")
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestGetAdminNotificationByStudentID(t *testing.T) {
	res, err := GetAdminNotificationByStudentID("2f29e0cf-f6cd-4651-b84a-89a50beb2860")
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestGetTestReportByReportID(t *testing.T) {
	_, res, err := GetTestReportByReportID("")
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestMarkReportAsRead(t *testing.T) {
	res, err := MarkReportAsRead("")
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestTeachersList(t *testing.T) {
	res, _ := TeachersList()

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestGetTeacherByTeacherID(t *testing.T) {
	fmt.Println(randomString(20))
	//res, _ := GetTeacherByTeacherID("c01b893c-a120-46e6-bb9c-fecc1d11cadb")
	//
	//e := json.NewEncoder(os.Stdout)
	//e.SetIndent(" ", " ")
	//e.Encode(&res)
}

func TestUpdateTeacherByTeacherID(t *testing.T) {
	var tr TeacherInput
	tr.ID = "c01b893c-a120-46e6-bb9c-fecc1d11cadb"
	tr.Username = "asif"
	tr.Password = "asif"
	res, err := UpdateTeacherByTeacherID(tr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res.Message)
}

func TestGetStudentListForAdmin(t *testing.T) {
	res, err := GetStudentListForAdmin()
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestGetStudentDetails(t *testing.T) {
	res, err := GetStudentDetails("308e5294-3f67-4ffc-ac58-eb1e1cb6fe17")
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}

func TestUpdateStudentReportByReportID(t *testing.T) {

	var r *ReportInput
	tr := make([]TestResultsInput, 2)

	r.ID = "0485ab97-ff37-42a5-b225-725f74894003"
	r.StudentID = "f410f2ba-ce55-4eb3-9cd5-685da0cfd75a"
	r.Date = "2018-11-11"
	r.HwCompletion = "partially"
	r.Feedback = "lacks interest in studies.."

	tr[0].ID = "07cd4c26-7407-4173-bd5b-af200865a795"
	tr[1].ID = "2c85799e-ce02-47e9-ba50-50516d09ccf0"
	tr[0].Date = "2018-10-30"
	tr[1].Date = "2018-10-30"
	tr[0].Reference = "test1"
	tr[1].Reference = "test2"

	res, err := UpdateStudentReportByReportID(r, tr)
	if err != nil {
		t.Fatal(err)
	}

	e := json.NewEncoder(os.Stdout)
	e.SetIndent(" ", " ")
	e.Encode(&res)
}