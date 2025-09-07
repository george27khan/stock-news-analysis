package domain

type Browser struct {
	Host    string
	Port    int
	ExePath string
	RunArgs []string
}
