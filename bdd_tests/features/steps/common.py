# Copyright © 2021 Pavel Tisnovsky, Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from shutil import which


from behave import given, then, when


@given(u"the system is in default state")
def system_in_default_state(context):
    """Check the default system state."""
    pass


@when(u"I look for executable file {filename}")
def look_for_executable_file(context, filename):
    """Try to find given executable file on PATH."""
    context.filename = filename
    context.found = which(filename)


@then(u"I should find that file on PATH")
def file_was_found(context):
    """Check if the file was found on PATH."""
    assert context.found is not None, \
        "executable filaname '{}' is not on PATH".format(context.filename)
