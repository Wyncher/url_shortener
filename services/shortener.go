package services

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"url_shortener/db"
	"url_shortener/models"
)

var conn = db.Connect()

func generateShrinkLink() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 15)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func Shrink(c *gin.Context) {
	var link models.Input_link
	if err := c.BindQuery(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	duration, err := time.ParseDuration("1h30m45s")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shinked := models.RedisStruct{Link: link.Link, ShrinkLink: generateShrinkLink()}
	val, err := models.MarshalBinary(shinked)
	countDB, err := conn.DBSize(ctx).Uint64()
	err = conn.Set(ctx, strconv.FormatUint(countDB, 10), val, duration).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var byteRedisStruct []byte
	err = conn.Get(ctx, strconv.FormatUint(countDB, 10)).Scan(&byteRedisStruct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var redisStruct models.RedisStruct
	redisStruct, err = models.UnmarshalBinary(byteRedisStruct, redisStruct)
	c.JSON(http.StatusBadRequest, gin.H{"link": redisStruct.Link, "short_link": redisStruct.ShrinkLink})
	return
}
