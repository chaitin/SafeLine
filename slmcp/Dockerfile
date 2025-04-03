FROM ghcr.io/astral-sh/uv:python3.13-alpine

COPY requirements.txt /app/requirements.txt

WORKDIR /app

RUN uv pip sync requirements.txt --system -i https://pypi.tuna.tsinghua.edu.cn/simple

COPY . /app

CMD ["python", "__main__.py"]