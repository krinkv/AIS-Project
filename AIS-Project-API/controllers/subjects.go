package controllers

import (
	"AIS-Project-API/database"
	"AIS-Project-API/utils/token"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GradeInput struct {
	CourseId  uint   `json:"courseId"`
	StudentId uint   `json:"studentId"`
	Grade     uint64 `json:"grade"`
}

func EditGrade(c *gin.Context) {
	var gradeInput GradeInput

	if err := c.ShouldBindJSON(&gradeInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade := database.Grade{
		StudentId: gradeInput.StudentId,
		CourseId:  gradeInput.CourseId,
		Grade:     gradeInput.Grade,
	}

	_, err := grade.Edit()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "grade edited successfully",
	})
}

type EnrollInput struct {
	StudentId string `json:"studentId" binding:"required"`
	CourseId  string `json:"courseId" binding:"required"`
}

func EnrollCourse(c *gin.Context) {
	teacherId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var enrollInput EnrollInput
	if err := c.ShouldBindJSON(&enrollInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	studentId, err := strconv.ParseInt(enrollInput.StudentId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	courseId, err := strconv.ParseInt(enrollInput.CourseId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	enrollment := database.Enrollment{
		StudentId: uint(studentId),
		TeacherId: teacherId,
		CourseId:  uint(courseId),
	}

	_, err = enrollment.Enroll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "enrolled successfully",
	})
}

func TeacherGrades(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var enrollment []database.Enrollment
	database.DB.Where("teacher_id = ?", user_id).Find(&enrollment)

	type ResultType struct {
		Student database.Student
		Grade   database.Grade
	}

	var grades []database.Grade
	for _, value := range enrollment {
		var tempGrades []database.Grade
		database.DB.Where("student_id = ?", value.StudentId).Find(&tempGrades)
		grades = append(grades, tempGrades...)
	}

	fmt.Println(grades)
	c.JSON(http.StatusOK, gin.H{"data": grades})

	// var grades []database.Grade
	// database.DB.Find(&grades)

	// var teacherGrades []database.Grade
	// for _, currentGrade := range grades {
	// 	if currentGrade.CourseId == user.ID {
	// 		teacherGrades = append(teacherGrades, currentGrade)
	// 	}
	// }

	// c.JSON(http.StatusOK, gin.H{"data": teacherGrades})
}

// за учител връща неговите курсове?
func CoursesPerTeacher(c *gin.Context) {
	teacherId, err := strconv.ParseUint(c.Param("teacherId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error while parsing parameter: ": err.Error()})
		return
	}

	var courses []database.Course
	database.DB.Where("teacher_id = ?", teacherId).Find(&courses)

	c.JSON(http.StatusOK, gin.H{"data": courses})
}

// for teacher's course return students and grades

// Returns array of Students that are enrolled in the subject passed in the url
func StudentsPerCourse(c *gin.Context) {
	subjectId, err := strconv.ParseUint(c.Param("subjectId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error while parsing parameter: ": err})
		return
	}

	var enrollments []database.Enrollment
	database.DB.Where("course_id = ?", subjectId).Find(&enrollments)

	var students []database.Student
	database.DB.Find(&students)

	var foundStudents []database.Student
	for _, currentEnrollments := range enrollments {
		for _, currentStudent := range students {
			if currentStudent.UserId == uint(currentEnrollments.StudentId) {
				foundStudents = append(foundStudents, currentStudent)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": foundStudents})
}

// Returns array of ResultGrade - tuple of CourseName and GradeNumber
func StudentGrades(c *gin.Context) {
	studentId, err := strconv.ParseUint(c.Param("studentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error while parsing parameter: ": err.Error()})
		return
	}

	var grades []database.Grade
	database.DB.Where("student_id = ?", studentId).Find(&grades)

	var courses []database.Course
	database.DB.Find(&courses)

	type ResultGrade struct {
		CourseName  string
		GradeNumber uint
	}

	var result []ResultGrade
	for _, currentGrade := range grades {
		for _, currentCourse := range courses {
			if currentCourse.ID == currentGrade.CourseId {
				result = append(result, ResultGrade{
					currentCourse.Name, uint(currentGrade.Grade),
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
