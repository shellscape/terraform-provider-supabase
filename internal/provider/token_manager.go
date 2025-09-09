package provider

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/supabase/cli/pkg/api"
)

// TokenManager handles automatic token exchange for different Supabase APIs
type TokenManager struct {
	managementClient *api.ClientWithResponses
	managementToken  string
	
	// Cache for project tokens
	projectTokens map[string]*ProjectTokens
	mutex         sync.RWMutex
}

// ProjectTokens holds the different tokens available for a project
type ProjectTokens struct {
	ServiceRoleKey string    // JWT token for storage/database operations
	AnonKey        string    // Anonymous JWT token  
	CachedAt       time.Time // When these tokens were cached
}

// NewTokenManager creates a new token manager
func NewTokenManager(managementClient *api.ClientWithResponses, managementToken string) *TokenManager {
	return &TokenManager{
		managementClient: managementClient,
		managementToken:  managementToken,
		projectTokens:    make(map[string]*ProjectTokens),
	}
}

// GetServiceRoleToken returns the service role JWT token for a project
// This token can be used for storage API operations
func (tm *TokenManager) GetServiceRoleToken(ctx context.Context, projectRef string) (string, error) {
	// Check if we have a cached token that's still fresh (cache for 1 hour)
	tm.mutex.RLock()
	tokens, exists := tm.projectTokens[projectRef]
	if exists && time.Since(tokens.CachedAt) < time.Hour && tokens.ServiceRoleKey != "" {
		defer tm.mutex.RUnlock()
		return tokens.ServiceRoleKey, nil
	}
	tm.mutex.RUnlock()

	// Fetch fresh tokens from the management API
	tokens, err := tm.fetchProjectTokens(ctx, projectRef)
	if err != nil {
		return "", fmt.Errorf("failed to fetch project tokens: %w", err)
	}

	// Cache the tokens
	tm.mutex.Lock()
	tm.projectTokens[projectRef] = tokens
	tm.mutex.Unlock()

	return tokens.ServiceRoleKey, nil
}

// GetManagementToken returns the management API token
func (tm *TokenManager) GetManagementToken() string {
	return tm.managementToken
}

// fetchProjectTokens fetches API keys from the management API
func (tm *TokenManager) fetchProjectTokens(ctx context.Context, projectRef string) (*ProjectTokens, error) {
	resp, err := tm.managementClient.V1GetProjectApiKeysWithResponse(ctx, projectRef, &api.V1GetProjectApiKeysParams{
		Reveal: true, // Required to get actual API key values instead of masked ones
	})
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), string(resp.Body))
	}

	if resp.JSON200 == nil || len(*resp.JSON200) == 0 {
		return nil, fmt.Errorf("no API keys found for project %s", projectRef)
	}

	tokens := &ProjectTokens{
		CachedAt: time.Now(),
	}

	// Find the service role and anon keys
	for _, key := range *resp.JSON200 {
		if key.Name == "service_role" {
			tokens.ServiceRoleKey = key.ApiKey
		} else if key.Name == "anon" {
			tokens.AnonKey = key.ApiKey
		}
	}

	if tokens.ServiceRoleKey == "" {
		return nil, fmt.Errorf("service_role key not found for project %s", projectRef)
	}

	return tokens, nil
}

// InvalidateProjectTokens clears the cached tokens for a project
// This can be called if authentication fails to force a refresh
func (tm *TokenManager) InvalidateProjectTokens(projectRef string) {
	tm.mutex.Lock()
	delete(tm.projectTokens, projectRef)
	tm.mutex.Unlock()
}