import json
import os

__author__ = 'yunfhu'


class FileLoader(object):
    def __load__(self):
        cf = open(os.path.join(os.path.dirname(__file__), './host.json'))
        try:
            self.config = json.load(cf)
        finally:
            cf.close()

    def __init__(self):
        self.__load__()

    def get(self, cls=None):
        if cls:
            return dict(self.config.get(cls))
        else:
            return dict(self.config)


class ClusterConfig(object):
    def __init__(self, cluster_key):
        self.cluster_key = cluster_key
        self.loader = FileLoader()
        self.config = {k: v[:-1] for k, v in self.loader.get(self.cluster_key).items() if v[4] is True}

    def get(self, host=None):
        if host:
            return self.config.get(host)
        else:
            return self.config


if __name__ == '__main__':
    loader = ClusterConfig("deploy")
    print(loader.get())
