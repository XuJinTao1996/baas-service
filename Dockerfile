From debian

COPY ./app /home
COPY ./conf/app.ini /etc/app.ini
COPY ./db/baas.db /home/baas.db

ENTRYPOINT /home/app