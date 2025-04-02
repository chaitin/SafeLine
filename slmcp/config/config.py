import os

class Config:
    SAFELINE_ADDRESS: str
    SAFELINE_API_TOKEN: str
    SECRET: str
    LISTEN_PORT: int
    LISTEN_ADDRESS: str

    def __init__(self):
        self.SAFELINE_ADDRESS = os.getenv("SAFELINE_ADDRESS")
        self.SAFELINE_API_TOKEN = os.getenv("SAFELINE_API_TOKEN")
        self.SECRET = os.getenv("SAFELINE_SECRET")
        env_listen_port = os.getenv("LISTEN_PORT")
        if env_listen_port and env_listen_port.isdigit():
            self.LISTEN_PORT = int(env_listen_port)
        else:
            self.LISTEN_PORT = 5678
        env_listen_address = os.getenv("LISTEN_ADDRESS")
        if env_listen_address:
            self.LISTEN_ADDRESS = env_listen_address
        else:
            self.LISTEN_ADDRESS = "0.0.0.0"

    @staticmethod
    def from_env():
        return Config()