## Dash

Docker build
`docker build -t dash .`

Docker run
`docker run -dp 8090:8090 dash`

Deploy
`fly deploy`

Secrets
`fly secrets import < .env.production`