package auth

import "context"

type PermissionHelper struct {
	keto KetoService
}

func NewPermissionHelper(keto KetoService) *PermissionHelper {
	return &PermissionHelper{keto: keto}
}

// Course-related permission checks
func (p *PermissionHelper) CanManageQR(ctx context.Context, facultyID, courseID string) (bool, error) {
	return p.keto.CheckPermission(ctx, "courses", facultyID, "manage_qr", courseID)
}

func (p *PermissionHelper) CanViewAttendance(ctx context.Context, userID, courseID string) (bool, error) {
	return p.keto.CheckPermission(ctx, "courses", userID, "view_attendance", courseID)
}

// Department-related permission checks
func (p *PermissionHelper) CanManageCourses(ctx context.Context, facultyID, departmentID string) (bool, error) {
	return p.keto.CheckPermission(ctx, "departments", facultyID, "manage_courses", departmentID)
}

// Resource-related permission checks
func (p *PermissionHelper) CanAccessResource(ctx context.Context, userID, resourceID, action string) (bool, error) {
	return p.keto.CheckPermission(ctx, "resources", userID, action, resourceID)
}
