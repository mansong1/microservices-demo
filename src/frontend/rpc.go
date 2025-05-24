// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"time"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/frontend/genproto"

	"github.com/pkg/errors"
)

const (
	avoidNoopCurrencyConversionRPC = false
)

func (fe *frontendServer) getCurrencies(ctx context.Context) ([]string, error) {
	currs, err := pb.NewCurrencyServiceClient(fe.currencySvcConn).
		GetSupportedCurrencies(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	var out []string
	for _, c := range currs.CurrencyCodes {
		if _, ok := whitelistedCurrencies[c]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

func (fe *frontendServer) getProducts(ctx context.Context) ([]*pb.Product, error) {
	resp, err := pb.NewProductCatalogServiceClient(fe.productCatalogSvcConn).
		ListProducts(ctx, &pb.Empty{})
	return resp.GetProducts(), err
}

func (fe *frontendServer) getProduct(ctx context.Context, id string) (*pb.Product, error) {
	start := time.Now()
	resp, err := pb.NewProductCatalogServiceClient(fe.productCatalogSvcConn).
		GetProduct(ctx, &pb.GetProductRequest{Id: id})
	duration := time.Since(start)
	
	status := "success"
	if err != nil {
		status = "error"
	}
	recordGRPCRequest("ProductCatalogService", "GetProduct", status, duration)
	
	return resp, err
}

func (fe *frontendServer) getCart(ctx context.Context, userID string) ([]*pb.CartItem, error) {
	start := time.Now()
	resp, err := pb.NewCartServiceClient(fe.cartSvcConn).GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	duration := time.Since(start)
	
	status := "success"
	if err != nil {
		status = "error"
	}
	recordGRPCRequest("CartService", "GetCart", status, duration)
	
	return resp.GetItems(), err
}

func (fe *frontendServer) emptyCart(ctx context.Context, userID string) error {
	start := time.Now()
	_, err := pb.NewCartServiceClient(fe.cartSvcConn).EmptyCart(ctx, &pb.EmptyCartRequest{UserId: userID})
	duration := time.Since(start)
	
	status := "success"
	if err != nil {
		status = "error"
	}
	recordGRPCRequest("CartService", "EmptyCart", status, duration)
	
	return err
}

func (fe *frontendServer) insertCart(ctx context.Context, userID, productID string, quantity int32) error {
	start := time.Now()
	_, err := pb.NewCartServiceClient(fe.cartSvcConn).AddItem(ctx, &pb.AddItemRequest{
		UserId: userID,
		Item: &pb.CartItem{
			ProductId: productID,
			Quantity:  quantity},
	})
	duration := time.Since(start)
	
	status := "success"
	if err != nil {
		status = "error"
	}
	recordGRPCRequest("CartService", "AddItem", status, duration)
	
	return err
}

func (fe *frontendServer) convertCurrency(ctx context.Context, money *pb.Money, currency string) (*pb.Money, error) {
	if avoidNoopCurrencyConversionRPC && money.GetCurrencyCode() == currency {
		return money, nil
	}
	
	start := time.Now()
	resp, err := pb.NewCurrencyServiceClient(fe.currencySvcConn).
		Convert(ctx, &pb.CurrencyConversionRequest{
			From:   money,
			ToCode: currency})
	duration := time.Since(start)
	
	status := "success"
	if err != nil {
		status = "error"
	}
	recordGRPCRequest("CurrencyService", "Convert", status, duration)
	
	return resp, err
}

func (fe *frontendServer) getShippingQuote(ctx context.Context, items []*pb.CartItem, currency string) (*pb.Money, error) {
	start := time.Now()
	quote, err := pb.NewShippingServiceClient(fe.shippingSvcConn).GetQuote(ctx,
		&pb.GetQuoteRequest{
			Address: nil,
			Items:   items})
	duration := time.Since(start)
	
	status := "success"
	if err != nil {
		status = "error"
	}
	recordGRPCRequest("ShippingService", "GetQuote", status, duration)
	
	if err != nil {
		return nil, err
	}
	localized, err := fe.convertCurrency(ctx, quote.GetCostUsd(), currency)
	return localized, errors.Wrap(err, "failed to convert currency for shipping cost")
}

func (fe *frontendServer) getRecommendations(ctx context.Context, userID string, productIDs []string) ([]*pb.Product, error) {
	start := time.Now()
	resp, err := pb.NewRecommendationServiceClient(fe.recommendationSvcConn).ListRecommendations(ctx,
		&pb.ListRecommendationsRequest{UserId: userID, ProductIds: productIDs})
	duration := time.Since(start)
	
	status := "success"
	if err != nil {
		status = "error"
		recordGRPCRequest("RecommendationService", "ListRecommendations", status, duration)
		return nil, err
	}
	recordGRPCRequest("RecommendationService", "ListRecommendations", status, duration)
	
	out := make([]*pb.Product, len(resp.GetProductIds()))
	for i, v := range resp.GetProductIds() {
		p, err := fe.getProduct(ctx, v)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get recommended product info (#%s)", v)
		}
		out[i] = p
	}
	if len(out) > 4 {
		out = out[:4] // take only first four to fit the UI
	}
	return out, err
}

func (fe *frontendServer) getAd(ctx context.Context, ctxKeys []string) ([]*pb.Ad, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	start := time.Now()
	resp, err := pb.NewAdServiceClient(fe.adSvcConn).GetAds(ctx, &pb.AdRequest{
		ContextKeys: ctxKeys,
	})
	duration := time.Since(start)
	
	status := "success"
	if err != nil {
		status = "error"
	}
	recordGRPCRequest("AdService", "GetAds", status, duration)
	
	return resp.GetAds(), errors.Wrap(err, "failed to get ads")
}
