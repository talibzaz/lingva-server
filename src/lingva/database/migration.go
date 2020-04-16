package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID 			int		`gorm:"primary_key;not null; auto_increment"`
	Username	string	`gorm:"unique;not null"`
	Password 	string
}

type Teacher struct {
	ID       		string	`gorm:"primary_key;not null"`
	Name     		string
	Subject  		string
	Username 		string
	Password 		string
	Count           int      `gorm:"-"`
}

type Student struct {
	ID          string 		`gorm:"primary_key;not null"`
	StudentName string
	Status      int
	TeacherID   string
	ParentName  string
	Username    string
	Password    string
	TeacherName	string 		`gorm:"-"`
	ImageUrl    string
}

type StudentReport struct {
	ID                   string	 	`gorm:"primary_key;not null"`
	StudentID            string 	`gorm:"not null"`
	Date                 string		`gorm:"type:varchar(10)"`
	HwCompletion         string
	VocabularyCompletion string
	ClassParticipation   string
	DictionaryUse        string
	ConversationClub     bool
	Workshops            bool
	MovieClub            bool
	ReadingClub          bool
	MarkRead             bool
	TestCorrection       string
	Feedback             string
	TeacherID            string
	Type                 string		`gorm:"default:'TEACHER'"`
}

type TestResults struct {
	ID             string 	`gorm:"primary_key;not null"`
	StudentID      string 	`gorm:"not null"`
	Date           string	`gorm:"type:varchar(10)"`
	Reference      string
	TotalQuestions int
	CorrectAnswers float64
	ReportID       string 		`gorm:"not null"`
}

type Notification struct {
	ID 				string		`gorm:"primary_key;not null"`
	StudentName 	string
	StudentID		string
	Message 		string		`gorm:"type:text"`
	CreatedAt 		string		`gorm:"unique;not null"`
	Type            string		`gorm:"default:'ADMIN'"`
	MarkRead        bool		`gorm:"default: false"`
	Total           int			`gorm:"-"`
}

type NotificationInput struct {
	StudentName 	string
	StudentID		string
}

type StudentInput struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     int    `json:"status"`
	TeacherID  string `json:"teacherId"`
	ParentName string `json:"parentName"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

type StudentList struct {
	Name 	string		`json:"name"`
	ID   	string 		`json:"id"`
	Date 	string 		`json:"date"`
	Subject	string		`json:"subject"`
	Teacher string		`json:"teacher"`
	ImageUrl string
}

type Response struct {
	ID      		string	`json:"id"`
	Status  		bool   	`json:"status"`
	Message 		string 	`json:"message"`
	Token 			string 	`json:"token"`
	Role  			string
	StudentName		string
}

type Data struct{
	StudentReport	StudentReport
	TestResults		[]TestResults
}

type ReportInput struct {
	ID                   string     `json:"id"`
	StudentID            string 	`json:"studentId"`
	TeacherID			 string		`json:"teacherId"`
	Date                 string 	`json:"date"`
	HwCompletion         string 	`json:"hwCompletion"`
	VocabularyCompletion string 	`json:"vocabularyCompletion"`
	ClassParticipation   string 	`json:"classParticipation"`
	DictionaryUse        string 	`json:"dictionaryUse"`
	TestCorrection       string 	`json:"testCorrection"`
	Feedback             string 	`json:"feedback"`
	ConversationClub     bool   	`json:"conversationClub"`
	Workshops            bool   	`json:"workshop"`
	MovieClub            bool   	`json:"movieClub"`
	ReadingClub          bool   	`json:"readingClub"`
}

type TestResultsInput struct {
	ID 			   string  `json:"id"`
	Date		   string  `json:"date"`
	Reference      string  `json:"reference"`
	TotalQuestions int     `json:"totalQuestions"`
	CorrectAnswers float64 `json:"correctAnswers"`
	ReportID	   *string
}

type TeacherInput struct {
	ID          string
	Name 		string
	Subject 	string
	Username 	string
	Password 	string
}

var DB *gorm.DB

func init() {
	db, err := gorm.Open("mysql", "root:root@tcp(138.68.70.81:3306)/lingva")
	if err != nil {
		panic(err)
	}

	DB = db

	migrationStatements()
}

func migrationStatements() {
	DB.AutoMigrate(&Teacher{}, &Student{}, &StudentReport{}, &TestResults{}, &User{}, &Notification{})

	DB.Model(&Student{}).AddForeignKey("teacher_id", "teachers(id)", "RESTRICT", "RESTRICT")
	DB.Model(&StudentReport{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	DB.Model(&StudentReport{}).AddForeignKey("teacher_id", "teachers(id)", "RESTRICT", "RESTRICT")
	DB.Model(&TestResults{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")
	DB.Model(&TestResults{}).AddForeignKey("report_id", "student_reports(id)", "RESTRICT", "RESTRICT")
	DB.Model(&Notification{}).AddForeignKey("student_id", "students(id)", "RESTRICT", "RESTRICT")

}
