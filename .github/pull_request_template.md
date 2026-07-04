## Summary

<!-- What does this PR change, and why? Link any related issue with "Closes #123". -->

## Type of change

<!-- Check all that apply. -->

- [ ] Bug fix (`fix`)
- [ ] New feature (`feat`)
- [ ] Refactor (`refactor`)
- [ ] Performance (`perf`)
- [ ] Documentation (`docs`)
- [ ] Tests (`test`)
- [ ] Build / CI (`build` / `ci`)

## Checklist

- [ ] My commit messages follow [Conventional Commits](CONTRIBUTING.md#commit-messages--conventional-commits)
- [ ] `gofmt -s -l .` reports no changes
- [ ] `go vet ./internal/...` passes
- [ ] `go test ./internal/...` passes
- [ ] `go test -race ./internal/reactive` passes (if I touched the reactive core)
- [ ] If this affects the WASM build: `cd examples/simple && GOOS=js GOARCH=wasm go build -o /tmp/app.wasm main.go` succeeds
- [ ] I did **not** add any third-party dependencies (Maya is standard-library only)
- [ ] I updated documentation where relevant

## Notes for reviewers

<!-- Anything that helps the review: tricky bits, follow-ups, screenshots, etc. -->
