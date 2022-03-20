FROM scratch

# Copy src files
COPY ./words-of-wisdom-client /words-of-wisdom-client

# Run app
CMD ["/words-of-wisdom-client"]
