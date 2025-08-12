package domain

type BlogWithPopValue struct {
	PopularityValue int
	Blog            *BlogDTO
}
type ByPopularityValue []BlogWithPopValue

// Len is the number of elements in the collection.
func (a ByPopularityValue) Len() int {
    return len(a)
}

// Less reports whether the element with index i should sort before the element with index j.
func (a ByPopularityValue) Less(i, j int) bool {
    return a[i].PopularityValue < a[j].PopularityValue
}

// Swap swaps the elements with indexes i and j.
func (a ByPopularityValue) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}


type ByPopularityValueDesc []BlogWithPopValue
// Len is the number of elements in the collection.
func (a ByPopularityValueDesc) Len() int {
    return len(a)
}

// Less reports whether the element with index i should sort before the element with index j.
func (a ByPopularityValueDesc) Less(i, j int) bool {
    return a[i].PopularityValue > a[j].PopularityValue
}

// Swap swaps the elements with indexes i and j.
func (a ByPopularityValueDesc) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}
