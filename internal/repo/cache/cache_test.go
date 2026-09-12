package cache_repo_test

import (
	"midtest/internal/auth"
	"midtest/internal/infrastructure/cache"
	cache_repo "midtest/internal/repo/cache"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {
	// Create a new cache instance
	cacheInstance := cache.NewCache(5*time.Minute, 10*time.Minute)

	// Create a new RepoCache instance
	repoCache := &cache_repo.RepoCache{
		Cache: cacheInstance,
	}

	// Test LoadFile method
	testData := "test data"
	repoCache.LoadFile("testKey", testData, 0)
	repoCache.LoadAuth(auth.LoginInfo{Login: "testLogin", Password: "testPassword"})
	// Test IsActive method
	if !repoCache.IsActive("testKey") {
		t.Errorf("Expected key to be active")
	}

	// Test GetAuth method
	authInfo := repoCache.GetAuth("testKey")
	if authInfo != testData {
		t.Errorf("Expected auth info to be '%s', got '%s'", testData, authInfo)
	}
}
