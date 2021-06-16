package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/net/idna"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	resp, err := http.Get("http://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("response body read failed: %w", err)
	}
	orig := strings.Split(string(b), "\n")

	tlds := make([]string, 0, len(orig))
	for _, v := range orig[1:] {
		if v != "" {
			tlds = append(tlds, v)
		}
	}

	g, ctx := errgroup.WithContext(ctx)
	for i, v := range tlds {
		i := i
		v := v
		g.Go(func() error {
			var err error
			tlds[i], err = idna.ToUnicode(strings.ToLower(v))
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	sort.Strings(tlds)

	fmt.Printf(`package ursgo

var tlds []string

func init() {
	tlds = []string{"%s"}
}
`, strings.Join(tlds, `", "`))

	return nil
}
