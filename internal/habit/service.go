package habit

import "errors"

var ErrNotFound = errors.New("habit not found")
var ErrValidation = errors.New("name is required")

// store is the in-memory backing storage. This will be replaced by a
// real repository (backed by Postgres) in the next milestone — nothing
// outside this file will need to change when that happens.
var store = []Habit{}
var nextID = 1

func List() []Habit {
	return store
}

func Create(input Habit) (Habit, error) {
	if input.Name == "" {
		return Habit{}, ErrValidation
	}

	input.ID = nextID
	nextID++
	store = append(store, input)
	return input, nil
}

func Get(id int) (Habit, error) {
	for _, h := range store {
		if h.ID == id {
			return h, nil
		}
	}
	return Habit{}, ErrNotFound
}

func Update(id int, input Habit) (Habit, error) {
	if input.Name == "" {
		return Habit{}, ErrValidation
	}

	for i, h := range store {
		if h.ID == id {
			input.ID = id
			store[i] = input
			return input, nil
		}
	}
	return Habit{}, ErrNotFound
}

func Delete(id int) error {
	for i, h := range store {
		if h.ID == id {
			store = append(store[:i], store[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
