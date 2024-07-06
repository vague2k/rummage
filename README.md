<h1 align="center">Rummage</h1>
<h3 align="center">A smart wrapper for "go get"</h3>

![rummage](https://github.com/vague2k/huez.nvim/assets/121782036/b9a85105-763e-4312-836b-eddb7b53408b)

## Installation

```
go install github.com/vague2k/rummage
```

## Alias

If you'd like a shorthand for the get command, you put this in your .zshrc/.bashrc file

```
alias rum="rummage get"
```

## Usage

Before using rummage regularly, It's reccommended that you use

```
rummage populate
```

This will get the database up to speed with the third party packages you have already installed.

| rummage [COMMAND] | Description                                                              |
| ----------------- | ------------------------------------------------------------------------ |
| `add`             | Add a package to the database.                                           |
| `remove`          | Removes a package from the database.                                     |
| `get`             | Gets a go package from the database, and increase its recency score.     |
| `populate`        | Add already installed packages to the database to quickstart your usage. |
| `query`           | Query the database to find a package.                                    |

## 📋 Contributing

Issues and PR's are always welcome and highly encouraged! I would love to learn more.

## License

[MIT](https://choosealicense.com/licenses/mit/)
