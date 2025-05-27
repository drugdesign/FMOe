import re
from argparse import ArgumentParser

p = ArgumentParser()
p.add_argument("svl")
a = p.parse_args()

data = []
level = ""
with open(a.svl, "r") as f:
    for line in f:
        line = line.strip()
        if "function " in line:
            if ";" in line:
                continue
            level = "\t"
            temp = line.split("function ")[1].strip()
            func_title = temp.split(" ")[0]
            func_attr = temp.replace(func_title, "").strip(" []")
            data.append(f"def {func_title} ({func_attr}):")
            
        if "/*" in line:
            data.append(level + '"""' + line.replace("/*", ""))
            for line_doc in f:
                line_doc = line_doc.strip()
                if "*/" in line_doc:
                    data.append(level + line_doc.replace("*/", ""))
                    data.append(level + '"""')
                    break
                data.append(level + line_doc)
        
        if "endfunction" in line:
            data.append("\tpass")
            level = ""

with open(a.svl.replace(".svl", ".py"), "w") as f:
    f.write("\n".join(data))


