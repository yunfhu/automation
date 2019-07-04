import logging
import os
import sys
import threading

from common import config, client

__author__ = 'yunfhu'


class MultiplePeer(object):
    def __init__(self, cluster_name, rws=None):
        self.cluster_name = cluster_name
        self.config = config.ClusterConfig(self.cluster_name).get()
        if rws:
            self.servers = [AutoSniffPeer(k, v, rws) for k, v in self.config.items()]
        else:
            self.servers = [Peer(k, v) for k, v in self.config.items()]

    def close(self):
        for s in self.servers:
            s.close()

    def __getitem__(self, item):
        return self.servers[item]


class Peer(object):

    def __init__(self, host, config) -> None:
        super().__init__()
        self.host = host
        self.type = config[0]
        self.config = tuple([host] + config[1:])
        self.__log__()

    def __log__(self):
        self.logger = logging.getLogger(self.__str__())
        self.logger.setLevel(logging.INFO)
        console = logging.StreamHandler(sys.stdout)
        console.setFormatter(logging.Formatter('%(asctime)s %(levelname)-8s: %(message)s'))
        self.logger.addHandler(console)

    def __str__(self) -> str:
        return "peer for host:%s" % self.host

    def get_ssh(self):
        if hasattr(self, "ssh") is False:
            self.ssh = client.SSHClient(self.config)
        return self.ssh

    def ssh_close(self):
        if hasattr(self, "ssh"):
            self.ssh.close()

    def get_sftp(self):
        if hasattr(self, "sftp") is False:
            self.sftp = client.SFTPClient(self.config)
        return self.sftp

    def sftp_close(self):
        if hasattr(self, "sftp"):
            self.sftp.close()

    def get_type(self):
        return self.type

    def close(self):
        self.sftp_close()
        self.ssh_close()


class AutoSniffPeer(Peer):

    def __init__(self, host, config, rws) -> None:
        super().__init__(host, config)
        threading.Thread(target=self.sniffer(rws), args=("sniffer",)).start()

    def sniffer(self, rws):
        self.logger.info("start sniffer thread to check whether prepare setup on host:%s is okay ", self.host)
        detect = "if test ! -d %s; then mkdir %s; fi; if test ! -e %s/replace; then echo 0; fi" % (
            rws, rws, rws)
        stdin, stdout, stderr = self.get_ssh().client.exec_command(detect)
        ans = stdout.read().decode().replace("\n", "")
        if ans:
            self.logger.warning("prepare setup is not okay on host:%s", self.host)
            self.get_sftp().scp(os.path.join(os.path.dirname(__file__), '../replace'), "%s/replace" % rws)
            self.logger.info("prepare setup on host:%s done ", self.host)
        else:
            self.logger.info(" prepare setup on host:%s is okay ", self.host)


if __name__ == '__main__':
    host = "127.0.0.1"
    config = ["fe", 22, "yunfei", "NextGen"]
    rws = "/home/yunfei/test/"
    server = AutoSniffPeer(host, config, rws)
    server.close()
