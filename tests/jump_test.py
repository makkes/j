from collections import Counter
import jump


def sort_history_prefers_exact_matches_test():
    history = Counter({
        "/home/max": 3,
        "/home/max/dev/boplish/core": 4
    })
    res = jump.sort_history("/home/max", history, True)
    assert res[0][0] == "/home/max"
    res = jump.sort_history("home/max", history, True)
    assert res[0][0] == "/home/max"


def sort_history_prefers_suffix_match_test():
    history = Counter({
        "/home/max": 3,
        "/home/max/dev/boplish/core": 4
    })
    res = jump.sort_history("max", history, True)
    assert res[0][0] == "/home/max"
    res = jump.sort_history("/max", history, True)
    assert res[0][0] == "/home/max"
    res = jump.sort_history("e/max", history, True)
    assert res[0][0] == "/home/max"


def sort_history_test():
    history = Counter({
        "/home/max": 3,
        "/home/max/dev/boplish/core": 4
    })
    res = jump.sort_history("/max", history, True)
    assert res[0][0] == "/home/max"
