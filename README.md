# paphos-backend

This is the core backend for the official Pygmalion service.

Very early work-in-progress, not usable in any capacity.

## Contributing

If you wish to contribute, this section contains some relevant information.

The stack is:

- [Crystal](https://crystal-lang.org/) as the language of choice
- [Lucky](https://luckyframework.org/) as the web framework
- [PostgreSQL](https://www.postgresql.org/) as the main database

To get started, you need all of the above installed and functional on your machine.

For development, the [default configuration](./config/database.cr) expects a `postgres` user with password `postgres` to be available. You can then create the databases manually, or use `lucky db.setup` to create them automatically I believe.

With PostgreSQL up and running and the development database created, you can boot up the development server:

```bash
$ lucky dev
```

The API should then be reachable at http://localhost:3000/. Routes are under `/api/v1` but are undocumented as of now, you'll need to [trek through the code](./src/actions/api/v1) to figure them out.

### Useful links

- [Lucky's official guides](https://luckyframework.org/guides/getting-started/installing): contains tutorials about the usual stuff (installation, route setup, handling data in the database, etc.)
