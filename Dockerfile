FROM iron/base

COPY ./lingva-server ./lingva-server

ENV PORT 8080
ENV IMAGE_URL http://138.68.70.81/image/


EXPOSE 80

CMD ["./lingva-server", "start"]
