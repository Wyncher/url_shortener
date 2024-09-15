package services

import (
	"errors"
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
func getRange(start uint, end uint) (result []uint) {
	result = make([]uint, 0, end-start)
	for value := start; value <= end; value++ {
		result = append(result, value)
	}

	return result
}
func searchLink(new_link string) (string, error) {
	ctx := context.Background()
	var byteRedisStruct []byte
	countDB, err := conn.DBSize(ctx).Uint64()
	if err != nil {
		return "", errors.New("Error DBsize")
	}
	for i := range getRange(0, uint(countDB)) {
		err := conn.Get(ctx, strconv.Itoa(i)).Scan(&byteRedisStruct)
		if err != nil {
			return "", errors.New("Error get DB")
		}
		var redisStruct models.RedisStruct
		redisStruct, err = models.UnmarshalBinary(byteRedisStruct, redisStruct)
		if redisStruct.Link == new_link {
			return redisStruct.ShrinkLink, nil
		}
	}
	return "", nil
}
func Shrink(c *gin.Context) {
	var link models.Input_link
	if err := c.BindQuery(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shrink_link, err := searchLink(link.Link)
	if shrink_link != "" && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"link": link.Link, "short_link": shrink_link})
		return
	}
	if err != nil {
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
	c.JSON(http.StatusBadRequest, gin.H{"link": shinked.Link, "short_link": shinked.ShrinkLink})
	return
}
func Redirect(c *gin.Context) {
	ctx := context.Background()
	var byteRedisStruct []byte
	var countDB uint64 = 1
	err := conn.Get(ctx, strconv.FormatUint(countDB, 10)).Scan(&byteRedisStruct)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var redisStruct models.RedisStruct
	redisStruct, err = models.UnmarshalBinary(byteRedisStruct, redisStruct)
	c.JSON(http.StatusBadRequest, gin.H{"link": redisStruct.Link, "short_link": redisStruct.ShrinkLink})
	return
}
