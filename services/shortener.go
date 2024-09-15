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
func searchLink(link string, shrink bool) (string, error) {
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
		if !shrink {
			if redisStruct.Link == link {
				return redisStruct.ShrinkLink, nil
			}
		} else {
			if redisStruct.ShrinkLink == link {
				return redisStruct.Link, nil
			}
		}

	}
	return "", nil
}
func Shrink(c *gin.Context) {
	var link models.Input_link
	if err := c.BindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shrink_link, err := searchLink(link.Link, false)
	if shrink_link != "" && err == nil {
		c.JSON(http.StatusOK, gin.H{"link": link.Link, "short_link": shrink_link})
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
	c.JSON(http.StatusOK, gin.H{"link": shinked.Link, "short_link": shinked.ShrinkLink})
	return
}
func Redirect(c *gin.Context) {
	var input_link models.Input_link
	if err := c.BindJSON(&input_link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link, err := searchLink(input_link.Link, true)
	if link != "" && err == nil {
		c.JSON(http.StatusOK, gin.H{"short_link": link})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "try use shrink before redirect"})
	return

}
