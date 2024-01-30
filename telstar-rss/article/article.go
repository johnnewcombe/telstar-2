package article



type Article struct {
	Title        string
	Description  string
	Date         string
}


/*
func (a *Article) Format(cols int) {
	a.Title = formatString(cleanText(a.Title), cols)
	a.Description = formatString(cleanText(a.Description), cols)
	a.Date = formatString(cleanText(a.Date), cols)
}
*/
