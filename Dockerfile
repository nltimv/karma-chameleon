FROM python:latest

WORKDIR /app

COPY src/* /app/

RUN pip install -r requirements.txt

ENTRYPOINT [ "python", "-u", "/app/karmabot.py" ]
