package businesserrors

import "errors"

var (
	ErrCatNotFound     = errors.New("cat with this id not found")
	ErrBreedIsNotValid = errors.New("breed is not valid")
)

var (
	ErrMissionNotFound = errors.New("mission with this id not found")
	ErrCatIsAssigned   = errors.New("a mission cannot be deleted, it is already assigned to a cat")
)

var (
	ErrTargetNotExists                   = errors.New("target not exists")
	ErrTargetNotesCantBeUpdated          = errors.New("cannot be updated if either the target or the mission is completed")
	ErrTargetCanNotBeDeleted             = errors.New("a target cannot be deleted if it is already completed")
	ErrMaxNumberOfTargetsExeeded         = errors.New("max number of targets exeeded")
	ErrCannotAddTargetToCompletedMission = errors.New("cannot add target to completed mission")
)
