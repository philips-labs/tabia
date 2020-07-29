package graphql

type RepositorySearch struct {
	RepositoryCount int
	PageInfo        PageInfo
	Edges           []Edge
}

type Edge struct {
	Node Node
}

type Node struct {
	Repository Repository `graphql:"... on Repository"`
}
