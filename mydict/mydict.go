package mydict

import "errors"

// Error
var (
	errNotFound   = errors.New("Not Found")
	errCantUpdate = errors.New("Cant update non-existing word")
	errIsExist    = errors.New("It's already exist")
)

// Dictionary type
type Dictionary map[string]string

// Search : Search for a word
func (d Dictionary) Search(word string) (string, error) {
	value, exist := d[word] // value : string, exist : boolean
	if exist {
		return value, nil
	}
	return "", errNotFound
}

// Add : Add word to Dictionary
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errIsExist
	}
	return nil
}

// Update : Update
func (d Dictionary) Update(word, definition string) error {
	_, err := d.Search(word)
	switch err {
	case nil:
		d[word] = definition
	case errNotFound:
		return errCantUpdate
	}
	return nil
}

//Delete : Delete word from Dictionary
func (d Dictionary) Delete(word string) {
	delete(d, word)
}
