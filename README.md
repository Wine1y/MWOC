# Minecraft World Owner Changer (MWOC)

Simple console tool to change the owner of a Minecraft world (Player NBT tag in level.dat).
Useful when playing on LAN to transfer your world to another player and save NBT (Inventory, mods data, xp)

More info about the problem [here](https://gaming.stackexchange.com/questions/197686/changing-minecraft-playerdata-when-transferring-a-saved-game-open-to-lan)

**WARNING:** The code is messy and it wasn't tested at all, probably it won't work in some cases.

 ### **CREATE BACKUP OF THE WORLD BEFORE USING!**

## Usage
```
MWOC.exe <path_to_world_folder> <new_owner_nickname>
```
If no path is specified, MWOC will try to find a world in the working directory. If no nickname is specified, MWOC will ask for it.
