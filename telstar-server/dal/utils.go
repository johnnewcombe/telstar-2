package dal

import (
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

// HashPassword hashes the password so that it can be stored in the database safely
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// CheckPasswordHash checks the supplied password hash against the specified hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidateUser checks the length, number of capitals, symbols and so on
func ValidateUser(user types.User) bool {

	var (
		userId   int64
		//password int64
		basePage int64
		err      error
	)
	// UserId   = ten digit numeric code e.g. 1000000000 - 9999999999
	// Password = four character numeric pin without leading zeros e.g. 1000 - 9999
	// BasePage = 3 chars or more e.g. 100 - 999999999
	if userId, err = strconv.ParseInt(user.UserId, 10, 64); err != nil ||
		userId < 1000000000 || userId > 9999999999 {
		return false
	}
	// support for non numeric passwords using CheckPasswordStrength() for checking
	err = utils.CheckPasswordStrength(user.Password)
	if err != nil{
		// just issue a warning
		logger.LogWarn.Print(fmt.Sprintf("User %d: %s.", userId, err))
	}
	/*
	if password, err = strconv.ParseInt(user.Password, 10, 64); err != nil ||
		password < 1000 || userId > 9999999999 {
		return false
	}

	if password, err = strconv.Atoi(user.Password); err != nil ||
		password < 1000 || password > 9999 {
		return false
	}
	 */
	if basePage, err = strconv.ParseInt(user.UserId, 10, 64); err != nil ||
		basePage < 100 || basePage > 9999999999 {
		return false
	}
	if len(user.Name) > 23 {
		return false
	}

	return true
}

func getCollectionName(pageNo int, primaryDB bool) (string, error) {

	var prefix string
	var result string

	if primaryDB {
		prefix = "p"
	} else {
		prefix = "s"
	}

	// the above means that if the user page area was 8080 then page 80802014 would be changed to 2014,
	// the code below means that the collection would end up being 201 preceeded by 8080 e.g. 8080201

	if pageNo < 100 {
		result = "000"
	} else {
		for pageNo > 999 {
			pageNo = pageNo / 10 // integer division ??
		}
		result = strconv.Itoa(pageNo)
	}

	return prefix + result, nil

}

// TODO find out what these are and if they are needed!
func primaryCollectionFilter(collection string) bool {
	//TODO: use regex
	return len(collection) == 4 && collection[:1] == "p"
}
func secondaryCollectionFilter(collection string) bool {
	//TODO: use regex
	return len(collection) == 4 && collection[:1] == "s"
}

// GetNewUserId gets a random userId in the range 100000000 - 999999999
func GetNewUserId() (string, error) {
	return "", nil
}
