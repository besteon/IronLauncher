#!/usr/bin/env python3

import os
import sys
import time
import hashlib
import json
import subprocess
from collections import defaultdict

def patch_recursive(romFullPath, patchesFolder, depth, patches):
    if depth == 0: return None

    orig_sha1 = ''
    patch_sha1 = ''
    new_sha1 = ''
    decendents = None
    rom = os.path.basename(romFullPath)
    with open(f'{romFullPath}', 'rb') as r:
        orig_sha1 = hashlib.sha1(r.read()).hexdigest()
    
    patch_results = { }
    for patch in os.listdir(patchesFolder):
        with open(f'{patchesFolder}/{patch}', 'rb') as p:
            patch_sha1 = hashlib.sha1(p.read()).hexdigest()

        if patch_sha1 not in patches:
            if patch.endswith('.ips') or patch.endswith('.bps'):
                new_rom = f'{patch}-{rom}'
                flips_args = ["flips", "--apply", f"{patchesFolder}/{patch}", f'{romFullPath}', f'{new_rom}']
                proc = subprocess.call(flips_args)
                with open(new_rom, 'rb') as new:
                    new_sha1 = hashlib.sha1(new.read()).hexdigest()

                p = patches[:]
                p.append(patch_sha1)
                res = patch_recursive(new_rom, patchesFolder, int(depth)-1, p)
                if res == None:
                    patch_results[patch_sha1] = {
                        'rom_hash': new_sha1,
                    }
                else:
                    patch_results[patch_sha1] = {
                        'rom_hash': new_sha1,
                        'patches': res
                    }

            elif patch.endswith('.xdelta'):
                continue

    return patch_results

def get_all_keys(d):
    for key, value in d.items():
        yield key
        if isinstance(value, dict):
            yield from get_all_keys(value)

def get_all_of_key(d, key):
    for k, v in d.items():
        if k == key:
            yield v
        if isinstance(v, dict):
            yield from get_all_of_key(v, key)

def inverse_dict(d, seen = {}):
  for k, v in d.items():
    if not isinstance(v, dict):
      yield {v: seen}
    else:
      seen[k] = v
      yield from inverse_dict(v, seen)

def trim_dict(d, seen = {}):
  for k, v in d.items():
    if not isinstance(v, dict):
      yield {k: seen}
    else:
      seen[k] = v
      yield from trim_dict(v, seen)

if __name__ == '__main__':

    romsFolder = sys.argv[1]
    patchesFolder = sys.argv[2]
    depth = sys.argv[3]

    results = { }

    for rom in os.listdir(romsFolder):
        hash = ''
        with open(f'{romsFolder}/{rom}', 'rb') as r:
            hash = hashlib.sha1(r.read()).hexdigest()
        results[rom] = {
            'rom_hash': hash,
            'patches': patch_recursive(f'{romsFolder}/{rom}', patchesFolder, depth, [ ])
        }

    target = ''

    inverse = list(inverse_dict(results))
    unique_roms = list(set(get_all_of_key(results, 'rom_hash')))

    r = []
    for x in unique_roms:
        for y in inverse:
            if x == next(iter(y)):
                r.append(y)
                break

    trimmed = list(trim_dict(inverse[0]))
    s = json.dumps(inverse, indent=2)
    print(s)
    print(len(s))

    print(unique_roms)
