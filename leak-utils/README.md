# 🛠 leak-utils: The Eleakxir Data Utility Toolkit

`leak-utils` is a powerful command-line tool built to help you manage, process,
and optimize data leaks for use with the **Eleakxir** search engine. It provides
a suite of utilities for data cleaning, format conversion, and file
manipulation, all designed to ensure your data wells are efficient and
standardized.

`leak-utils` is written in **Go** and leverages **DuckDB** for its
high-performance in-memory processing, ensuring fast and reliable operations on
large datasets.

## 🚀 Features

- **Parquet File Management**: Clean and inspect existing `.parquet` files.
- **Format Conversion**: Seamlessly convert `.csv`, `.txt`, `.json` files into
  the optimized `.parquet` format.
- **Schema Uniformity**: Tools designed to help you standardize and normalize
  your data to align with the
  [Eleakxir data leak normalization rules](./DATALEAKS-NORMALIZATION.md). This
  ensures a consistent schema across all your files, which is crucial for
  efficient searching and consistent results.
- **High Performance**: Built with Go and DuckDB for fast and efficient data
  processing.

## ⚙️ How to Use

The tool operates via a single executable with different commands, each
corresponding to a specific action. You can find the executable in the
`leak-utils` directory of the Eleakxir project.

### Install

#### With go

```bash
go install "github.com/anotherhadi/eleakxir/leak-utils@latest"
```

#### With Nix/NixOS

<details>
<summary>Click to expand</summary>

**From anywhere (using the repo URL):**

```bash
nix run "github:anotherhadi/eleakxir#leak-utils" -- action [--flags value]
```

**Permanent Installation:**

```bash
# add the flake to your flake.nix
{
  inputs = {
    eleakxir.url = "github:anotherhadi/eleakxir";
  };
}

# then add it to your packages
environment.systemPackages = with pkgs; [ # or home.packages
    eleakxir.packages.${pkgs.system}.leak-utils
];
```

</details>

### Available Actions

#### `cleanParquet`

Optimizes and cleans an existing Parquet file. This can be used to change
columns, clean rows, ...

See:

```bash
leak-utils cleanParquet --help
```

#### `infoParquet`

Displays metadata and schema information for a given Parquet file. Useful for
inspecting file structure and column types.

#### `csvToParquet`

Converts a `.csv` file into a highly compressed and efficient `.parquet` file.
This is the recommended way to prepare your data for Eleakxir.

#### `mergeFiles`

Merges multiple files (of the same type) into a single, larger file. This is
useful for combining smaller data leaks.

#### `removeUrlSchemeFromUlp`

This utility prevents the colon (`:`) in URL schemes like `https://` from being
mistakenly parsed as a column separator when processing ULP data in flat files
like CSV or TXT.

## 🤝 Contributing

[Contributions](../CONTRIBUTING.md) to `leak-utils` are welcome! Feel free to
open issues or submit pull requests for new features, bug fixes, or performance
improvements.
