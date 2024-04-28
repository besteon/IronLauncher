#!/usr/bin/env python3

import os
import sys
import time
import hashlib
import json
import subprocess
import networkx as nx

FILE_HASHES = { }

def patch_recursive_array(romFullPath, patchesFolder, orig_depth, depth, patches, patch_results = [ ]):
    if depth == 0: return None

    orig_sha1 = ''
    patch_sha1 = ''
    new_sha1 = ''

    rom = os.path.basename(romFullPath)
    with open(f'{romFullPath}', 'rb') as r:
        orig_sha1 = hashlib.sha1(r.read()).hexdigest()
    FILE_HASHES[orig_sha1] = rom
    
    for patch in os.listdir(patchesFolder):
        with open(f'{patchesFolder}/{patch}', 'rb') as p:
            patch_sha1 = hashlib.sha1(p.read()).hexdigest()
        FILE_HASHES[patch_sha1] = patch

        if patch_sha1 not in patches:
            if patch.endswith('.ips') or patch.endswith('.bps'):
                new_rom = f'buildscripts/{patch}-{rom}'
                flips_args = ["files/tools/binaries/flips", "--apply", f"{patchesFolder}/{patch}", f'{romFullPath}', f'{new_rom}']
                proc = subprocess.call(flips_args)
                with open(new_rom, 'rb') as new:
                    new_sha1 = hashlib.sha1(new.read()).hexdigest()

                root_rom = depth == orig_depth
                patch_results.append({
                    "root_rom": root_rom,
                    "orig": orig_sha1,
                    "patch": patch_sha1,
                    "new": new_sha1
                })

                p = patches[:]
                p.append(patch_sha1)
                patch_recursive_array(new_rom, patchesFolder, orig_depth, int(depth)-1, p, patch_results)

            elif patch.endswith('.xdelta'):
                continue

    return patch_results

if __name__ == '__main__':

    romsFolder = sys.argv[1]
    patchesFolder = sys.argv[2]
    depth = sys.argv[3]

    results = [ ]

    romsFolder = os.path.join(os.path.dirname(__file__), romsFolder)
    patchesFolder = os.path.join(os.path.dirname(__file__), patchesFolder)

    for rom in os.listdir(romsFolder):
        hash = ''
        with open(f'{romsFolder}/{rom}', 'rb') as r:
            hash = hashlib.sha1(r.read()).hexdigest()
        results += patch_recursive_array(f'{romsFolder}/{rom}', patchesFolder, depth, depth, [ ])

    for rom in os.listdir(os.getcwd() + '/buildscripts'):
        if rom.endswith('.gba'):
            os.remove('buildscripts/'+rom)

    patchgraphpath = os.path.join(os.path.dirname(__file__), '../files/tools/data/patchgraph.json')
    with open(patchgraphpath, 'w') as f:
        json.dump(results, f, indent=2)