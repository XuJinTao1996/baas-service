From debian

COPY ./conf /home/conf
COPY ./assets /home/assets
COPY ./db /home/db
COPY ./app /home/app

WORKDIR /home

ENTRYPOINT /home/app