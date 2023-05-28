package utils

import (
	"os"
	"path"
	"regexp"
)

type PlayerDataMeta struct {
	UUID string
	Path string
}

type WorldDataPaths struct {
	BasePath        string
	PlayerdataMetas []PlayerDataMeta
	LevelPath       string
}

var uuidRegex *regexp.Regexp = regexp.MustCompile(`[\da-z]{8}-[\da-z]{4}-[\da-z]{4}-[\da-z]{4}-[\da-z]{12}`)

func TryFindWorldDataPathsAt(searching_path string) (paths *WorldDataPaths, exists bool) {
	levelPath := path.Join(searching_path, "level.dat")
	if _, err := os.Stat(levelPath); err != nil {
		return nil, false
	}

	playerDataPath := path.Join(searching_path, "playerdata")
	if _, err := os.Stat(playerDataPath); err != nil {
		return nil, false
	}

	playerdata_entries, err := os.ReadDir(playerDataPath)
	if err != nil {
		return nil, false
	}
	var playerDataMetas = make([]PlayerDataMeta, 0, 1)
	for _, entry := range playerdata_entries {
		if uuid := uuidRegex.FindString(entry.Name()); len(uuid) > 0 {
			playerDataMetas = append(playerDataMetas, PlayerDataMeta{
				UUID: uuid,
				Path: path.Join(playerDataPath, entry.Name()),
			})
		}

	}

	return &WorldDataPaths{searching_path, playerDataMetas, levelPath}, true
}

func TryFindLocalWorldDataPaths() (paths *WorldDataPaths, exists bool) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, false
	}
	paths, exists = TryFindWorldDataPathsAt(workingDir)
	if exists {
		return paths, exists
	}
	wdEntries, err := os.ReadDir(workingDir)
	if err != nil {
		return nil, false
	}
	for _, entry := range wdEntries {
		if entry.IsDir() {
			paths, exists = TryFindWorldDataPathsAt(path.Join(workingDir, entry.Name()))
			if exists {
				return paths, exists
			}
		}
	}
	return nil, false
}
