package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/example/wasm-sandbox-runtime/internal/application"
	"github.com/example/wasm-sandbox-runtime/internal/domain/execution"
	"github.com/example/wasm-sandbox-runtime/internal/domain/module"
	"github.com/example/wasm-sandbox-runtime/internal/domain/policy"
	"github.com/example/wasm-sandbox-runtime/internal/domain/quota"
	"github.com/example/wasm-sandbox-runtime/pkg/api"
)

type Server struct {
	service *application.Service
	log     *slog.Logger
	maxBody int64
	mux     *http.ServeMux
}

func New(s *application.Service, l *slog.Logger) *Server {
	h := &Server{service: s, log: l, maxBody: 4 << 20, mux: http.NewServeMux()}
	h.routes()
	return h
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	s.mux.ServeHTTP(w, r)
	s.log.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
}
func (s *Server) routes() {
	m := s.mux
	m.HandleFunc("/healthz", s.health)
	m.HandleFunc("/readyz", s.health)
	m.HandleFunc("/metrics", s.metrics)
	m.HandleFunc("/api/v1/modules", s.modules)
	m.HandleFunc("/api/v1/modules/", s.moduleSub)
	m.HandleFunc("/api/v1/executions", s.executions)
	m.HandleFunc("/api/v1/executions/", s.executionSub)
	m.HandleFunc("/api/v1/policies", s.policies)
	m.HandleFunc("/api/v1/quotas", s.quotas)
	m.HandleFunc("/api/v1/runtime", s.runtimeInfo)
	m.HandleFunc("/api/v1/admission-reviews", s.notImplemented)
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	jsonWrite(w, 200, map[string]any{"status": "ok", "time": time.Now().UTC()})
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintln(w, "sandbox_http_requests_total 1")
}
func (s *Server) modules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tenant := r.URL.Query().Get("tenant_id")
		offset, limit := page(r)
		v, e := s.service.ListModules(r.Context(), tenant, offset, limit)
		if e != nil {
			writeErr(w, e)
			return
		}
		jsonWrite(w, 200, map[string]any{"items": v, "offset": offset, "limit": limit})
	case http.MethodPost:
		var in module.Module
		if !decode(w, r, &in, s.maxBody) {
			return
		}
		v, e := s.service.CreateModule(r.Context(), in)
		if e != nil {
			writeErr(w, e)
			return
		}
		jsonWrite(w, 201, v)
	default:
		w.WriteHeader(405)
	}
}
func (s *Server) moduleSub(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeErr(w, api.NotFound("module route"))
		return
	}
	id := parts[3]
	if len(parts) == 4 && r.Method == http.MethodGet {
		v, e := s.service.GetModule(r.Context(), id)
		if e != nil {
			writeErr(w, e)
			return
		}
		jsonWrite(w, 200, v)
		return
	}
	if len(parts) >= 5 && parts[4] == "versions" {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			return
		}
		var v module.Version
		if !decode(w, r, &v, s.maxBody) {
			return
		}
		out, e := s.service.AddVersion(r.Context(), id, v)
		if e != nil {
			writeErr(w, e)
			return
		}
		jsonWrite(w, 201, out)
		return
	}
	if len(parts) >= 5 && parts[4] == "inspect" {
		s.inspect(w, r, id)
		return
	}
	writeErr(w, api.NotFound("module route"))
}
func (s *Server) inspect(w http.ResponseWriter, r *http.Request, id string) {
	v, e := s.service.GetModule(r.Context(), id)
	if e != nil {
		writeErr(w, e)
		return
	}
	if len(v.Versions) == 0 {
		writeErr(w, api.Invalid("module has no versions"))
		return
	}
	data, e := io.ReadAll(http.MaxBytesReader(w, r.Body, s.maxBody))
	if e != nil {
		writeErr(w, e)
		return
	}
	out, e := s.service.Inspect(r.Context(), v.Versions[len(v.Versions)-1], data)
	if e != nil {
		writeErr(w, e)
		return
	}
	jsonWrite(w, 200, out)
}
func (s *Server) executions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		offset, limit := page(r)
		v, e := s.service.ListExecutions(r.Context(), r.URL.Query().Get("tenant_id"), offset, limit)
		if e != nil {
			writeErr(w, e)
			return
		}
		jsonWrite(w, 200, map[string]any{"items": v, "offset": offset, "limit": limit})
	case http.MethodPost:
		var in execution.Execution
		if !decode(w, r, &in, s.maxBody) {
			return
		}
		v, e := s.service.Submit(r.Context(), in)
		if e != nil {
			writeErr(w, e)
			return
		}
		jsonWrite(w, 202, v)
	default:
		w.WriteHeader(405)
	}
}
func (s *Server) executionSub(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeErr(w, api.NotFound("execution route"))
		return
	}
	id := parts[3]
	if len(parts) == 4 && r.Method == http.MethodGet {
		v, e := s.service.GetExecution(r.Context(), id)
		if e != nil {
			writeErr(w, e)
			return
		}
		jsonWrite(w, 200, v)
		return
	}
	if len(parts) >= 5 {
		switch parts[4] {
		case "cancel":
			if r.Method != http.MethodPost {
				w.WriteHeader(405)
				return
			}
			if e := s.service.Cancel(r.Context(), id); e != nil {
				writeErr(w, e)
				return
			}
			jsonWrite(w, 200, map[string]string{"status": "cancelled"})
		case "checkpoint":
			if r.Method != http.MethodPost {
				w.WriteHeader(405)
				return
			}
			v, e := s.service.Checkpoint(r.Context(), id)
			if e != nil {
				writeErr(w, e)
				return
			}
			jsonWrite(w, 200, v)
		case "events":
			s.events(w, r, id)
		case "logs":
			jsonWrite(w, 200, map[string]any{"items": []string{"execution logs are redacted by policy"}})
		case "result":
			v, e := s.service.GetExecution(r.Context(), id)
			if e != nil {
				writeErr(w, e)
				return
			}
			jsonWrite(w, 200, map[string]any{"digest": v.ResultDigest, "state": v.State})
		default:
			writeErr(w, api.NotFound("execution operation"))
		}
		return
	}
	writeErr(w, api.NotFound("execution route"))
}
func (s *Server) events(w http.ResponseWriter, r *http.Request, id string) {
	ch, closeFn := s.service.Events().Subscribe(id)
	defer closeFn()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	f, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, api.Internal("stream unsupported"))
		return
	}
	for {
		select {
		case <-r.Context().Done():
			return
		case e := <-ch:
			b, _ := json.Marshal(e)
			fmt.Fprintf(w, "event: state\ndata: %s\n\n", b)
			f.Flush()
		}
	}
}
func (s *Server) policies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var in policy.Policy
	if !decode(w, r, &in, s.maxBody) {
		return
	}
	v, e := s.service.PutPolicy(r.Context(), in)
	if e != nil {
		writeErr(w, e)
		return
	}
	jsonWrite(w, 201, v)
}
func (s *Server) quotas(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("tenant_id")
		v, e := s.service.GetQuota(r.Context(), id)
		if e != nil {
			writeErr(w, e)
			return
		}
		jsonWrite(w, 200, v)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var in quota.Quota
	if !decode(w, r, &in, s.maxBody) {
		return
	}
	v, e := s.service.PutQuota(r.Context(), in)
	if e != nil {
		writeErr(w, e)
		return
	}
	jsonWrite(w, 201, v)
}
func (s *Server) runtimeInfo(w http.ResponseWriter, _ *http.Request) {
	jsonWrite(w, 200, map[string]any{"adapter": "go-safe-simulator", "capabilities": []string{"inspect", "execute", "cancel", "checkpoint"}})
}
func moduleStatus(m module.Module) int {
	if m.LifecycleReady() {
		return 200
	}
	return 409
}
func (s *Server) notImplemented(w http.ResponseWriter, _ *http.Request) {
	writeErr(w, api.Error{Code: "not_implemented", Message: "admission review adapter is not configured"})
}
func page(r *http.Request) (int, int) {
	o, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	l, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if o < 0 {
		o = 0
	}
	if l <= 0 || l > 100 {
		l = 20
	}
	return o, l
}
func decode(w http.ResponseWriter, r *http.Request, v any, max int64) bool {
	r.Body = http.MaxBytesReader(w, r.Body, max)
	defer r.Body.Close()
	if e := json.NewDecoder(r.Body).Decode(v); e != nil {
		writeErr(w, api.Invalid("invalid JSON: "+e.Error()))
		return false
	}
	return true
}
func jsonWrite(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, e error) {
	status := 500
	v := api.Internal(e.Error())
	if errors.Is(e, context.Canceled) {
		status = 499
		v = api.Error{Code: "cancelled", Message: e.Error()}
	}
	if strings.Contains(e.Error(), "required") || strings.Contains(e.Error(), "invalid") || strings.Contains(e.Error(), "limit") {
		status = 400
		v = api.Invalid(e.Error())
	}
	if strings.Contains(e.Error(), "not found") {
		status = 404
		v = api.NotFound(e.Error())
	}
	jsonWrite(w, status, v)
}
