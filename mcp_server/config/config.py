import os
import logging

class Config:
    SAFELINE_ADDRESS: str
    SAFELINE_API_TOKEN: str
    SECRET: str
    LISTEN_PORT: int
    LISTEN_ADDRESS: str
    DEBUG: bool

    def __init__(self):
        set_log_level()

        if os.getenv("MCP_SERVER_DEBUG"):
            self.DEBUG = True
        else:
            self.DEBUG = False
        self.SAFELINE_ADDRESS = os.getenv("SAFELINE_ADDRESS")
        if self.SAFELINE_ADDRESS:
            self.SAFELINE_ADDRESS = self.SAFELINE_ADDRESS.removesuffix("/")
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


def set_log_level():
    level = logging.WARN
    log_level = os.getenv("MCO_SERVER_LOG_LEVEL")
    if log_level:
        match log_level.lower():
            case "debug":
                level = logging.DEBUG
            case "info":
                level = logging.INFO
            case "warn":
                level = logging.WARN
            case "error":
                level = logging.ERROR
            case "critical":
                level = logging.CRITICAL

    logging.basicConfig(level=level, format="%(asctime)s - %(levelname)s - %(message)s")