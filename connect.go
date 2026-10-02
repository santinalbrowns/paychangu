package paychangu

import (
	"context"
	"net/url"
)

// CreateConnectLink generates a PayChangu Connect authorization URL.
func (p *PayChangu) CreateConnectLink(request ConnectLinkRequest) (*ConnectLinkResponse, error) {
	return p.CreateConnectLinkContext(context.Background(), request)
}

// CreateConnectLinkContext generates a Connect link with context.
func (p *PayChangu) CreateConnectLinkContext(ctx context.Context, request ConnectLinkRequest) (*ConnectLinkResponse, error) {
	q := url.Values{}
	q.Set("client_id", request.ClientID)
	q.Set("redirect_uri", request.RedirectURI)
	q.Set("mode", request.Mode)
	if request.Scope != "" {
		q.Set("scope", request.Scope)
	}
	if request.WebhookURL != "" {
		q.Set("wh_url", request.WebhookURL)
	}
	if request.WebhookSecret != "" {
		q.Set("wh_secret", request.WebhookSecret)
	}

	var response ConnectLinkResponse
	if err := p.post(ctx, queryPath("/connect/authorize-link", q), map[string]any{}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetConnectUser retrieves user info for a Connect access token.
func (p *PayChangu) GetConnectUser(accessToken string) (*ConnectUserResponse, error) {
	return p.GetConnectUserContext(context.Background(), accessToken)
}

// GetConnectUserContext retrieves Connect user info with context.
func (p *PayChangu) GetConnectUserContext(ctx context.Context, accessToken string) (*ConnectUserResponse, error) {
	q := url.Values{}
	if accessToken != "" {
		q.Set("access_token", accessToken)
	}
	var response ConnectUserResponse
	if err := p.get(ctx, queryPath("/connect/user", q), &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// RevokeAccessToken revokes a Connect access token.
// Note: official docs list a placeholder path; this uses /connect/revoke.
func (p *PayChangu) RevokeAccessToken(token string) error {
	return p.RevokeAccessTokenContext(context.Background(), token)
}

// RevokeAccessTokenContext revokes a Connect token with context.
func (p *PayChangu) RevokeAccessTokenContext(ctx context.Context, token string) error {
	q := url.Values{}
	q.Set("token", token)
	return p.post(ctx, queryPath("/connect/revoke", q), map[string]any{}, nil)
}
