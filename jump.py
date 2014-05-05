#!/usr/bin/env python3

import os
import pickle
import re
from collections import Counter
from datetime import datetime

DATFILE = os.getenv('HOME') + '/.jump.dat'
LOGFILE = os.getenv('HOME') + '/.jump.log'
SPECIALDIRS = ['.', '..']


def log(msg):
    f = open(LOGFILE, 'a')
    f.write(str(datetime.today()) + ': ' + str(msg))
    f.write('\n')


def write_history(obj):
    f = open(DATFILE, 'wb')
    pickle.dump(obj, f)
    f.close()


def read_history():
    if not os.path.exists(DATFILE):
        return Counter()
    f = open(DATFILE, 'rb')
    obj = pickle.load(f)
    f.close()
    return obj


def sort_history(history):
    return sorted(history.items(), key=lambda e: e[1], reverse=True)


def delete_entry(entry):
    print('deleting "' + entry + '"')
    history = read_history()
    del history[entry]
    write_history(history)


def find_dir(pattern, history):
    p = re.compile(pattern)
    for d, cnt in sort_history(history):
        if p.search(d) and os.path.exists(d):
            yield d


def strip_common_prefix(l):
    if len(l) <= 1:
        return l
    prefix = os.path.commonprefix(l)
    # don't strip beginning of directory names
    prefix = prefix[:prefix.rfind('/')+1]
    prefixlen = len(prefix)
    return list(map(lambda s:
                    s[prefixlen:]
                    # only strip prefix if the result wouldn't be empty
                    if len(s[prefixlen:]) > 0
                    else s,
                    l))


def strip_leading_slash(l):
    if len(l) <= 1:
        return l
    return list(map(lambda s: s.lstrip('/'), l))


def only_special_dirs(name):
    return re.compile('^[./]+$').match(name) is not None


def find_candidates(name, history, canonical):
    try:
        # a directory that is reachable via 'name' has precedence
        islink = os.path.islink(name)
        oldcwd = os.getcwd()
        os.chdir(name)
        cwd = os.getcwd()
        # return to the previous directory so completion works upon continuous
        # calls to this function
        os.chdir(oldcwd)
        # if we followed a symlink, return that, else return the CWD
        return [cwd if not islink else
                os.path.join(os.path.split(cwd)[0], name)]
    except OSError:
        # a directory 'name' doesn't exist so carry on
        pass
    if not only_special_dirs(name):
        candidates = []
        for found in find_dir(name, history):
            candidates.append(found)
        if not canonical:
            candidates = strip_leading_slash(strip_common_prefix(candidates))
        return sorted(candidates, key=lambda e: e.count('/'))
    return [name]


def find_best_dir(d):
    history = read_history()
    candidates = find_candidates(d, history, True)
    if not candidates:
        return d
    best_dir = candidates[0]
    history[best_dir] += 1
    write_history(history)
    return best_dir


def complete(d):
    res = find_candidates(d, read_history(), False)
    width = len(str(len(res)-1))
    if len(res) > 1:
        res = list(map(
            lambda e:
            "{:0{width}} {}".format(res.index(e), e, width=width), res))
    print("\n".join(res))


def dump():
    history = read_history()
    width = len(str(history.most_common(1)[0][1]))
    for directory, count in sort_history(history):
        print("{:{width}} {}".format(count, directory, width=width))
