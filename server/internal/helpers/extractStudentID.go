package helpers

import "github.com/labstack/echo/v4"

func ExtractStudentID(c echo.Context) (int, error) {
	studentID := c.Get("student_id")
	studentIDInt, ok := studentID.(int)
	if !ok {
		return 0, Error(c, "invalid student_id", 400)
	}
	return studentIDInt, nil
}
