package module

import "time"

type Metadata struct {
	Labels      map[string]string
	Annotations map[string]string
	SBOMDigest  string
	BuildID     string
	SourceRef   string
	CreatedAt   time.Time
}

func (m Metadata) Label(key string) string { return m.Labels[key] }
func (m *Metadata) SetLabel(key, value string) {
	if m.Labels == nil {
		m.Labels = map[string]string{}
	}
	m.Labels[key] = value
}
func (m Metadata) Clone() Metadata {
	out := Metadata{SBOMDigest: m.SBOMDigest, BuildID: m.BuildID, SourceRef: m.SourceRef, CreatedAt: m.CreatedAt, Labels: map[string]string{}, Annotations: map[string]string{}}
	for k, v := range m.Labels {
		out.Labels[k] = v
	}
	for k, v := range m.Annotations {
		out.Annotations[k] = v
	}
	return out
}
