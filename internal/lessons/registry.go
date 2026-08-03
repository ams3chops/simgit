package lessons

import "simgit/internal/lesson"

// All returns all lessons in order.
func All() []lesson.Lesson {
	return []lesson.Lesson{
		lessonConfig(),
		lessonInit(),
		lessonStaging(),
		lessonCommit(),
		lessonLog(),
		lessonBranch(),
		lessonSwitch(),
		lessonMerge(),
		lessonCheckout(),
		lessonRemote(),
	}
}
