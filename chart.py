#!/usr/bin/env python3

# Converts the output from the benchmark into a html document.
# Usage:
#  ./hashmap/benchmark.sh | tee /tmp/bench.out
#  ./hashmap/chart.py < /tmp/bench.out > /tmp/chart.html
#  firefox /tmp/chart.html

import argparse
import sys
from collections import defaultdict

info = {
    "RandomShuffleInserts":'''
    Before the test, a vector with the values [0, n) is generated and shuffled.
    Then for each value k in the vector, the key-value pair (k, 1) is inserted into the hash map.
    ''',
    "RandomFullInserts":'''
    Before the test, a vector with n random values in the whole range of the integer size is generated.
    Then for each value k in the vector, the key-value pair (k, 1) is inserted into the hash map.
    ''',
    "RandomInserts":'''
    Before the test, a vector with n random UUIDs is generated.
    Then for each value k in the vector, the key-value pair (k, 1) is inserted into the hash map.
    ''',
    "InsertsWithReserve":'''
    Same as the UUID random inserts test but the reserve method of the hash map is called beforehand
    to avoid any rehash during the insertion. It provides a fair comparison even if the growth factor
    of each hash map is different.
    ''',
    "RandomFullWithReserveInserts":'''
    Same as the random full inserts test but the reserve method of the hash map is called beforehand
    to avoid any rehash during the insertion. It provides a fair comparison even if the growth factor
    of each hash map is different.
    ''',
    "RandomFullDeletes":'''
    Before the test, n elements in the same way as in the random full insert test are added.
    Each key is deleted one by one in a different and random order than the one they were inserted.
    ''',
    "RandomShuffleReads":'''
    Before the test, n elements are inserted in the same way as in the random shuffle inserts test.
    Each key-value pair is look up in a different and random order than the one they were inserted.
    ''',
    "FullReads":'''
    Before the test, n elements are inserted in the same way as in the random full inserts test.
    Each key-value pair is look up in a different and random order than the one they were inserted.
    ''',
    "RandomReads":'''
    Before the test, n elements are inserted in the same way as in the random UUID insert test.
    Read each key-value pair is look up in a different and random order than the one they were inserted.
    ''',
    "ReadsMisses":'''
     Before the test, n elements are inserted in the same way as in the random UUID insert test.
    Then a another vector of n random elements different from the inserted elements is generated
    which is tried to search in the hash map.
    ''',
    "FullReadsMisses":'''
    Before the test, n elements are inserted in the same way as in the random full inserts test.
    Then a another vector of n random elements different from the inserted elements is generated
    which is tried to search in the hash map.
    ''',
    "RandomFullReadsAfterDeletingHalf": '''
    Before the test, n elements are inserted in the same way as in the random full inserts test
    before deleting half of these values randomly. Then all the original values are tried to read
    in a different order, which will lead to 50% hits and 50% misses.
    ''',
    "RandomFullIteration":'''
    Before the test, n elements are inserted in the same way as in the random full inserts test.
    Then the hash map iterators is used to read all the key-value pairs.
    ''',
    "_50Reads_25Inserts_25Deletes":'''
    Before the test, a vector with n random values is generated, but only n/2 elements are inserted.
    Then the full vector is shuffled and randomly processed where 50% reads, 25% inserts, 25% deletes 
    operations are executed (successful vs unsuccessful rate 50/50). That benchmark seems to be the
    closest to reality.
    ''',
    "MemoryConsumption":'''
    Memory consumption of the random insert benchmark.
    '''
}

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
        benchName = firstRaw[0].replace('Benchmark','')
        annotationList = firstRaw[1].split('-')
        mapName = annotationList[0]
        n = int(annotationList[1])
        time_ms = int(lineRaw[2].strip().split(' ')[0]) / (1000 * 1000)
        load = 'load=' + lineRaw[4].strip().split(' ')[0]
        mapping[benchName][mapName].append((n,time_ms,load))
        if "U32RandomFullInserts" in benchName:
            memory_bytes = int(lineRaw[3].strip().split(' ')[0]) / (1024 * 1024)
            mapping["MemoryConsumptionU32"][mapName].append((n,memory_bytes,load))
        if "U64RandomFullInserts" in benchName:
            memory_bytes = int(lineRaw[3].strip().split(' ')[0]) / (1024 * 1024)
            mapping["MemoryConsumptionU64"][mapName].append((n,memory_bytes,load))
        if "UUIDRandomInserts" in benchName:
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
        <title>Golang Hashmap Benchmark</title>
        <script src='https://cdn.plot.ly/plotly-2.20.0.min.js'></script>
    </head>
    <body>
        <h1 align="center">Golang Hashmap Benchmark</h1>
        <hr>
    ''')

    for benchmark in sorted(mapping):
        b = mapping[benchmark]
        fd_out.write("<div id='"+benchmark+"'><script>\n")
        names = []
        y_naming = "time (ms)"
        if "MemoryConsumption" in benchmark:
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
        fd_out.write("var layout_" + benchmark + " = {title:'" + benchmark + "', xaxis: {title: 'number of entries in hash table (n)'},yaxis: {title: '" + y_naming + "'}};\n");
        fd_out.write("Plotly.newPlot('" + benchmark + "', data_"+ benchmark + ", layout_" + benchmark + ");\n"),
        fd_out.write("</script></div>")
        info_name = benchmark.replace('U64','').replace('U32','').replace('UUID','')
        fd_out.write('<center><p style="width: 700px;padding: 20px;"> '+info[info_name]+' </p></center>\n')
        fd_out.write('<hr>\n')

    # fd_out.write rest of body
    fd_out.write("</body>\n")


if __name__ == "__main__":
    main()
