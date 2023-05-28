package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Wine1y/MWOC/utils"
)

var argsWorldPath *string
var argsNewOwnerName *string

func init() {
	args := os.Args[1:]
	if len(args) > 0 {
		argsWorldPath = &args[0]
	}
	if len(args) > 1 {
		argsNewOwnerName = &args[1]
	}
}

func readStringInput(input_message string) string {
	if len(input_message) > 0 {
		print(input_message)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func exitGracefully(exit_text string, exit_code int) {
	println(exit_text)
	print("Press any key to exit...")
	readStringInput("")
	os.Exit(exit_code)
}

func main() {
	println("Searching for valid Minecraft world...")
	var worldDataPaths utils.WorldDataPaths = resolveWorldDataPaths()
	fmt.Printf("World found at %s\n", worldDataPaths.BasePath)

	levelNBT := resolveLevelNBT(worldDataPaths)
	newWorldOwnerMeta := resolveNewWorldOwnerMeta(worldDataPaths)
	newWorldOwnerData := resolveNewWorldOwnerData(newWorldOwnerMeta)
	applyNewOwnerDataToLevelNBT(newWorldOwnerData, levelNBT)
	resolveSavingLevelNBT(levelNBT, worldDataPaths)
}

func resolveWorldDataPaths() utils.WorldDataPaths {
	if argsWorldPath != nil {
		worldDataPaths, exists := utils.TryFindWorldDataPathsAt(*argsWorldPath)
		if !exists {
			exitGracefully(fmt.Sprintf("Can't find Minecraft world at %s\n", *argsWorldPath), 1)
		}
		return *worldDataPaths
	} else {
		worldDataPaths, exists := utils.TryFindLocalWorldDataPaths()
		if !exists {
			exitGracefully("Can't find Minecraft world at current directory", 1)
		}
		return *worldDataPaths
	}
}

func resolveLevelNBT(worldDataPaths utils.WorldDataPaths) map[string]interface{} {
	levelNBT, err := utils.TryExtractDatNBT(worldDataPaths.LevelPath)
	if err != nil {
		exitGracefully(fmt.Sprintf("Can't read level.dat at %s", worldDataPaths.LevelPath), 1)
	}
	_, levelDataExists := levelNBT["Data"]
	if !levelDataExists {
		exitGracefully(fmt.Sprintf("Invalid level.dat at %s", worldDataPaths.LevelPath), 1)
	}
	return levelNBT
}

func resolveNewWorldOwnerMeta(worldDataPaths utils.WorldDataPaths) *utils.PlayerDataMeta {
	if argsNewOwnerName == nil {
		name := readStringInput("New world owner nickname: ")
		argsNewOwnerName = &name
	}

	offlineMeta, onlineMeta := utils.TryMapUsernameToPlayerdataMeta(
		*argsNewOwnerName,
		worldDataPaths.PlayerdataMetas,
	)
	var newOwnerMeta *utils.PlayerDataMeta

	switch {
	case offlineMeta != nil && onlineMeta == nil:
		newOwnerMeta = offlineMeta
		fmt.Printf("Using offline playerdata of %s (%s)\n", *argsNewOwnerName, newOwnerMeta.UUID)
	case onlineMeta != nil && offlineMeta == nil:
		newOwnerMeta = onlineMeta
		fmt.Printf("Using online playerdata of %s (%s)\n", *argsNewOwnerName, newOwnerMeta.UUID)
	case onlineMeta != nil && offlineMeta != nil:
	WaitingForTwoMetasResolve:
		for {
			user_resp := readStringInput(fmt.Sprintf(`Found 2 playerdatas for %s, one online and one offline, which one to use ? (Enter 1 or 2):`, *argsNewOwnerName))
			switch {
			case user_resp == "1":
				newOwnerMeta = offlineMeta
				fmt.Printf("Using offline playerdata of %s (%s)\n", *argsNewOwnerName, newOwnerMeta.UUID)
				break WaitingForTwoMetasResolve
			case user_resp == "2":
				newOwnerMeta = onlineMeta
				fmt.Printf("Using online playerdata of %s (%s)\n", *argsNewOwnerName, newOwnerMeta.UUID)
				break WaitingForTwoMetasResolve
			default:
				println("Invalid response, try again (Enter 1 to use offline playerdata and 2 to use online playerdata)")
			}
		}
	default:
	WaitingForEmptyMetaResolve:
		for {
			user_resp := readStringInput(fmt.Sprintf("Can't find playerdata for %s, use empty ? (First player logged to the world will become new owner) Yes/No: ", *argsNewOwnerName))
			switch {
			case strings.ToLower(user_resp) == "yes":
				newOwnerMeta = nil
				println("Using empty playerdata")
				break WaitingForEmptyMetaResolve
			case strings.ToLower(user_resp) == "no":
				os.Exit(0)
			default:
				println("Invalid response, try again (Enter 1 to use offline playerdata and 2 to use online playerdata)")
			}
		}
	}
	return newOwnerMeta
}

func resolveNewWorldOwnerData(meta *utils.PlayerDataMeta) map[string]interface{} {
	if meta == nil {
		return nil
	}
	playerdataNBT, err := utils.TryExtractDatNBT(meta.Path)
	if err != nil {
		exitGracefully(fmt.Sprintf("Can't read playerdata at %s", meta.Path), 1)
	}
	return playerdataNBT
}

func applyNewOwnerDataToLevelNBT(newWorldOwnerData map[string]interface{}, levelNBT map[string]interface{}) {
	if newWorldOwnerData != nil {
		levelNBT["Data"].(map[string]interface{})["Player"] = newWorldOwnerData
	} else {
		delete(levelNBT["Data"].(map[string]interface{}), "Player")
	}
}

func resolveSavingLevelNBT(levelNBT map[string]interface{}, worldDataPaths utils.WorldDataPaths) {
	saved := utils.TrySaveNBTToDat(levelNBT, worldDataPaths.LevelPath)
	if saved {
		exitGracefully(fmt.Sprintf("Successfully change world owner to %s", *argsNewOwnerName), 0)
	}
	exitGracefully(fmt.Sprintf("Can't save new level.dat to %s", worldDataPaths.LevelPath), 1)
}
