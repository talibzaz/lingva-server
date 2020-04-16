package database

import (
	"fmt"
	"crypto/sha256"
	"lingva/auth"
	"strconv"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"regexp"
	"encoding/base64"
	"io/ioutil"
	"math/rand"
	"strings"
	"os"
	"time"
)

//Logs in user based on his role of teacher or parent.
func Login(username, password string) (Response, error){
	if username == "" || password == "" {
		return Response{
			Status: false,
			Message: "username / password cannot be empty",
		}, nil
	}
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	var res Response
	var teacher Teacher

	if err := DB.Select("id").Where("username = ? AND password = ? ", username, hashedPassword).First(&teacher).Error; err != nil {

		//If we get empty id we will check in students table to see whether user may be parent.
		if err.Error() == "record not found"{

			var student Student
			if err := DB.Select("id, student_name").Where("username = ? AND password = ? ", username, hashedPassword).First(&student).Error; err != nil {
				if err.Error() == "record not found" {
					return Response{
						Status: false,
						Message: "incorrect email / password",
					}, nil
				}
				return Response{}, err
			}
			res.ID = student.ID
			res.Status = true
			res.Message = "Login Success"
			token, err := auth.GenerateToken()
			if err != nil {
				return Response{}, err
			}
			res.Token = token
			res.Role = "parent"
			res.StudentName = student.StudentName

			return res, nil
		}
		return Response{}, err
	}

	res.ID = teacher.ID
	res.Status = true
	res.Message = "Login Success"
	token, err := auth.GenerateToken()
	if err != nil {
		return Response{}, err
	}
	res.Token = token
	res.Role = "teacher"

	return res, nil
}

//Login for Admin..
func AdminLogin(username, password string) (Response, error) {
	if username == "" || password == "" {
		return Response{
			Status: false,
			Message: "username / password cannot be empty",
		}, nil
	}

	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	var user User
	var res Response

	if err := DB.Select("id").Where("username = ? && password = ? ", username, hashedPassword).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return Response{
				Status: false,
				Message: "incorrect email / password",
			}, nil
		}
		return Response{}, err
	}

	res.ID = strconv.Itoa(user.ID)
	res.Status = true
	res.Message = "login success"
	token, err := auth.GenerateToken()
	if err != nil {
		return Response{}, err
	}
	res.Token = token
	res.Role = "admin"

	return res, nil
}

//Fetches the list of students for given teacher id.
func GetStudentList(teacherId string) ([]StudentList, error){

	var stdList []StudentList

	query := "select MAX(sr.date) as date, s.id as id, t.subject as subject, " +
		"s.student_name as name, s.image_url as image_url " +
		"from students s " +
		"left join student_reports sr on s.id = sr.student_id " +
		"inner join teachers t on s.teacher_id = t.id " +
		"and t.id = ? " +
		"group by id "

	if err := DB.Raw(query, teacherId).Scan(&stdList).Error; err != nil {
		return nil, err
	}

	return stdList, nil
}

//Saves new student report
func SaveStudentReport(report *ReportInput, testResults []TestResultsInput) (Data, error) {
	var stdReport StudentReport
	var result TestResults

	id, _ := uuid.NewV4()

	stdReport.ID = id.String()
	stdReport.StudentID = report.StudentID
	stdReport.TeacherID = report.TeacherID
	stdReport.Date = report.Date
	stdReport.HwCompletion = report.HwCompletion
	stdReport.VocabularyCompletion = report.VocabularyCompletion
	stdReport.ClassParticipation = report.ClassParticipation
	stdReport.DictionaryUse = report.DictionaryUse
	stdReport.TestCorrection = report.TestCorrection
	stdReport.Feedback = report.Feedback
	stdReport.ConversationClub = report.ConversationClub
	stdReport.Workshops = report.Workshops
	stdReport.MovieClub = report.MovieClub
	stdReport.ReadingClub = report.ReadingClub

	res :=  DB.Create(&stdReport)
	res.Last(&stdReport)

	var data Data
	data.StudentReport = stdReport

	for _, v := range testResults {

		result.ID = v.ID
		result.Date = v.Date
		result.StudentID = stdReport.StudentID
		result.TotalQuestions = v.TotalQuestions
		result.CorrectAnswers = v.CorrectAnswers
		result.Reference = v.Reference
		result.ReportID =stdReport.ID
		res := DB.Create(&result)
		res.Last(&result)
		data.TestResults = append(data.TestResults, result)
	}

	return data, nil
}

//Updates the previous student report..
func UpdateStudentReportByReportID(report *ReportInput, testResults []TestResultsInput) (StudentReport, error) {

	var stdReport StudentReport
	var result TestResults

	stdReport.Date = report.Date
	stdReport.HwCompletion = report.HwCompletion
	stdReport.VocabularyCompletion = report.VocabularyCompletion
	stdReport.ClassParticipation = report.ClassParticipation
	stdReport.DictionaryUse = report.DictionaryUse
	stdReport.TestCorrection = report.TestCorrection
	stdReport.Feedback = report.Feedback
	stdReport.ConversationClub = report.ConversationClub
	stdReport.Workshops = report.Workshops
	stdReport.MovieClub = report.MovieClub
	stdReport.ReadingClub = report.ReadingClub
	stdReport.MarkRead = false

	//if err := DB.Model(&StudentReport{}).Where("id = ?", report.ID).Update(stdReport).Error; err != nil {
	//	log.Info(err)
	//	return Response{}, err
	//}

	res := DB.Model(&StudentReport{}).Where("id = ?", report.ID).Update(stdReport)

	if err := res.Error; err != nil {
		log.Info(err)
		return StudentReport{}, err
	}

	res.Last(&stdReport)

	for _, v := range testResults {
		result.Date = v.Date
		result.TotalQuestions = v.TotalQuestions
		result.CorrectAnswers = v.CorrectAnswers
		result.Reference = v.Reference
		result.ReportID = report.ID
		result.ID = v.ID
		result.StudentID = report.StudentID

		if err := DB.Where("id = ?", v.ID).Assign(result).FirstOrCreate(&result).Error; err != nil {
			log.Info(err)
			return StudentReport{}, err
		}
	}

	return  stdReport, nil

}

//This generates all the Test Reports of a particular student based on his StudentID.
func GetStudentTestReportsByStudentID(studentId string) ([]StudentReport, error){
	var report []StudentReport

	if err := DB.Where("student_id = ?", studentId).Find(&report).Error; err != nil {
		return nil, err
	}


	return report, nil
}

func GetAdminNotificationByStudentID(studentId string) ([]Notification, error) {
	var notification []Notification

	if err := DB.Where("student_id = ?", studentId).Find(&notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

//This generates the report of a paricular student based on his/her reportID.
func GetTestReportByReportID(reportId string) (StudentReport, []TestResults, error) {
	var report StudentReport
	var testResult []TestResults

	if err := DB.Where("id = ?", reportId).First(&report).Error; err != nil {
		return StudentReport{}, nil, err
	}

	if err := DB.Where("report_id = ?", reportId).Find(&testResult).Error; err != nil {
		return StudentReport{}, nil, err
	}

	return report, testResult, nil
}

//Mark the report as read.
func MarkReportAsRead(reportId string) (StudentReport, error) {
	res := DB.Model(&StudentReport{}).Where("id = ?", reportId).Update("mark_read", true)

	if err := res.Error; err != nil  {
		return StudentReport{}, err
	}

	var stdReport StudentReport
	res.Last(&stdReport)

	return stdReport, nil
}

func MarkNotificationAsRead(notificationId string) (Notification, error) {
	res := DB.Model(&Notification{}).Where("id = ?", notificationId).Update("mark_read", true)

	if err := res.Error; err != nil {
		return Notification{}, err
	}

	var notification Notification
	res.Last(&notification)

	return notification, nil
}

func AllNotifications() ([]Notification, error) {
	var notification []Notification

	query := `select created_at, student_name, message, count(id) as total
              from notifications
              group by created_at`

	if err := DB.Raw(query).Scan(&notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

func NotificationsStudentList(createdAt string) ([]Notification, error) {
	var notification []Notification

	query := `select student_name from notifications 
			where created_at = ?`

	if err := DB.Raw(query, createdAt).Scan(&notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

//List of teachers for admin..
func TeachersList() ([]Teacher, error) {
	var teacher []Teacher

	query := "select t.id as id, t.name as name, t.subject as subject, COUNT(s.id) as count " +
		"from teachers t left join students s on t.id = s.teacher_id group by t.id"

	if err := DB.Raw(query).Scan(&teacher).Error; err != nil {
		return nil, err
	}

	return teacher, nil
}

//Gets teacher details by teacher id for admin
func GetTeacherByTeacherID(teacherId string) (Teacher, error){
	var teacher Teacher

	if err := DB.Where("id = ? ", teacherId).First(&teacher).Error; err != nil {
		return Teacher{}, err
	}

	return teacher, nil
}

func UpdateTeacherByTeacherID(input TeacherInput) (Response, error) {
	var teacher Teacher

	teacher.Name = input.Name
	teacher.Username = input.Username
	teacher.Password = input.Password
	teacher.Subject = input.Subject

	if err := DB.Model(&Teacher{}).Where("id = ?", input.ID).Update(teacher).Error; err != nil {
		return Response{}, err
	}

	return Response{
		Status: true,
		Message: "update successful.",
	}, nil
}

//Adds new teacher into the db...
func AddNewTeacher(input TeacherInput) (Response, error) {
	var teacher Teacher

	id, _ := uuid.NewV4()

	teacher.ID = id.String()
	teacher.Name = input.Name
	teacher.Subject = input.Subject
	teacher.Username = input.Username
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(input.Password)))
	teacher.Password = hashedPassword

	if err := DB.Create(&teacher).Error; err != nil {
		return Response{}, err
	}

	return Response{
		Status: true,
		Message: "successfully created new teacher.",
	}, nil
}

//Generates student list for admin
func GetStudentListForAdmin() ([]StudentList, error){
	var stdList []StudentList

	query := "select s.id as id, t.subject as subject, s.student_name as name, " +
	"t.name as teacher from students s " +
	"inner join teachers t on s.teacher_id = t.id " +
	"group by id "

	if err := DB.Raw(query).Scan(&stdList).Error; err != nil {
		return nil, err
	}

	return stdList, nil
}

//Fetches student details for admin
func GetStudentDetails(studentId string) (Student, error) {
	var student Student

	query := "select s.id as id, s.student_name as student_name, " +
		"t.name as teacher_name, s.parent_name as parent_name, " +
		"s.status as status, s.username as username, " +
		"s.password as password, t.id as teacher_id " +
		"from students s inner join teachers t on s.teacher_id = t.id "+
	"where s.id = ? "

	if err := DB.Raw(query, studentId).Scan(&student).Error; err != nil {
		return Student{}, err
	}

	return student, nil
}

//Adds new student in db..
func AddNewStudent(input StudentInput, imageData string) (Response, error) {
	var student Student
	var imageName string
	var err error

	id, _ := uuid.NewV4()

	if imageData != "" {
		imageName, err = saveImage(imageData)
		if err != nil {
			return Response{}, err
		}
	}

	student.ID = id.String()
	student.StudentName = input.Name
	student.TeacherID = input.TeacherID
	student.ParentName = input.ParentName
	student.Status = input.Status
	student.Username = input.Username
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(input.Password)))
	student.Password = hashedPassword
	student.ImageUrl = os.Getenv("IMAGE_URL") +imageName

	if err := DB.Create(&student).Error;err != nil {
		return Response{}, err
	}

	return Response{
		Status: true,
		Message: "successfully added new student.",
	}, nil
}

func saveImage(imageData string) (string, error){
	re := regexp.MustCompile("[:;,]")
	result := re.Split(imageData, -1)

	extension := strings.Split(result[1], "/")

	decodedImage, _ := base64.StdEncoding.DecodeString(result[3])

	imageName := randomString(20)+"."+extension[1]

	err := ioutil.WriteFile("files/"+imageName, decodedImage, 0666)
	if err != nil {
		log.Error("err", err)
		return "", err
	}
	return imageName, nil
}

func randomString(len int) string {
	bytes := make([]byte, len)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25))
	}
	return string(bytes)
}

func UpdateStudentByStudentID(input StudentInput) (Response, error) {
	var student Student

	student.StudentName = input.Name
	student.TeacherID = input.TeacherID
	student.ParentName = input.ParentName
	student.Status = input.Status
	student.Username = input.Username
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(input.Password)))
	student.Password = hashedPassword

	if err := DB.Model(&student).Where("id = ?", input.ID).Update(student).Error; err != nil {
		return Response{}, err
	}

	return Response{
		Status: true,
		Message: "student updated successfully",
	}, nil
}

func AddStudentNotification(notificationInput []NotificationInput, message string) (Response, error){
	var notification Notification

	date := time.Now().String()

	for _, v := range notificationInput {
		id, _ := uuid.NewV4()

		notification.ID = id.String()
		notification.StudentName = v.StudentName
		notification.StudentID = v.StudentID
		notification.Message = message
		notification.CreatedAt = date

		if err := DB.Create(&notification).Error; err != nil {
			return Response{}, err
		}
	}
	return Response{Status: true, Message: "Notification Sent."}, nil
}