# JSON to Database Migration
This Go application converts a JSON file into a database migration file, with the ability to output to SQLite or Eloquent(more TBD). The goal of this application is adding context to DB migrations so that the business purposes are more visible and resistant to time.

The application takes in a json structure that outlines a migration and business need, we output a timestamped file with the appropriate format to be ingested by your migration tool/script.

## Usage
```sh
Usage of ./jason2migration:
  -f string
        Input file (default "input.json")
  -s string
        Strategy (sqlite, eloquent) (default "sqlite")
  -w    Write to file
```

### Example JSON input:
`Please refer to the example.json file in the root of the project`

## Known Limitations
- The application currently only supports SQLite, MySQL and Eloquent migration formats.
- When altering tables, we only support adding columns, not removing/renaming them.
- The application does not handle indexes (WIP)
- The application does not support foreign keys, primary keys, or other constraints
- The application is fairly brittle with constraints and data types, and will likely break if you provide an unsupported type or constraint.

## Unknown Limitations
- Probably a lot

## Running the Application
### To generate a SQLite migration file:

```sh
./jason2migration -f path/to/your/input.json -s sqlite -w
```

## Installation
To install the application, ensure you have Go installed, then run:

```sh
go get github.com/linkinlog/jason2migration
```

Binaries are also available for direct download in the releases tab.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
