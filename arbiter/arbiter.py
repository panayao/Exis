#!/usr/bin/python
"""
This program provides the final word of how each programming language
should implement function calls while using Exis, so the name is actually quite fitting.

It can document and test real live examples of how to use Exis for every langugage.

Please run '$0 -ls all' for more info.

Environment Variables:
    EXIS_REPO - the path to the Exis repo

TODO:
    - Implement Object checking for Python
    - Implement other languages
"""
import sys
import os
import time
import glob
import argparse
import re
import platform
from collections import defaultdict as ddict

import multiprocessing
import threading
from multiprocessing.dummy import Pool as ThreadPool


EXISREPO = os.environ.get("EXIS_REPO", None)

if EXISREPO is None:
    print("!" * 50)
    print("!! $EXIS_REPO not found, this may not work")
    print("!" * 50)
    sys.path.append("..")
else:
    sys.path.append(EXISREPO)

from utils import functionizer as funcizer
from utils import utils

import exampler
import repl

if platform.system() == "Darwin":
    repl.STUB_REPL = True
    sys.stderr.write("Warning: Darwin detected, switching REPL scripts to stub\n")
else:
    repl.STUB_REPL = False


def findTasks(lang=None, task=None, verbose=False, shouldPrint=True):
    """
    Searches for all example files in the Exis repo.
    Args:
        lang : One of {python, go, js, swift} or None which means get all.
        task : Matching task with wildcard support (ie. "Pub/Sub*")
        verbose : T/F on verbose printing
    """
    examples = exampler.Examples.find(EXISREPO, lang)
    allTasks = examples.getTasks(lang, task)

    if shouldPrint:
        for t in allTasks:
            if(verbose):
                print(t.details())
            else:
                print(t)

    return allTasks


def findTask(lang, task):
    """
    Finds and prints reference to a specific task in a specific language.
    Args:
        lang : lang to search for
        task : Task name
    """
    examples = exampler.Examples.find(EXISREPO, lang)
    ts = examples.getTask(task)
    if(ts):
        print(ts.details())
    else:
        print("No Task found")

TASK_DEF_RE = re.compile("(.*)? (.*):(.*)$")


def _ripTaskDef(t, kwargs):
    """
    Internal function that rips apart a task definition like "language action:example"
    """
    # Check for optional lang so you don't have to type it out
    l = kwargs.get('lang', None)
    a = "{} {}".format(l, t) if l else t
    m = TASK_DEF_RE.match(a)
    if not m:
        # Check if they used the lang, if so it can look different
        if l:
            return l, None, t
        else:
            return [None] * 3
    return m.groups()


def test(*tasks, **kwargs):
    """
    Executes potentially many tasks provided as individual arguments.
    NOTE: Please order your tasks intelligently - this means place subs/regs before calls/pubs.
    Arguments:
        tasks... : potentially many tasks to execute, in the format "language action:example name"
        -v       : if the last arg is -v then print extra data about the tasks found

    Example:
        test("python register:Reg/Call", "swift call:Reg/Call")
            This will setup the reg of the Reg/Call example in python and call it with the Reg/Call
            example from Swift.

        test("python register:Reg/Call Basic", "swift publish:Pub/Sub")
            This will setup the python reg of "Reg/Call Basic" and the swift publish from Pub/Sub
            obviously nothing will happen and that is the point - know what you are doing...
    """
    if(tasks[-1] == "-v"):
        tasks = tasks[:-1]
        verbose = True
    examples = exampler.Examples.find(EXISREPO)

    taskList = list()
    actionList = list()
    for t in tasks:
        lang, action, taskName = _ripTaskDef(t, kwargs)
        ts = examples.getTask(taskName, lang)
        if not ts:
            print("!! No TaskSet found")
        elif action is None:
            # This means we need to add each of the tasks from the taskset
            for t in ts.getOrderedTasks():
                taskList.append(t)
                actionList.append(t.action)
        else:
            taskList.append(ts.getTask(action))
            actionList.append(action)

    # Exec all of them
    if repl.executeTasks(taskList, actionList):
        exit(0)
    else:
        exit(1)


def testAll(lang, stopOnFail=False, multithreaded=False):
    """
    Executes all found tests for the language provided.
    Args:
        lang       : language to test, or "all"
        stopOnFail : If true, we stop testing when a failure is found, default False
        multithreaded: if True, run the tests in parallel
    """

    if lang == "all":
        langs = exampler.LANGS.keys()
    else:
        langs = [lang]

    if multithreaded:
        def runner(taskset):
            res = repl.executeTaskSet(taskset)
            if res is True:
                print('-' * 80)
            elif res is False:
                if stopOnFail:
                    exit()
                else:
                    print('-' * 80)

        for lang in langs:
            examples = exampler.Examples.find(EXISREPO, lang)
            tasks = examples.getTasks(lang)

            pool = ThreadPool(10)
            pool.map(runner, tasks)
            pool.close()
            pool.join()

            # threads = [threading.Thread(target=runner, args=(x,)) for x in tasks]
            # [x.start() for x in threads]
            # [x.join() for x in threads]
    else:
        hasFailed = False
        for lang in langs:
            examples = exampler.Examples.find(EXISREPO, lang)
            for t in examples.getTasks(lang):
                res = repl.executeTaskSet(t)
                if res is False:
                    hasFailed = True
                    if stopOnFail:
                        break

    if hasFailed:
        exit(1)
    else:
        exit(0)


def genTemplate(langs=exampler.LANGS.keys(), actions=["Pub/Sub", "Reg/Call"]):
    """
    Use the generator.js code to generate basic templates and print to stdout.
    Args:
        langs   : LIST of langs ["python", "swift", "js"]
        actions : LIST of actions ["Pub/Sub", "Reg/Call"]
    """
    print langs


def cleanup():
    """
    Cleans up all tmp test folders.
    """
    print "Cleaning up tmp directories:"
    repl.cleanupTests()
    print "DONE"


def genDocs():
    """
    Takes all code found and generates a JSON style doc system that can be consumed by the website.
    Looks like:
           {
            TaskName: {
                    lang: {
                        action: {
                            file: str,
                            lines: str,
                            code: ...,
                            expectType: str,
                            expectVal: str
                        }
                    }
                }
            }
    """
    examples = exampler.Examples.find(EXISREPO)

    docs = ddict(lambda: {k: dict() for k in exampler.LANGS.keys()})
    for t in examples.getTasks():
        for tt in t.tasks:
            d = dict(file=tt.fileName, lineStart=tt.lineStart, lineEnd=tt.lineEnd, code=tt.code,
                     expectType=tt.expectType, expectVal=tt.expectVal)
            docs[t.getName()][t.getLangName()][tt.action] = d

    # Strip out anything that isn't populated
    for k, v in docs.iteritems():
        l = list()
        for kk, vv in v.iteritems():
            if not vv or kk == "browser":
                l.append(kk)
        for ll in l:
            v.pop(ll)
        p = v.pop('nodejs', None)
        if p:
            v['js'] = p

    print utils.jsonPretty(docs)


def _getArgs():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('-node', help='Launch a node to test on', action='store_true')
    parser.add_argument('-debug', help='Debug mode', action='store_true')
    parser.add_argument('-v', '--verbose', help='Verbose mode', action='store_true')
    return parser

if __name__ == "__main__":
    parser = _getArgs()
    funcizer.init(parser)
    args = parser.parse_args()

    # Startup a node upon request
    if args.node:
        repl.launchNode()
    if args.debug:
        repl.debugMode()
    if args.verbose:
        repl.enableVerbose()

    # Now make the call that decides which of our functions to run
    funcizer.performFunctionalize(args, __name__, modSearch="__main__")

    if args.node:
        repl.killNode()
