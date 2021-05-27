package csv

type Output struct {
	Column    string `csv:"column"`
	Number    int    `csv:"number"`
	Title     string `csv:"title"`
	URL       string `csv:"url"`
	Path      string `csv:"path"`
	Assignees string `csv:"assignees"`
	Labels    string `csv:"labels"`
}
