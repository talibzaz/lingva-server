package server

import (
	"context"
	"lingva/gql"
	db "lingva/database"
	"errors"
)

func (s *server) Query() gql.QueryResolver {
	return queryResolver{s}
}

func (s *server) Mutation() gql.MutationResolver {
	return mutationResolver{s}
}

type queryResolver struct {*server}

type mutationResolver struct {*server}

func (s *server) GetStudentList(ctx context.Context, teacherId string) ([]db.StudentList, error) {
	res, err := db.GetStudentList(teacherId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) SaveReport(ctx context.Context, report *db.ReportInput, testResults []db.TestResultsInput)(db.Data, error){

	res, err := db.SaveStudentReport(report, testResults)

	if err != nil {
		return db.Data{}, err
	}

	return res, nil
}

func (s *server) UpdateReportByReportID(ctx context.Context, report *db.ReportInput, testResults []db.TestResultsInput)(db.StudentReport, error){

	res, err := db.UpdateStudentReportByReportID(report, testResults)

	if err != nil {
		return db.StudentReport{}, err
	}

	return res, nil
}

func (s *server) TestReportsByStudentID(ctx context.Context, studentId string) ([]db.StudentReport, error) {
	res, err := db.GetStudentTestReportsByStudentID(studentId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) AdminNotificationByStudentID(ctx context.Context, studentId string) ([]db.Notification, error) {
	res, err := db.GetAdminNotificationByStudentID(studentId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) TestReportByReportID(ctx context.Context, reportId string) (db.Data, error) {
	report, testResult, err := db.GetTestReportByReportID(reportId)
	if err != nil {
		return db.Data{}, err
	}

	var data db.Data
	data.StudentReport = report
	data.TestResults = testResult

	return data, nil
}

func (s *server) MarkAsRead(ctx context.Context, reportId string) (db.StudentReport, error) {
	res, err := db.MarkReportAsRead(reportId)
	if err != nil {
		return db.StudentReport{}, err
	}

	return res, nil
}

func (s *server) MarkNotificationAsRead(ctx context.Context, notificationId string) (db.Notification, error) {
	res, err := db.MarkNotificationAsRead(notificationId)
	if err != nil {
		return db.Notification{}, err
	}

	return res, nil
}

func (s *server) GetTeacherList(ctx context.Context) ([]db.Teacher, error) {
	res, err := db.TeachersList()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) GetTeacherByID(ctx context.Context, teacherId string) (db.Teacher, error) {
	res, err := db.GetTeacherByTeacherID(teacherId)
	if err != nil {
		return db.Teacher{}, err
	}

	return res, nil
}

func (s *server) UpdateTeacherByID(ctx context.Context, input db.TeacherInput) (db.Response, error) {
	res, err := db.UpdateTeacherByTeacherID(input)
	if err != nil {
		return db.Response{}, err
	}

	return res, nil
}

func (s *server) AddNewTeacher(ctx context.Context, input db.TeacherInput) (db.Response, error) {
	if input.Name == "" || input.Password == "" || input.Username == "" || input.Subject == "" {
		return db.Response{}, errors.New("please fill the empty fields")
	}
	res, err := db.AddNewTeacher(input)
	if err != nil {
		return db.Response{}, err
	}

	return res, err
}

func (s *server) AddNewStudent(ctx context.Context, input db.StudentInput, imageData string) (db.Response, error) {
	if input.Username == "" || input.Password == "" || input.Name == "" || input.Username == "" || input.ParentName == "" || input.TeacherID == "" {
		return db.Response{}, errors.New("please fill the empty fields")
	}
	res, err := db.AddNewStudent(input, imageData)
	if err != nil {
		return db.Response{}, err
	}

	return res, err
}

func (s *server) StudentListForAdmin(ctx context.Context) ([]db.StudentList, error) {
	res, err := db.GetStudentListForAdmin()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) StudentDetails(ctx context.Context, studentId string) (db.Student, error) {
	res, err := db.GetStudentDetails(studentId)
	if err != nil {
		return db.Student{}, err
	}

	return res, nil
}

func (s *server) UpdateStudentByStudentID(ctx context.Context, input db.StudentInput) (db.Response, error) {
	res, err := db.UpdateStudentByStudentID(input)
	if err != nil {
		return db.Response{}, err
	}

	return res, nil
}

func (s *server) SendNotificationToStudents(ctx context.Context, notification []db.NotificationInput, message string) (db.Response, error) {
	res, err := db.AddStudentNotification(notification, message)
	if err != nil {
		return db.Response{}, err
	}

	return res, nil
}

func (s *server) AllSentNotifications(ctx context.Context) ([]db.Notification, error) {
	res, err := db.AllNotifications()
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *server) NotificationsStudentList(ctx context.Context, createdAt string) ([]string, error) {
	res, err := db.NotificationsStudentList(createdAt)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)

	for _, v := range res {
		names = append(names, v.StudentName)
	}

	return names, nil
}