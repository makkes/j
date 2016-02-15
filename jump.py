#!/usr/bin/env python3

import os
import pickle
import re
import json
from collections import Counter
from datetime import datetime

JSONFILE = os.getenv('HOME') + '/.jump.json'
LOGFILE = os.getenv('HOME') + '/.jump.log'
SPECIALDIRS = ['.', '..']


def log(msg):
    f = open(LOGFILE, 'a')
    f.write(str(datetime.today()) + ': ' + str(msg))
    f.write('\n')


def migrate_old_history():
    OLD_DATFILE = os.getenv('HOME') + '/.jump.dat'
    if not os.path.exists(OLD_DATFILE):
        return
    f = open(OLD_DATFILE, 'rb')
    obj = pickle.load(f)
    f.close()
    os.remove(OLD_DATFILE)
    write_history(obj)


def write_history(obj):
    f = open(JSONFILE, 'w', encoding="utf-8")
    f.write(json.dumps(obj, indent=4))
    f.close()


def read_history():
    migrate_old_history()

    if not os.path.exists(JSONFILE):
        return Counter()
    f = open(JSONFILE, 'r', encoding="utf-8")
    history = Counter(json.load(f, encoding="utf-8"))
    f.close()
    return history


def history_sorter(pattern):
    stripped_pattern = pattern.lstrip("/")

    def key_func(element):
        return int(
            element[0].lstrip("/") == stripped_pattern or
            element[0].split("/")[-1] == stripped_pattern), element[1]
    return key_func


def sort_history(pattern, history, reverse):
    return sorted(history.items(), key=history_sorter(pattern),
                  reverse=reverse)


def delete_entry(entry):
    print('deleting "' + entry + '"')
    history = read_history()
    del history[entry]
    write_history(history)


def find_dir(pattern, history):
    p = re.compile(pattern)
    for d, cnt in sort_history(pattern, history, True):
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
        # a directory that is reachable from CWD via 'name' has precedence
        islink = os.path.islink(name)
        oldcwd = os.getcwd()
        os.chdir(name)
        cwd = os.getcwd()
        #cwd = os.getcwd()
        # return to the previous directory so completion works upon
        # continuous calls to this function
        os.chdir(oldcwd)
        # if we followed a symlink, return that, else return the CWD
        return [cwd if not islink else
                os.path.join(oldcwd, name)]
    except OSError:
        # a directory 'name' doesn't exist so carry on
        pass
    if not only_special_dirs(name):
        candidates = []
        for found in find_dir(name, history):
            candidates.append(found)
        if not canonical:
            candidates = strip_leading_slash(strip_common_prefix(candidates))
        return candidates
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


def dump(reverse):
    history = read_history()
    width = len(str(history.most_common(1)[0][1]))
    for directory, count in sort_history("", history, reverse):
        print("{:{width}} {}".format(count, directory, width=width))
