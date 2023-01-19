# paphos-backend

This is the backend for the official Pygmalion service.

Very early work-in-progress, not usable in any capacity.

## Contributing

If you wish to contribute, this section contains some relevant information.

The stack is:

- [Golang](https://go.dev/) as the language of choice
- [Buffalo](https://gobuffalo.io/) as the web framework
- [PostgreSQL](https://www.postgresql.org/) as the main database

To get started, you need all of the above installed and functional on your machine.

For development, the [default configuration](./database.yml) expects a `paphos` user with password `paphos` to be available. You can then create the databases manually, or use `buffalo pop create -a`.

With PostgreSQL up and running and the development database created, you can apply the migrations:

```bash
$ buffalo pop migrate
```

And then boot up the development server:

```bash
$ PORT=3000 buffalo dev
```

The API should then be reachable at http://localhost:3000/.

### Useful links

- [fizz](https://github.com/gobuffalo/fizz/blob/main/README.md)'s documentation (migration DSL)
