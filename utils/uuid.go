package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const mojangApiUrl string = "https://api.mojang.com/users/profiles/minecraft"

type uUIDResponse struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func addStripesToUUID(uuid *string) {
	*uuid = fmt.Sprintf("%s-%s-%s-%s-%s", (*uuid)[:8], (*uuid)[8:12], (*uuid)[12:16], (*uuid)[16:20], (*uuid)[20:])
}

func OfflineUsernameToUUID(username string) string {
	username = fmt.Sprintf("OfflinePlayer:%s", username)
	hash := md5.Sum([]byte(username))
	hash[6] = hash[6]&0x0f | 0x30
	hash[8] = hash[8]&0x3f | 0x80
	uuid := hex.EncodeToString(hash[:])
	addStripesToUUID(&uuid)
	return uuid
}

func UsernameToUUID(username string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", mojangApiUrl, username))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("mojang api returned %s", resp.Status)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var uuidRespones *uUIDResponse = &uUIDResponse{}
	err = json.Unmarshal(content, uuidRespones)
	if err != nil {
		return "", err
	}
	addStripesToUUID(&uuidRespones.ID)
	return uuidRespones.ID, nil
}

func TryMapUsernameToPlayerdataMeta(
	username string,
	metas []PlayerDataMeta,
) (
	offlineMeta *PlayerDataMeta,
	onlineMeta *PlayerDataMeta,
) {
	offlineUUID := OfflineUsernameToUUID(username)
	onlineUUID, onlineErr := UsernameToUUID(username)
	for _, meta := range metas {
		if meta.UUID == offlineUUID {
			newMeta := meta
			offlineMeta = &newMeta
		}
		if onlineErr == nil && meta.UUID == onlineUUID {
			newMeta := meta
			onlineMeta = &newMeta
		}
	}
	return offlineMeta, onlineMeta
}
