from collections import Counter
import jump
import unittest


class TestMatchingMethods(unittest.TestCase):

    def test_sort_history_prefers_exact_matches(self):
        history = Counter({
            "/home/max": 3,
            "/home/max/dev/boplish/core": 4
            })
        res = jump.sort_history("/home/max", history, True)
        self.assertEqual(res[0][0], "/home/max")
        res = jump.sort_history("home/max", history, True)
        self.assertEqual(res[0][0], "/home/max")

    def test_sort_history_prefers_suffix_match(self):
        history = Counter({
            "/home/max": 3,
            "/home/max/dev/boplish/core": 4
            })
        res = jump.sort_history("max", history, True)
        self.assertEqual(res[0][0], "/home/max")
        res = jump.sort_history("/max", history, True)
        self.assertEqual(res[0][0], "/home/max")
        res = jump.sort_history("e/max", history, True)
        self.assertEqual(res[0][0], "/home/max")

    def test_sort_history(self):
        history = Counter({
            "/home/max": 3,
            "/home/max/dev/boplish/core": 4
            })
        res = jump.sort_history("/max", history, True)
        self.assertEqual(res[0][0], "/home/max")


if __name__ == '__main__':
    unittest.main()
