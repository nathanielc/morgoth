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

def get_col_for_metric(db, metric):
    """
    Return the collection connection for a given metric


    NOTE: the current implementation uses just one collection for all metric data.
    I plan to change this use multiple collections later

    @param db: the db connection to use
    @param metric: the name of the metric
    """
    return db.metrics
