package model

import "time"

// Secret is a stored credential. Value is encrypted at rest and must never be
// logged, printed, or returned to an AI agent.
type Secret struct {
	Name           string    `json:"name"`
	Value          string    `json:"value"`
	Namespace      string    `json:"namespace"`
	Tags           []string  `json:"tags,omitempty"`
	AllowedTargets []string  `json:"allowed_targets,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	LastUsedAt     time.Time `json:"last_used_at,omitempty"`
}

// SecretMeta is the value-free projection safe to return to agents and logs.
type SecretMeta struct {
	Name      string   `json:"name"`
	Namespace string   `json:"namespace"`
	Tags      []string `json:"tags,omitempty"`
}

// Store is the full collection of secrets, serialized into the encrypted vault.
type Store struct {
	Secrets []Secret `json:"secrets"`
}

// Find returns a pointer to the secret matching name+namespace.
func (s *Store) Find(name, namespace string) (*Secret, bool) {
	for i := range s.Secrets {
		if s.Secrets[i].Name == name && s.Secrets[i].Namespace == namespace {
			return &s.Secrets[i], true
		}
	}
	return nil, false
}

// Upsert inserts or replaces a secret identified by name+namespace.
func (s *Store) Upsert(sec Secret) {
	if existing, ok := s.Find(sec.Name, sec.Namespace); ok {
		sec.CreatedAt = existing.CreatedAt
		*existing = sec
		return
	}
	if sec.CreatedAt.IsZero() {
		sec.CreatedAt = time.Now().UTC()
	}
	s.Secrets = append(s.Secrets, sec)
}

// Meta returns value-free metadata, optionally filtered by namespace ("" = all).
func (s *Store) Meta(namespace string) []SecretMeta {
	out := make([]SecretMeta, 0, len(s.Secrets))
	for _, sec := range s.Secrets {
		if namespace != "" && sec.Namespace != namespace {
			continue
		}
		out = append(out, SecretMeta{Name: sec.Name, Namespace: sec.Namespace, Tags: sec.Tags})
	}
	return out
}
