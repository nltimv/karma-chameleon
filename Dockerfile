FROM python:latest

WORKDIR /app

COPY src/requirements.txt /app/

RUN pip install -r requirements.txt

COPY src/* /app/

ENTRYPOINT [ "python", "-u", "/app/karmabot.py" ]
