import re
from argparse import ArgumentParser

p = ArgumentParser()
p.add_argument("svl")
a = p.parse_args()

re_docs_header = re.compile("^/\*(.+?)\*/", re.MULTILINE | re.DOTALL)
re_function = re.compile("function (.+?) +(.+?)\n.*?/\*(.+?)\*/.+?endfunction", re.MULTILINE | re.DOTALL)

with open(a.svl, "r") as f:
    body = f.read()
    body = re.sub("function .+;", "", body)

with open(a.svl.replace(".svl", ".py"), "w") as f:
    try:
        header = next(re_docs_header.finditer(body))
        f.write('"""' + header[1] + '\n"""\n')
    except:
        pass

    for line in re_function.findall(body):
        f.write("def " + line[0] + "(" + line[1].strip("[]") + "):")
        f.write('\n\t"""' + line[2] + '\n\t"""\n\tpass\n')

