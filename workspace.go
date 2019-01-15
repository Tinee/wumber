package wumber

type Workspace struct {
	ID   string `dynamodbav:"PK"`
	Name string `dynamodbav:"Name"`
}
