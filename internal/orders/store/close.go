package store

func (s *Store) Close() { s.pool.Close() }
