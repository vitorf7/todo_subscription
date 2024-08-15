# Testing gqlgen issue
After cloning the repo, run:

```bash
go mod download

make dev
```

Then to test run the following curl command in a new terminal:

```bash
curl -N --request POST --url http://localhost:8888/graphql \
    -H "accept: text/event-stream" \
    -H 'content-type: application/json' \
    --verbose \
    --data '{"query":"subscription {todoState { notes { ID}}}"}'
```

The issue is mimicking a gRPC stream that returns a stream of events and then the GraphQL server would send the new data
via a channel for the GraphQL subscription.
