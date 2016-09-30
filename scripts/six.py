import sys
from math import sqrt

def read_matrix(fp):
    rows = []
    while True:
        line = fp.readline().strip()
        if line == "": break

        row = [float(x) for x in line.replace("\t", " ").split(" ") if x != ""]
        rows.append(row)
    if rows == []:
        return None
    return rows

def get_n(n, hist):
    for row in hist:
        if row[1] == n:
            return row[2]
    return 0
get_six = lambda h: get_n(6, h)

def vsub(x, y):
    diff = []
    n = min(len(x), len(y))
    for i in range(0, n):
        diff.append(x[i] - y[i])
    return diff

if __name__ == '__main__':
    base = sys.argv[1]

    hists = open(base, "r")
    read_matrix(hists) ## throw away the first

    mats = open(base + ".mat", "r")

    last_hist = read_matrix(hists)
    last_mat  = read_matrix(mats)
    last_six  = get_six(last_hist)
    last_tr6  = last_mat[6]

    i = 1
    while True:
        hist = read_matrix(hists)
        mat  = read_matrix(mats)
        if hist is None or mat is None:
            break

        six = get_six(hist)

        tr6 = mat[6]
        ntr6 = sqrt(sum(x**2 for x in vsub(tr6, last_tr6)))

        print "%d\t%f\t%f" % (i, six - last_six, ntr6)

        last_six = six
        last_tr6 = tr6
        i += 1

