package main

import "net/http"

type apiError struct {
	Msg    string
	Status int
}

var ErrMethodNotAllowed = apiError{
	Msg:    "Invalid Method!",
	Status: http.StatusMethodNotAllowed,
}

var ErrWriteJSON = apiError{
	Msg:    "Couldn't marshall object into JSON format!",
	Status: http.StatusInternalServerError,
}

var ErrReadJSON = apiError{
	Msg:    "Couldn't unmarshall JSON into object!",
	Status: http.StatusInternalServerError,
}

var ErrQueryDatabase = apiError{
	Msg:    "Querying the database resulted in an error!",
	Status: http.StatusInternalServerError,
}

var ErrStrToInt = apiError{
	Msg:    "Couldn't convert the String to an Int!",
	Status: http.StatusInternalServerError,
}

var ErrIDNotSet = apiError{
	Msg:    "The ID is not set!",
	Status: http.StatusInternalServerError,
}

var ErrIDNegative = apiError{
	Msg:    "Negative IDs are invalid!",
	Status: http.StatusInternalServerError,
}

var ErrNameNotSet = apiError{
	Msg:    "The category name so be searched for is empty!",
	Status: http.StatusInternalServerError,
}
