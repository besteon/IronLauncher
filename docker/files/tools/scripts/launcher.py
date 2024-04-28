import configparser
import os
import shutil
import hashlib
import json
import subprocess
import networkx as nx


def recoverRom(romsFolder, patchesFolder, patchGraphFile, target_hash, base_rom, out_folder):

    FILE_HASHES = {}

    for rom in os.listdir(romsFolder):
        with open(f'{romsFolder}/{rom}', 'rb') as r:
            orig_sha1 = hashlib.sha1(r.read()).hexdigest()
        FILE_HASHES[orig_sha1] = rom
    
    for patch in os.listdir(patchesFolder):
        with open(f'{patchesFolder}/{patch}', 'rb') as p:
            patch_sha1 = hashlib.sha1(p.read()).hexdigest()
        FILE_HASHES[patch_sha1] = patch

    results = [ ]

    with open(patchGraphFile, 'r') as j:
        results = json.load(j)

    G = nx.DiGraph()

    for x in results:
        G.add_node(x['orig'], root_rom=x['root_rom'])
        G.add_node(x['new'], root_rom=False)
        G.add_edge(x['orig'], x['new'], patch=x['patch'])

    rom_nodes = [x for x,y in G.nodes(data=True) if y['root_rom'] == True]

    shortest = 1000
    path = None
    for x in rom_nodes:
        sp = nx.shortest_path(G, source=x, target=target_hash)
        if len(sp) < shortest:
            shortest = len(sp)
            path = sp

    pathGraph = nx.path_graph(path)

    edges = list(pathGraph.edges())
    if len(edges) > 0:
        from_rom = f'{romsFolder}/{base_rom}'
        to_rom = f'/tmp/{base_rom}'
        for edge in pathGraph.edges():
            patch_file = FILE_HASHES[G.edges[edge[0], edge[1]]['patch']]
            cmd = f'/home/launcher/files/tools/binaries/flips --apply "{patchesFolder}/{patch_file}" "{from_rom}" "{to_rom}"'
            from_rom = to_rom

            proc = subprocess.call(cmd, shell=True)
            with open(to_rom, 'rb') as new:
                new_sha1 = hashlib.sha1(new.read()).hexdigest()
                if new_sha1 == target_hash:
                    #shutil.copy(to_rom, f"{out_folder}/{base_rom}")
                    print('Successfully recovered target ROM.')
                    return True
                else:
                    print(f'Applied patch {patch_file}. More patches to go.')                

        print("There was an issue recovering the target ROM")
        return False
    else:
        print('No patches needed')
        shutil.copy(f'{romsFolder}/{base_rom}', f'/tmp/{base_rom}')
        return True


if __name__ == '__main__':

    HOME = os.environ['HOME']

    # If config files don't exist, copy the default ones
    BIZHAWKINI = f"{HOME}/BizHawk/config.ini"
    GBASETTINGSINI = f"{HOME}/BizHawk/Lua/gba/Ironmon-Tracker/Settings.ini"
    NDSSETTINGSINI = f"{HOME}/BizHawk/Lua/nds/Ironmon-Tracker/Settings.ini"
    if os.path.getsize(GBASETTINGSINI) == 0:
        shutil.copyfile(f"{HOME}/GbaSettings.ini", GBASETTINGSINI)
    if os.path.getsize(NDSSETTINGSINI) == 0:
        shutil.copyfile(f"{HOME}/NdsSettings.ini", NDSSETTINGSINI)
    if os.path.getsize(BIZHAWKINI) == 0:
        shutil.copyfile(f"{HOME}/Bizhawk.ini", BIZHAWKINI)

    bizhawk_config = { }
    with open(BIZHAWKINI, 'r') as b_config:
        bizhawk_config = json.load(b_config)

    # Read ironlauncher.ini
    config = configparser.ConfigParser()
    config.read(f'{HOME}/ironlauncher.ini')

    rom = config['settings']['defaultRom']
    mode = config['settings']['defaultMode']
    qolPatches = config['settings']['qolPatches']

    PATCH_GRAPH = os.path.join(os.path.dirname(__file__), f"../data/patchgraph.json")
    GAME_MODES = os.path.join(os.path.dirname(__file__), f"../data/gamemodes.json")

    game_modes = {}
    with open(GAME_MODES, 'r') as g_modes:
        game_modes = json.load(g_modes)

    selected_rom_hash = ""
    with open(f"/roms/{rom}", 'rb') as new:
        selected_rom_hash = hashlib.sha1(new.read()).hexdigest()

    modeObj = game_modes[selected_rom_hash]["Modes"][mode]
    target_hash = modeObj["QolHash"] if qolPatches and "QolHash" in modeObj else modeObj['RomHash']

    print(f'Target: {target_hash}')
    out_folder = "/tmp"
    recoveredRom = recoverRom("/roms", f"{HOME}/files/patches", PATCH_GRAPH, target_hash, rom, out_folder)

    if modeObj['Tracker']:
        ironmonSettings = configparser.ConfigParser()
        ironmonSettings.optionxform = str
        system = game_modes[selected_rom_hash]['System'].lower()
        settingsFile = f'/home/launcher/BizHawk/Lua/{system}/Ironmon-Tracker/Settings.ini'
        ironmonSettings.read(settingsFile)

        bizhawk_config['RecentLua']['recentlist'] = [f"/home/launcher/BizHawk/Lua/{system}/Ironmon-Tracker/Ironmon-Tracker.lua"]
        bizhawk_config['RecentLua']['AutoLoad'] = True

        rando = modeObj['RandomizerSettings']
        ironmonSettings['config']['Settings_File'] = f"/home/launcher/BizHawk/Lua/{system}/Ironmon-Tracker/ironmon_tracker/RandomizerSettings/{rando}"
        ironmonSettings['config']['Source_ROM'] = f"{out_folder}/{rom}"

        if modeObj['Randomizer'] == 'default':
            ironmonSettings['config']['Randomizer_JAR'] = '/home/launcher/PokeRandoZX.jar'
        elif modeObj['Randomizer'] == 'natdex':
            ironmonSettings['config']['Randomizer_JAR'] = f'/home/launcher/BizHawk/Lua/{system}/Ironmon-Tracker/extensions/natdex/randomizer-1.1.2.jar'
            ironmonSettings['config']['Settings_File'] = f"/home/launcher/BizHawk/Lua/{system}/Ironmon-Tracker/extensions/natdex/rnqs_files/{rando}"

        with open(settingsFile, 'w') as ironmonConfigFile:
            ironmonSettings.write(ironmonConfigFile, space_around_delimiters=False)
    else:
        bizhawk_config['RecentLua']['recentlist'] = []

    bizhawk_config['RecentRoms']['recentlist'] = [f"*OpenRom*{out_folder}/{rom}"]
    bizhawk_config['RecentRoms']['AutoLoad'] = True
    with open(BIZHAWKINI, 'w') as b_config:
        json.dump(bizhawk_config, b_config, indent=1)

        