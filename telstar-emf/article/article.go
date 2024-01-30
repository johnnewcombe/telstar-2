package article


type Article struct {
	Title        string `json:"title" bson:"title"`
	Description  string `json:"description" bson:"description"`
	StartDate    string `json:"start_date" bson:"start_date"` // FIXME consider combining start and end into one field
	EndDate    string `json:"end_date" bson:"end_date"`
	StartTime    string `json:"start_time" bson:"start_time"`
	EndTime    string `json:"end_time" bson:"end_time"`
	Venue string `json:"venue" bson:"venue"`
}


/*
func (a *Article) Format(cols int) {
	a.Title = formatString(cleanText(a.Title), cols)
	a.Description = formatString(cleanText(a.Description), cols)
	a.Date = formatString(cleanText(a.Date), cols)
}
*/
