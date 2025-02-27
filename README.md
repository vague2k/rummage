<div align="center">
    <h1>Rummage</h1>
    <h3>A smart wrapper for "go get"</h3>
    <p>Rummage lets you get the packages you use most often with only a few keystrokes. <br>No more typing full package paths or copy pasting.</br></p>
</div>

![rummage](https://github.com/vague2k/huez.nvim/assets/121782036/b9a85105-763e-4312-836b-eddb7b53408b)

## Installation

```
go install github.com/vague2k/rummage@3.1.0
```

## Getting Started

Here's a list of commands available to rummage. Any command that supports arguements can take multiple arguement by default.

Use the help flag on any command to get more info about it

| rummage [COMMAND] | Description                                                                                    |
| ----------------- | ---------------------------------------------------------------------------------------------- |
| `add`             | Add an item that resembles a go package manually to the database.                              |
| `remove`          | Remove items from the database or be prompted to confirm to remove all items.                  |
| `get`             | Get a go package from the database using a substring, or get a package how you normally would. |
| `populate`        | Populate the database with third party packages already known by go.                           |
| `query`           | Query the database to find an entry by highest score, or using an exact match.                 |

Before using rummage regularly, It's reccommended that you use `populate` as
this will get the database up to speed with the third party packages you have already installed.

## Contributing

Issues and PR's are always welcome and highly encouraged! I would love to learn more.

## License

[MIT](https://choosealicense.com/licenses/mit/)
