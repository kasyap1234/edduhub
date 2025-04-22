package auth

import "fmt"

var validNamespaces = map[string]map[string]bool{
	"courses": {
		"faculty":            true,
		"student":            true,
		"admin":              true,
		"manage_qr":          true,
		"view_attendance":    true,
		"mark_attendance":    true,
		"manage_assignments": true,
		"submit_assignments": true,
		"grade_assignments":  true,
	},
	"departments": {
		"head":           true,
		"faculty_member": true,
		"manage_courses": true,
		"view_analytics": true,
	},
	"resources": {
		"owner":    true,
		"viewer":   true,
		"editor":   true,
		"uploader": true,
		"download": true,
	},
	"assignments": {
		"creator":   true,
		"submitter": true,
		"grader":    true,
		"viewer":    true,
	},
	"announcements": {
		"publisher": true,
		"viewer":    true,
		"manager":   true,
	},
}

func validateNamespaceAndRelation(namespace, relation string) error {
	relations, exists := validNamespaces[namespace]
	if !exists {
		return fmt.Errorf("invalid namespace: %s", namespace)
	}

	if !relations[relation] {
		return fmt.Errorf("invalid relation %s for namespace %s", relation, namespace)
	}

	return nil
}
