package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/bpeters-cmu/dns-threat-analyser/graph/generated"
	"github.com/bpeters-cmu/dns-threat-analyser/graph/model"
	"github.com/bpeters-cmu/dns-threat-analyser/pkg/dns"
)

func (r *mutationResolver) Enque(ctx context.Context, ips []string) ([]model.EnqueStatus, error) {
	resultsChan := make(chan model.EnqueStatus, len(ips))
	response := make([]model.EnqueStatus, len(ips))
	for _, ip := range ips {
		go dns.HandleDnsLookup(ip, resultsChan)
	}
	for range resultsChan {
		response = append(response, <-resultsChan)
	}
	return response, nil
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) ([]*model.IP, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
