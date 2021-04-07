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

"""Implementation of test steps that run Insights Aggregator Mock service and check its output."""

import subprocess


def process_executable_output(context, out, return_code):
    """Process cleaner output."""
    assert out is not None

    # interact with the process:
    # read data from stdout and stderr, until end-of-file is reached
    stdout, stderr = out.communicate()

    # basic checks
    assert stderr is None, "Error during check"
    assert stdout is not None, "No output from cleaner"

    # check the return code of a process
    assert out.returncode == 0 or out.returncode == return_code, \
        "Return code is {}".format(out.returncode)

    # try to decode output
    output = stdout.decode('utf-8').split("\n")

    assert output is not None

    # update testing context
    context.output = output
    context.stdout = stdout
    context.stderr = stderr


@when(u"I start the mock service with the {flag} command line flag")
def start_mock_service_with_flag(context, flag):
    """Start the mock service with given command-line flag."""
    out = subprocess.Popen(["insights-results-aggregator-mock", flag],
                           stdout=subprocess.PIPE,
                           stderr=subprocess.STDOUT)

    assert out is not None
    process_executable_output(context, out, 2)


@then(u"I should see help messages displayed on standard output")
def check_help_from_cleaner(context):
    """Check if help is displayed by cleaner."""
    expected_output = """
Service to provide content for OCP rules

Usage:

    insights-results-aggregator-mock [command]

The commands are:

    <EMPTY>                      starts content service
    start-service                starts content service
    help     print-help          prints help
    config   print-config        prints current configuration set by files & env variables
    version  print-version-info  prints version info
    authors  print-authors       prints authors

"""

    assert context.stdout is not None
    stdout = context.stdout.decode("utf-8").replace("\t", "    ")

    # preliminary checks
    assert stdout is not None, "stdout object should exist"
    assert type(stdout) is str, "wrong type of stdout object"

    # check the output
    assert stdout.strip() == expected_output.strip(), "{} != {}".format(stdout, expected_output)


@then(u"I should see version info displayed on standard output")
def check_version_from_cleaner(context):
    """Check if version info is displayed by cleaner."""
    # preliminary checks
    assert context.output is not None
    assert type(context.output) is list, "wrong type of output"

    # check the output
    assert "Version:\t0.1" in context.output, \
        "Caught output: {}".format(context.output)


@then(u"I should see info about authors displayed on standard output")
def check_authors_info_from_cleaner(context):
    """Check if information about authors is displayed by cleaner."""
    # preliminary checks
    assert context.output is not None
    assert type(context.output) is list, "wrong type of output"

    # check the output
    assert "Authors:" in context.output, \
        "Authors: header is expected"

    assert "Pavel Tisnovsky <ptisnovs@redhat.com>" in context.output, \
        "Caught output: {}".format(context.output)
