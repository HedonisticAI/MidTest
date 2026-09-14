package cache_repo

import (
	"midtest/internal/auth"
	"midtest/internal/infrastructure/cache"
	"time"
)

const CacheDefaultFileExpire = 10 * time.Minute

type RepoCache struct {
	Cache *cache.Cache
}

func (RepoCache *RepoCache) LoadFile(Name string, Data interface{}, Duration time.Duration) {
	if Duration == 0 {
		Duration = CacheDefaultFileExpire
	}
	RepoCache.Cache.Set(Name, Data, Duration)
}

func (RepoCache *RepoCache) LoadAuth(Info auth.LoginInfo) {
	RepoCache.Cache.Set(Info.Login, Info.Password, time.Hour)
}

func (RepoCache *RepoCache) LoadToken(AuthInfo auth.AuthInfo) {
	RepoCache.Cache.Set(AuthInfo.Token, AuthInfo.ID, 5*time.Minute)
}

func (RepoCache *RepoCache) IsActive(Token string) bool {
	_, res := RepoCache.Cache.Get(Token)
	return res
}

func (RepoCache *RepoCache) CheckPWD(Info auth.LoginInfo) bool {
	Get, exists := RepoCache.Cache.Get(Info.Login)
	if !exists {
		return false
	}
	if Get != Info.Password {
		return false
	}

	return true
}

func (RepoCache *RepoCache) GetAuth(Token string) string {
	res, exists := RepoCache.Cache.Get(Token)
	if !exists {
		return ""
	}
	return res.(string)
}

func (RepoCache *RepoCache) DeleteItem(key string) error {
	return RepoCache.Cache.Delete(key)
}

func (RepoCache *RepoCache) GetFile(ID string) (interface{}, bool) {
	return RepoCache.Cache.Get(ID)
}
