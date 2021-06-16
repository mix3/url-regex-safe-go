update: goinstall
	gentlds > tlds.gen.go


goinstall:
	go install ./cmd/...
