package store

// Close releases the pool owned by Store.
func (s *Store) Close() { s.pool.Close() }
