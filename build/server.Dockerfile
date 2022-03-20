FROM scratch

# Copy src files
COPY ./words-of-wisdom-server /words-of-wisdom-server
COPY ./data/quotes.json /data/quotes.json

# Run app
CMD ["/words-of-wisdom-server"]
