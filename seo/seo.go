package seo


type SEO struct {
	Name       string
	Varibles   []string
	Context    func(...interface{}) map[string]string
	collection *Collection
}