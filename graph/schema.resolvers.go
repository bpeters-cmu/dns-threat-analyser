package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/bpeters-cmu/dns-threat-analyser/graph/generated"
	"github.com/bpeters-cmu/dns-threat-analyser/graph/model"
	"github.com/bpeters-cmu/dns-threat-analyser/pkg/database"
	"github.com/bpeters-cmu/dns-threat-analyser/pkg/dns"
)

const (
	notFoundError = "NOT_FOUND"
)

func (r *mutationResolver) Enque(ctx context.Context, ips []string) ([]model.Status, error) {
	resultsChan := make(chan model.Status, len(ips))
	response := make([]model.Status, 0)
	for _, ip := range ips {
		go dns.HandleDnsLookup(ip, &database.SqliteDB{}, &dns.DigClient{}, resultsChan)
	}
	for i := 0; i < cap(resultsChan); i++ {
		response = append(response, <-resultsChan)
	}
	return response, nil
}

func (r *queryResolver) GetIPDetails(ctx context.Context, ip string) (model.Status, error) {
	err := dns.ValidateIp(ip)
	if err != nil {
		return nil, err
	}
	db := database.SqliteDB{}
	ipDetails, err := db.GetIp(ip)
	if err != nil {
		return model.ErrorStatus{Error: &model.Error{IPAddress: ip, ErrorMessage: err.Error(), ErrorCode: notFoundError}}, nil
	}
	return model.SuccessStatus{IP: ipDetails}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
