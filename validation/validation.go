package validation


func foldSourceData () error {
	return nil
}

func DoPerformKFoldValidation (numberOfSlots int) error {

	if numberOfSlots < 1 {
		return KErrorInvalidNumberOfFolds
	}

	err := foldSourceData ()

	return err
}