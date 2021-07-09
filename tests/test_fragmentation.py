import subprocess
import platform
import unittest
from unittest import TestCase
import filecmp


class TestFragmentation(TestCase):
    def test_simple1(self):
        svl_command = "run ['../fragmentation.svl',[" \
                      "test:'fragmentation'," \
                      "moe:'resources/test.moe'," \
                      "renderer:'../bin/fill_template.win64.exe'," \
                      "template:'../templates/sample.ajf'," \
                      "fbase:'test'," \
                      "pdb:'temp/test.pdb'," \
                      "ajf:'temp/test.ajf'," \
                      "basis_set:'6-31g*']" \
                      "]".replace('/', '\\\\' if platform.system() == 'Windows' else '/')
        subprocess.run(["moebatch.exe", "-exec", svl_command])
        self.assertTrue(len(filecmp.cmpfiles('temp', "references/test1", ["test.pdb", "test.ajf"])[1]) == 0)


if __name__ == "__main__":
    unittest.main()