<div align="center">
    <img alt="nixy logo" src="https://raw.githubusercontent.com/anotherhadi/eleakxir/main/.github/assets/logo.png" width="120px" />
</div>

<br>

# Eleakxir — Self-hosted search engine for leaked data.

[Eleakxir](https://eleakxir.hadi.diy) is a **self-hosted search engine** that
lets you connect to your own **private and secure server**, **explore data
wells** (parquet files) from multiple sources, and visualize results in a clean,
modern web interface.

> ✨ 100% open-source — you control your data, you control your server.

## 🚀 Features

- 🔐 **Private by design** — connect to your own Eleakxir server with a custom
  URL + password.
- 🎨 **Beautiful UI** — built with Svelte, TailwindCSS, and DaisyUI.
- 🛠 **Open source & extensible** — hack it, self-host it, extend it.
- **📁 Efficient File Format**: Uses the columnar **Parquet** format for high
  compression and rapid query performance.
- **🤖 Automatic Discovery**: Automatically detects new `.parquet` files in your
  folders and updates its metadata cache on a configurable schedule.
- **📜 Standardized Schema**: Includes a detailed guide on how to normalize your
  data leaks for consistent and effective searching across different breaches.
  (See [here](./leak-utils/DATALEAKS-NORMALIZATION.md))
- **🧰 Data Utility Toolkit**: Includes a dedicated command-line tool
  [leak-utils](./leak-utils/README.md) for managing, cleaning, and converting
  data files to the standardized Parquet format.
- **🔍 OSINT Tools**: Integration of various OSINT tools:
  - [github-recon](https://github.com/anotherhadi/github-recon)
  - [gravatar-recon](https://github.com/anotherhadi/gravatar-recon) (To-do)
  - sherlock (To-do)
  - holehe (To-do)
  - ghunt (To-do)

## ⚙️ How it works

1. You run an **Elixir server** that manages parquet files from various leaked
   data sources and multiple OSINT tools.
2. Eleakxir (the web client) connects to your server via HTTPS and authenticated
   headers.
3. You can:
   - Search across indexed leaks and OSINT tools
   - Browse results interactively
   - Review history and stats

## 🚨 Disclaimer

Eleakxir is provided **for educational and research purposes only**. You are
solely responsible for how you use this software. Accessing, storing, or
distributing leaked data may be illegal in your jurisdiction. The authors and
contributors **do not condone or promote illegal activity**. Use responsibly and
only with data you are legally permitted to process.

## 🧑‍💻 Tech stack

- **Frontend**: [Svelte](https://svelte.dev/),
  [sv-router](https://github.com/colinlienard/sv-router),
  [TailwindCSS](https://tailwindcss.com/), [DaisyUI](https://daisyui.com/)
- **Backend**: [Golang](https://go.dev), [Gin](https://gin-gonic.com),
  [duckdb](https://duckdb.org)

## 📦 Getting started

### Install with NixOS

1. In the `flake.nix` file, add `eleakxir` in the `inputs` section and import
   the `eleakxir.nixosModules.default` module:

```nix
{
  inputs = {
    eleakxir.url = "github:anotherhadi/eleakxir";
  };
  outputs = { 
    # ...
    modules = [
      inputs.eleakxir.nixosModules.eleakxir
    ];
    # ...
  }
}
```

2. Enable the backend service:

```nix
services.eleakxir = {
    enable = true;
    # port = 9198;
    folders = ["/var/lib/eleakxir/leaks/"] # Folders with parquet files
};
```

## 🔑 Configuration

### Backend

Check the [back.nix](./nix/back.nix) file to see configuration options.

### Client

Before searching, configure your server in the client:

1. Open [https://eleakxir.hadi.diy](https://eleakxir.hadi.diy) in your browser
   and add your server.
2. Click **“Connect your server”** in the UI.
3. Enter your **server URL** and **password**.
4. Start searching 🚀

## 🤝 Contributing

[Contributions are welcome](./CONTRIBUTING.md)! Feel free to open issues or
submit PRs.
