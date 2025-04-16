package helpers

import "github.com/labstack/echo/v4"

func ExtractCollegeID(c echo.Context) (int, error) {

	collegeID := c.Get("college_id")
	collegeIDInt, ok := collegeID.(int)
	if !ok {
		return 0, echo.NewHTTPError(400, "Invalid college_id")

	}

	return collegeIDInt, nil
}
