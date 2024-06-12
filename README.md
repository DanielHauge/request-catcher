# Request catcher

This small project aims to spin up a small server that will catch all incoming requests and log them to files.
Request headers and body will be logged to a file in the logs folder.

```bash
docker run -dit --rm --name rc -p 8080:8080 --volume $(pwd)/logs:/go/logs  danielhauge/request-catcher
```

**note: pwd may require a -W flag on windows**
