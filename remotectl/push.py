#!/usr/bin/env python3
import os
import sys

from common import cluster

__author__ = 'yunfhu'
rws = "/root/rebuild_jars/"
cms = ["python /root/rebuild_jars/replaceJar"]


def env():
    os.environ["JAVA_HOME"] = "/usr/java/jdk1.8.0_172-amd64"
    os.environ["PERL5LIB"] = "..:%s/src/common/cfg/:%s/src/msc/java/msc/rc/model/scripts/" % (lws, lws)


def push():
    for server in cls:
        jar_dir = "%s/build/dist/%s/jars/" % (lws, server.get_type())
        server.get_sftp().upload(jar_dir, rws, format=(lambda x: x.endswith(".jar") or x.endswith(".manifest")))


def replace():
    for server in cls:
        server.get_ssh().exec_command(cms)


def build():
    for mdl in set([s.get_type() for s in cls]):
        bf = "%s/src/%s/java/build.xml stageall" % (lws, mdl)
        rs = os.system("ant -f %s" % bf)
        if rs is not 0:
            sys.exit()


if __name__ == '__main__':
    # flows = [env, push, replace]
    # lws = "/home/yunfei/IdeaProjects/HOUSTON_SCEF_AAC_yunfei"
    lws = "/".join(os.path.abspath(__file__).split("/")[:-1])
    flows = [env, build, push, replace]
    cls = cluster.MultiplePeer("deploy", rws)
    for f in flows:
        f()
    cls.close()
