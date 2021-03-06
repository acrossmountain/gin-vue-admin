package service

import (
	"errors"
	"gin-vue-admin/model"

	"github.com/go-redis/redis"
	"github.com/go-spring/spring-boot"
	"github.com/jinzhu/gorm"
)

func init() {
	SpringBoot.RegisterBean(new(JwtBlackListService))
}

type JwtBlackListService struct {
	Db    *gorm.DB      `autowire:""`
	Redis redis.Cmdable `autowire:"?"`
}

// @title    JsonInBlacklist
// @description   create jwt blacklist
// @param     jwtList         model.JwtBlacklist
// @auth                     （2020/04/05  20:22）
// @return    err             error

func (service *JwtBlackListService) JsonInBlacklist(jwtList model.JwtBlacklist) (err error) {
	err = service.Db.Create(&jwtList).Error
	return
}

// @title    IsBlacklist
// @description   check if the Jwt is in the blacklist or not, 判断JWT是否在黑名单内部
// @auth                     （2020/04/05  20:22）
// @param     jwt             string
// @param     jwtList         model.JwtBlacklist
// @return    err             error

func (service *JwtBlackListService) IsBlacklist(jwt string, jwtList model.JwtBlacklist) bool {
	isNotFound := service.Db.Where("jwt = ?", jwt).First(&jwtList).RecordNotFound()
	return !isNotFound
}

// @title    GetRedisJWT
// @description   Get user info in redis
// @auth                     （2020/04/05  20:22）
// @param     userName        string
// @return    err             error
// @return    redisJWT        string

func (service *JwtBlackListService) GetRedisJWT(userName string) (err error, redisJWT string) {
	if service.Redis == nil {
		return errors.New("redis is nil"), ""
	}
	redisJWT, err = service.Redis.Get(userName).Result()
	return err, redisJWT
}

// @title    SetRedisJWT
// @description   set jwt into the Redis
// @auth                     （2020/04/05  20:22）
// @param     jwtList         model.JwtBlacklist
// @param     userName        string
// @return    err             error

func (service *JwtBlackListService) SetRedisJWT(jwtList model.JwtBlacklist, userName string) (err error) {
	if service.Redis == nil {
		return errors.New("redis is nil")
	}
	err = service.Redis.Set(userName, jwtList.Jwt, 1000*1000*1000*60*60*24*7).Err()
	return err
}
