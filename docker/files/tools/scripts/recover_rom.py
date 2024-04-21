#!/usr/bin/env python3

import os
import sys
import time
import hashlib
import json
import subprocess
import networkx as nx

FILE_HASHES = { }



if __name__ == '__main__':

    romsFolder = sys.argv[1]
    patchesFolder = sys.argv[2]
    patchGraphFile = sys.argv[3]
    target_hash = sys.argv[4]
    out_file = sys.argv[5]

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

    path = None
    shortest = 1000
    for x in rom_nodes:
        sp = nx.shortest_path(G, source=x, target=target_hash)
        if len(sp) < shortest:
            shortest = len(sp)
            path = sp

    pathGraph = nx.path_graph(sp)

    base_rom = FILE_HASHES[list(pathGraph.edges())[0][0]]
    from_rom = f'{romsFolder}/{base_rom}'
    to_rom = f'/tmp/{base_rom}'
    for edge in pathGraph.edges():
        patch_file = FILE_HASHES[G.edges[edge[0], edge[1]]['patch']]
        flips_args = ["flips", "--apply", f"{patchesFolder}/{patch_file}", f'{from_rom}', f'{to_rom}']
        from_rom = to_rom

        proc = subprocess.call(flips_args)
        with open(to_rom, 'rb') as new:
            new_sha1 = hashlib.sha1(new.read()).hexdigest()
            if new_sha1 == target_hash:
                os.rename(to_rom, out_file)
                print('Successfully recovered target ROM.')
            else:
                print('There was an issue recovering the target rom')
