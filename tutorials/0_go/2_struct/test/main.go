package main

type Student struct {
    Name string
    ID   int
}

func (s Student) GetName() string {
    return s.Name
}

func (s Student) GetID() int {
    return s.ID
}

type InfoProvider interface {
    GetName() string
    GetID() int
}
