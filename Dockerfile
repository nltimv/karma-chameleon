FROM python:latest

WORKDIR /app

COPY requirements.txt /app/

RUN pip install -r requirements.txt

COPY karmabot.py /app/

ENTRYPOINT [ "python", "-u", "/app/karmabot.py" ]
