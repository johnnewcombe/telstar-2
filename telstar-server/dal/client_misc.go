package dal

func BackupCollection(connectionUrl string, collection string, destination_directory string, authUser string) error {

	// FIXME This is not backup code, just demo code.
	/*
		pNames, sNames, err := getCollectionNames(connectionUrl)
		if err != nil {
			return err
		}
		fmt.Println(pNames)
		fmt.Println(sNames)
	*/
	return nil
}

// TODO These are not used in Telstar 1.x, are they needed!
//func createPkIndex(collection string)                                  {} // not used
//func dropPkIndex(collection string)                                    {} // not used
//func createDirectory(directory string)                                 {} // should this be here? only used by backupDatabase amd that isn't used
//func BackupDatabase(destination_directory string)                      {} // not used ???
//func CreatePkIndex()                                                   {} // not used
//func InsertDocumentFromPath(pathToDoc string, primaryDb bool)          {} // not used
//func GetAllStaticDocuments(collection string)                          {} // not used
//func GetAllNonStaticDocuments(collection string)                       {} // not used
