# go-talks

Code and slides from Golang Ghana talk sessions.

## Sessions

- **session2** — P2P NAT traversal with Go (slides)
- **session3** — Building APIs with GORM and sqlc (code + slides)
- **session4** — Testing Kubernetes with Go (e2e-framework)
- **session6** — Go tooling: deadcode, govulncheck (slides)
- **session7** — Understanding `context` in Go (cancel, cause, deadlines, timeouts)

## Layout

Each session is self-contained in its own directory. Code sessions ship a `go.mod`; slide-only sessions just hold the deck.

## Run

```bash
cd session<N>
go run ./...
```

See per-session `README.md` / `SETUP.md` where present (e.g. `session3/SETUP.md`, `session3/API_README.md`).

## License

GPL-3.0 — see [LICENSE](LICENSE).
