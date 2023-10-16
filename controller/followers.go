package controller

import (
	"compress/gzip"
	"context"
	"diary_api/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

type request struct {
	First  uint `json:"first"`
	Second uint `json:"second"`
}

func FindFollowers(c *gin.Context) {
	userIdParam := c.Query("userid")

	if userIdParam == "" {
		fmt.Println("Error: No query provided")
		c.String(http.StatusBadRequest, "No query provided")
		return
	}

	userId, err := strconv.ParseUint(userIdParam, 10, 0)
	if err != nil {
		c.String(http.StatusBadRequest, "userid param is not a number")
		return
	}

	user, err := model.FindUserById(uint(userId))

	if err != nil {
		fmt.Println("Error in FindUserById for user:", err)
		c.String(http.StatusInternalServerError, "Error finding user(s)")
		return
	}

	var userFollower []model.Follower
	findFollowersRequest(user.InstaId, &userFollower)
	c.JSON(http.StatusBadRequest, gin.H{"data": userFollower})
}

func FindAndStoreFollowers(c *gin.Context) {
	userIdParam := c.Query("userid")

	if userIdParam == "" {
		fmt.Println("Error: No query provided")
		c.String(http.StatusBadRequest, "No query provided")
		return
	}

	userId, err := strconv.ParseUint(userIdParam, 10, 0)
	if err != nil {
		c.String(http.StatusBadRequest, "userid param is not a number")
		return
	}

	user, err := model.FindUserById(uint(userId))

	if err != nil {
		fmt.Println("Error in FindUserById for user:", err)
		c.String(http.StatusInternalServerError, "Error finding user(s)")
		return
	}

	var userFollower []model.Follower

	findFollowersRequest(user.InstaId, &userFollower)

	if err := user.SaveFollowers(userFollower); err != nil {
		// Handle the error
		fmt.Println("Error saving followers for user:", err)
		c.String(http.StatusInternalServerError, "Error storing followers")
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": userFollower})
}

func FindCommon(c *gin.Context) {
	var profiles request
	if err := c.ShouldBindJSON(&profiles); err != nil {
		fmt.Println("Error:", err)
		c.String(http.StatusBadRequest, "Invalid JSON provided")
		return
	}

	user1, err1 := model.FindUserById(profiles.First)
	user2, err2 := model.FindUserById(profiles.Second)

	if err1 != nil || err2 != nil {
		// Handle any error that occurred during the FindUserById calls
		if err1 != nil {
			fmt.Println("Error in FindUserById for user1:", err1)
		}
		if err2 != nil {
			fmt.Println("Error in FindUserById for user2:", err2)
		}
		// Handle the error accordingly, e.g., return an error response
		c.String(http.StatusInternalServerError, "Error finding user(s)")
		return
	}

	followers1 := user1.Followers
	followers2 := user2.Followers

	user1Map := make(map[string]bool)
	var commonUsers []model.Follower

	for _, user := range followers1 {
		user1Map[user.UserName] = true
	}

	for _, user := range followers2 {
		if _, exists := user1Map[user.UserName]; exists {
			commonUsers = append(commonUsers, user)
			delete(user1Map, user.UserName)
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"data": commonUsers})

}

func prepareRequest() *http.Request {
	req, err := http.NewRequest("GET", "https://www.instagram.com/api/v1/friendships/2217206142/following/?count=100", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set custom headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", `dpr=2; mid=ZSwcCQAEAAHtVCOM8X0o3bxVcpjU; ig_did=D1F6950C-FC23-4C23-8DCE-99E412DCC16B; ig_nrcb=1; datr=CBwsZSR847kWJPYz3GkXBBeu; csrftoken=Nrqq50mNbzSgESsYWUYmqyzKHiLcRjIC; ds_user_id=708148013; sessionid=708148013%3Aq2Cc4ldXgT38pp%3A5%3AAYdWtWuCJwV726bXqThjW6_feHlrkPTtRWcG1K0IoQ; shbid="8953\054708148013\0541728925592:01f74396b41fd8a708aca8de5d058925a5cdde5fd02f21f0f2c2a2aff350368450c7c0b2"; shbts="1697389592\054708148013\0541728925592:01f7e3c3fb7afec8e2ca213c94cad3b3064aa083a2a96a026f847f204dae57ced15af3fd"; rur="ODN\054708148013\0541728996351:01f79a63e3c470e4b7bff2b034260c8f4c3f6afaa9f29dc1bdae1a3b60a08081495d795f"`) // shortened for brevity
	req.Header.Set("Dpr", "2")
	req.Header.Set("Referer", "https://www.instagram.com/gornersisyan_/following/")
	req.Header.Set("Sec-Ch-Prefers-Color-Scheme", "light")
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"`)
	req.Header.Set("Sec-Ch-Ua-Full-Version-List", `"Google Chrome";v="117.0.5938.132", "Not;A=Brand";v="8.0.0.0", "Chromium";v="117.0.5938.132"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Model", `""`)
	req.Header.Set("Sec-Ch-Ua-Platform", "macOS")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "13.3.1")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Set("Viewport-Width", "584")
	req.Header.Set("X-Asbd-Id", "129477")
	req.Header.Set("X-Csrftoken", "hnIJEeRPY43U32mFVTBNpnJpg3vaOp1Q")
	req.Header.Set("X-Ig-App-Id", "936619743392459")
	req.Header.Set("X-Ig-Www-Claim", "hmac.AR2Oam_oTOK2CyMnfOg7e6_0hCg_X1vcHd_eQUXZzrx1TX01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	return req
}

func findFollowersRequest(profileId string, userPointer *[]model.Follower) {
	maxId := 0
	var err error
	var users []model.Follower

	for maxId != -1 {
		baseRequest := prepareRequest()
		newReq := baseRequest.Clone(context.TODO()) // This creates a shallow copy of the request
		newURL := fmt.Sprintf("https://www.instagram.com/api/v1/friendships/%s/following/?count=100&max_id=%d", profileId, maxId)
		fmt.Println("URL requested: ", newURL)
		newReq.URL, err = url.Parse(newURL)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new HTTP client and execute the request
		client := &http.Client{}
		resp, err := client.Do(newReq)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		var reader io.Reader = resp.Body

		// Check if the response is gzip-encoded.
		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err = gzip.NewReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer reader.(*gzip.Reader).Close()
		}

		type Response struct {
			Users []model.Follower `json:"users"`
			Next  string           `json:"next_max_id"`
		}

		var result Response
		if err := json.NewDecoder(reader).Decode(&result); err != nil {
			log.Fatal(err)
		}

		if result.Next != "" {
			maxId, err = strconv.Atoi(result.Next)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			maxId = -1
		}
		// Print the `users` field.
		users = append(users, result.Users...)
	}
	*userPointer = users
}
