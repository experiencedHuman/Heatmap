FROM alexcpn/fb_prophet_python:3.7.13-buster

RUN pip install requests

WORKDIR /usr/app/src

COPY forecasting.py ./

CMD ["python", "./forecasting.py"]