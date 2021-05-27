package v4client

type LabelEdge struct {
	Cursor string    `json:"cursor"`
	Node   LabelNode `json:"node"`
}

type LabelNode struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

type AssigneesEdge struct {
	Cursor string       `json:"cursor"`
	Node   AssigneeNode `json:"node"`
}

type AssigneeNode struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatarUrl"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	Url       string `json:"url"`
}

type CardEdge struct {
	Cursor string   `json:"cursor"`
	Node   CardNode `json:"node"`
}

type CardNode struct {
	Content    Issue  `json:"content"`
	Note       string `json:"note"`
	IsArchived bool   `json:"isArchived"`
	ID         string `json:"id"`
}

type ColumnEdge struct {
	Cursor string     `json:"cursor"`
	Node   ColumnNode `json:"node"`
}

type ColumnNode struct {
	Name  string `json:"name"`
	Cards struct {
		Edges []CardEdge `json:"edges"`
	} `json:"cards"`
	ID string `json:"id"`
}

type ProjectEdge struct {
	Cursor string      `json:"cursor"`
	Node   ProjectNode `json:"node"`
}

type ProjectNode struct {
	Columns struct {
		Edges []ColumnEdge `json:"edges"`
	} `json:"columns"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Issue struct {
	Title  string `json:"title"`
	Number int    `json:"number"`
	Url    string `json:"url"`
	Body   string `json:"body"`
	Labels struct {
		Edges []LabelEdge `json:"edges"`
	} `json:"labels"`
	Assignees struct {
		Edges []AssigneesEdge `json:"edges"`
	} `json:"assignees"`
}

type Repository struct {
	Name     string `json:"name"`
	Projects struct {
		Edges []ProjectEdge `json:"edges"`
	} `json:"projects"`
}
