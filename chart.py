#!/usr/bin/env python3

# Converts the output from the benchmark into a html document.
# Usage:
#  ./hashmap/benchmark.sh | tee /tmp/bench.out
#  ./hashmap/chart.py < /tmp/bench.out > /tmp/chart.html
#  firefox /tmp/chart.html

import argparse
import sys
from collections import defaultdict

def setup_arg_parser():
    """
    Set up argument parser and default values
    :return: parsed argument list
    """
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-f",
        "--file",
        type=str,
        nargs='?',
        default=None,
        help='raw benchmark file'
    )
    parser.add_argument(
        "-o",
        "--out",
        type=str,
        nargs='?',
        default=None,
        help='output html file'
    )
    return parser.parse_args()

def main():

    args = setup_arg_parser()

    fd_in = sys.stdin
    if args.file is not None:
        fd_in = open(args.file, "r")

    fd_out = sys.stdout
    if args.out is not None:
        fd_out = open(args.out, "w")

    #
    # collect metric values
    #
    mapping = defaultdict(lambda: defaultdict(list))
    for line in fd_in:
        lineRaw = line.split('\t')
        if len(lineRaw) != 8 or not lineRaw[0].startswith('Benchmark'):
            continue
        firstRaw = lineRaw[0].split('/')
        benchName = firstRaw[0]
        annotationList = firstRaw[1].split('-')
        mapName = annotationList[0]
        n = int(annotationList[1])
        time_ms = int(lineRaw[2].strip().split(' ')[0]) / (1000 * 1000)
        load = 'load=' + lineRaw[4].strip().split(' ')[0]
        mapping[benchName][mapName].append((n,time_ms,load))
        if "BenchmarkU32RandomFullInserts" in benchName:
            memory_bytes = int(lineRaw[3].strip().split(' ')[0]) / (1024 * 1024)
            mapping["MemoryConsumptionU32"][mapName].append((n,memory_bytes,load))
        if "BenchmarkU64RandomFullInserts" in benchName:
            memory_bytes = int(lineRaw[3].strip().split(' ')[0]) / (1024 * 1024)
            mapping["MemoryConsumptionU64"][mapName].append((n,memory_bytes,load))
        if "BenchmarkUUIDRandomInserts" in benchName:
            memory_bytes = int(lineRaw[3].strip().split(' ')[0]) / (1024 * 1024)
            mapping["MemoryConsumptionUUID"][mapName].append((n,memory_bytes,load))


    #
    # fd_out.write html document
    #

    # fd_out.write html header
    fd_out.write('''<!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <title>Hashmap Benchmark</title>
        <script src='https://cdn.plot.ly/plotly-2.20.0.min.js'></script>
    </head>
    <body>
        <h1 align="center">Hashmap Benchmark</h1>
        <hr>
    ''')
        
    for benchmark in sorted(mapping):
        b = mapping[benchmark]
        fd_out.write("<div id='"+benchmark+"'><script>\n")
        names = []
        y_naming = "time (ms)"
        if benchmark == "MemoryConsumption":
            y_naming = "memory (MB)"
        for mapName in b:
            points = b[mapName]
            x_values = map(lambda x: x[0], points)
            y_values = map(lambda x: x[1], points)
            load_values = map(lambda x: x[2], points)
            name = benchmark+'_'+mapName
            names.append(name)
            fd_out.write('var ' + name + ' = {\n')
            fd_out.write("name: '" + mapName + "',\n")
            fd_out.write('    x: ' + str(list(x_values)) + ',\n')
            fd_out.write('    y: ' + str(list(y_values)) + ',\n')
            fd_out.write('    text: ' + str(list(load_values)) + ',\n')
            fd_out.write('''   mode: 'lines+markers', type: 'scatter'
    };\n''')
        fd_out.write("var data_" + benchmark + "=" + '[%s]' % ', '.join(map(str, names)) + ";\n")
        fd_out.write("var layout_" + benchmark + " = {title:'" + benchmark + "', xaxis: {title: 'number of entries in hash table'},yaxis: {title: '" + y_naming + "'}};\n");
        fd_out.write("Plotly.newPlot('" + benchmark + "', data_"+ benchmark + ", layout_" + benchmark + ");\n"),
        fd_out.write("</script></div><hr>\n")

    # fd_out.write rest of body
    fd_out.write("</body>\n")


if __name__ == "__main__":
    main()