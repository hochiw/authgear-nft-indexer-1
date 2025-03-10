# AuthGear NFT Indexer

## Prerequisite

1. Install asdf

2. Run the following to install all dependencies in .tool-versions

```
asdf install
```

## Environment Setup

1. Run the following to generate a config file

```
make setup
```

2. Edit `authgear-nft-indexer` for applicable configurations

## Database setup

1. Start the db container

```
docker-compose up -d db
```

2. Apply database schema migrations:

   make sure the db container is running

   ```sh
   go run ./cmd/indexer database migrate up
   ```

To create new migration:

```sh
# go run ./cmd/authgear database migrate new <migration name>
go run ./cmd/indexer database migrate new add user table
```

## Run everything

```
docker-compose up -d
```

Then run the following command to start up the indexer worker

```
# in project root
make start-worker
```
