from os.path import exists
import math

sizes =[1,5,10]

line = """\t{
        "key1": "value",
        "array": [],
        "obj": {},
        "atomArray": [11201,1e112,true,false,null,"str"]
    }"""

def write_data(size: int): 
    name = f"{size}MB.json"
    if not exists(name):
        with open(name, mode="w", encoding="utf8") as f:
            f.write("[\n")
            size = math.floor((size*1000000)/len(line))
            f.write(",\n".join([line for _ in range(0, size)]))
            f.write("\n]")

[write_data(size) for size in sizes]
