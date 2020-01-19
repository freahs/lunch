FROM ubuntu
RUN apt-get update && apt-get install -y python-dev libxml2-dev libxslt1-dev antiword unrtf poppler-utils pstotext tesseract-ocr \
    flac ffmpeg lame libmad0 libsox-fmt-mp3 sox libjpeg-dev swig python-pip curl && apt-get clean
RUN pip install textract
ENV TESSDATA_PREFIX=/textract/tessdata
RUN mkdir -p ${TESSDATA_PREFIX} &&\
    cd ${TESSDATA_PREFIX} &&\
    curl -L https://github.com/tesseract-ocr/tessdata/blob/master/swe.traineddata?raw=true -o swe.traineddata
RUN mkdir /downloads /scripts && echo '#!/bin/bash\n\
url=$1\n\
file=${url##*/}\n\
[[ ! -f /downloads/$file ]] && cd /downloads && curl -OL $url\n\
textract -m tesseract -O language=swe /downloads/$file' > /scripts/run.sh && chmod +x /scripts/run.sh