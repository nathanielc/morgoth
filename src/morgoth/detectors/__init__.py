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

from morgoth.plugin_loader import PluginLoader
import os

import logging

logger = logging.getLogger(__name__)

_LOADER = None

def configure_detectors(app):
    """
    Configure the detector plugin loader

    @param app: the morgoth application instance
    @type app: morgoth.app.App
    """
    global _LOADER
    from morgoth.detectors.detector import Detector
    dirs = [os.path.dirname(__file__)]
    dirs.extend(app.config.get(['plugin_dirs', 'detectors'], []))
    _LOADER = PluginLoader(app, dirs, Detector)

def get_loader():
    """
    Return the Detectors plugin loader

    configure_detectors must be called before this method.
    """
    assert _LOADER is not None
    return _LOADER
