#
# Copyright 2014 Nathaniel Cook
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


from morgoth.config import Config

import os
import tempfile
import unittest

class TestConfig(unittest.TestCase):
    """
    morgoth.config.Config unittest
    """
    simple_yaml = """
    test_key0: value0
    test_key1:
        - li0
        - li1
        - li2
    test_key2:
        k0: v0
        k1: v1
        k2: v2
        k3: { inline_key: value }
    """
    include_yaml = """
    included: !include included.yml
    """
    included_yaml = """
    k0: v0
    k1: v1
    k2: v2
    k3: v3
    """

    def test_conifg_00(self):
        """
        Test basic parsing and access patterns
        """
        config = Config.loads(self.simple_yaml)

        # Test different access patterns
        self.assertEqual('value0', config.test_key0)
        self.assertEqual('value0', config.get('test_key0', None))
        self.assertEqual('value0', config.get(['test_key0'], None))


        # Test list access
        list_data = ['li0', 'li1', 'li2']
        i = 0
        for li in config.test_key1:
            self.assertEqual(list_data[i], li)
            i += 1

        # Test dictionary access
        self.assertEqual('v0', config.test_key2.k0)
        self.assertEqual('v1', config.test_key2.k1)
        self.assertEqual('v2', config.test_key2.k2)
        self.assertEqual('value', config.test_key2.k3.inline_key)

        # Test defaulting behaviors
        self.assertEqual(1, config.get('nonexistant', 1))
        self.assertEqual(1, config.nonexistant)




    def test_config_01(self):
        """
        Test '!include' behavior
        """
        include = None
        included = None
        tdir = tempfile.mkdtemp()
        try:

            include = open(os.path.join(tdir, 'include.yml'), 'w')
            include.write(self.include_yaml)
            include.close()

            included = open(os.path.join(tdir, 'included.yml'), 'w')
            included.write(self.included_yaml)
            included.close()

            config = Config.load(include.name)

            self.assertEqual('v0', config.included.k0)
            self.assertEqual('v1', config.included.k1)
            self.assertEqual('v2', config.included.k2)
            self.assertEqual('v3', config.included.k3)

        finally:
            if include:
                os.remove(include.name)
            if included:
                os.remove(included.name)
            os.rmdir(tdir)

if __name__ == '__main__':
    unittest.main()

