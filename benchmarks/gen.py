from os.path import exists
import math
import json

sizes =[1,5,10,100]

line = json.dumps({
    "id": 12345,
    "name": "very_long_string_with_escapes_and_unicode_abcdefghijklmnopqrstuvwxyz_0123456789",
    "description": "This string contains\nmultiple\nlines\nand \"quotes\" and unicode ❤❤❤",
    "nested": {
        "level1": {
            "level2": {
                "level3": {
                    "level4": {
                        "array": [
                            "short",
                            "string_with_escape\\n",
                            "another\\tvalue",
                            "unicode\u2603",
                            "escaped_quote_\"_and_backslash_\\",
                            1234567890,
                            -1.2345e67,
                            3.1415926535897932384626433832795028841971,
                            True,
                            False,
                            None,
                            "\u0041\u0042\u0043\u00A9\u20AC",
                            "mix\\n\\t\\r\\\\\\\"end"
                        ]
                    }
                }
            }
        }
    }
})

def write_data(size: int): 
    name = f"{size}MB.json"
    if not exists(name):
        with open(name, mode="w", encoding="utf8") as f:
            f.write("[\n")
            size = math.floor((size*1000000)/len(line))
            f.write(",\n".join([line for _ in range(0, size)]))
            f.write("\n]")

[write_data(size) for size in sizes]
