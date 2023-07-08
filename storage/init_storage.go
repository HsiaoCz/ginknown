package storage

type InitStore interface {
	InitStorage() error
}

type IStore struct {
	Iss []InitStore
}

func NewIStore() *IStore {
	return &IStore{
		Iss: make([]InitStore, 0),
	}
}

func (s *IStore) StartConn(iss ...InitStore) (err error) {
	s.Iss = append(s.Iss, iss...)
	for _, init := range s.Iss {
		if err = init.InitStorage(); err != nil {
			return
		}
	}
	return
}
