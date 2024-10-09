package gofilemirror

import (
	"time"
)

type OperationList []*Operation

func (ol *OperationList) WaitForStart(duration time.Duration) {
	currentDuration := time.Duration(0)

	for {
		started := 0

		for _, operation := range *ol {
			if operation.started {
				started++
			}
		}

		if started == len(*ol) {
			return
		}

		sleepDuration := 10 * time.Millisecond

		time.Sleep(sleepDuration)
		currentDuration += sleepDuration

		if currentDuration >= duration {
			return
		}
	}
}

func (ol *OperationList) WaitForDone(duration time.Duration) {
	currentDuration := time.Duration(0)

	for {
		done := 0

		for _, operation := range *ol {
			if operation.done {
				done++
			}
		}

		if done == len(*ol) {
			return
		}

		sleepDuration := 10 * time.Millisecond

		time.Sleep(sleepDuration)
		currentDuration += sleepDuration

		if currentDuration >= duration {
			return
		}
	}
}

func (ol *OperationList) GetFirstAsyncOperation() *Operation {
	for _, operation := range *ol {
		if operation.async {
			return operation
		}
	}

	return nil
}

func (ol *OperationList) GetFirstNonAsyncOperation() *Operation {
	for _, operation := range *ol {
		if !operation.async {
			return operation
		}
	}

	return nil
}

func (ol *OperationList) GetAsyncOperations() *OperationList {
	operationList := OperationList{}

	for _, operation := range *ol {
		if operation.async {
			operationList = append(operationList, operation)
		}
	}

	return &operationList
}

func (ol *OperationList) GetNonAsyncOperations() *OperationList {
	operationList := OperationList{}

	for _, operation := range *ol {
		if !operation.async {
			operationList = append(operationList, operation)
		}
	}

	return &operationList
}

func (ol *OperationList) GetPendingOperations() *OperationList {
	operationList := OperationList{}

	for _, operation := range *ol {
		if !operation.done {
			operationList = append(operationList, operation)
		}
	}

	return &operationList
}
