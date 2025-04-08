package helpers

import "hash/fnv"

func RollNoToStudentID(rollNo string) (int, error) {
	hasher := fnv.New64a()
	hasher.Write([]byte(rollNo))
	studentID := int(hasher.Sum64())
	return studentID, nil

}
