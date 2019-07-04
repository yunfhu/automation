import logging
import os
import sys
import paramiko
from abc import abstractmethod

__author__ = 'yunfhu'


class ClientBase(object):
    def __init__(self, config, client=None):
        self.config = config
        self.__setup__(self.config, client)
        self.__log__()
        self.logger.debug("%s connected", self)

    def __log__(self):
        self.logger = logging.getLogger(self.__str__())
        self.logger.setLevel(logging.INFO)
        console = logging.StreamHandler(sys.stdout)
        console.setFormatter(logging.Formatter('%(asctime)s %(levelname)-8s: %(message)s'))
        self.logger.addHandler(console)

    @abstractmethod
    def __setup__(self, config, client):
        pass

    def close(self):
        if self.client:
            self.client.close()
            self.logger.debug("%s disconnected", self)


class SSHClient(ClientBase):
    def __setup__(self, config, client):
        client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        client.connect(*self.config)
        self.client = client

    def __init__(self, config):
        super().__init__(config, paramiko.SSHClient())

    def __str__(self) -> str:
        return "ssh client for host %s" % self.config[0]

    def exec_command(self, commands):
        try:
            for command in commands:
                stdin, stdout, stderr = self.client.exec_command(command)
                self.logger.info("executed command:%s on host: %s", command, self.config[0])
                for out in stdout:
                    self.logger.info(out.replace("\n", ""))
                for err in stderr:
                    self.logger.warning(err.replace("\n", ""))
        except Exception as e:
            self.logger.warning(e)


class SFTPClient(ClientBase):
    def __setup__(self, config, client=None):
        transport = paramiko.Transport(*self.config[:1])
        transport.connect(username=self.config[2], password=self.config[3])
        self.client = paramiko.SFTPClient.from_transport(transport)

    def __init__(self, config):
        super().__init__(config)

    def __str__(self) -> str:
        return "sftp client for host %s" % self.config[0]

    def download(self, *args, format=None):
        remote_path = args[0]
        local_repository = args[1]
        remote_files_attr = list(filter(format, self.client.listdir_attr(remote_path)))
        try:
            for file_attr in remote_files_attr:
                self.client.get(os.path.join(remote_path, file_attr.filename),
                                os.path.join(local_repository, file_attr.filename))
                self.logger.info("download %s to %s", file_attr, local_repository)
        except Exception as e:
            self.logger.warning(e)

    def upload(self, *args, format=None):
        lws = args[0]
        rws = args[1]
        self.logger.info("sync %s with %s:/%s" % (lws, self.config[0], rws))
        lws_files = list(filter(format, os.listdir(lws)))
        try:
            for file in lws_files:
                self.client.put(os.path.join(lws, file), os.path.join(rws, file))
                self.logger.info("upload %s%s", lws, file)
        except Exception as e:
            self.logger.warning(e)

    def scp(self, sour, dest):
        try:
            self.client.put(sour, dest)
        except Exception as e:
            self.logger.warning(e)
