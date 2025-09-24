# Rules for handling Data Leaks

This normalization framework is designed to standardize data leaks for
[Eleakxir](https://github.com/anotherhadi/eleakxir), the open-source search
engine, using
[leak-utils](https://github.com/anotherhadi/eleakxir-temp/blob/main/leak-utils/README.md),
a dedicated CLI tool that converts and cleans files for efficient indexing and
searching.

## The Relevance of Parquet for Data Leaks

Parquet is an efficient, open-source columnar storage file format designed to
handle complex data in bulk. When dealing with data leaks, its choice is highly
relevant for several reasons:

- **Compression**: Parquet files offer superior compression compared to
  row-based formats like CSV. By storing data column by column, it applies more
  effective compression algorithms, which significantly reduces disk space. For
  data leaks, where file sizes can range from gigabytes to terabytes, this is
  crucial for minimizing storage costs.
- **Query Performance**: As a columnar format, Parquet allows you to read only
  the specific columns you need for a query. In a data leak, you might only be
  interested in emails and passwords, not full addresses or phone numbers. This
  selective reading drastically speeds up search operations, as the system
  doesn't have to scan through entire rows of irrelevant data.
- **Efficiency**: The format is optimized for analytics. It stores data with
  metadata and statistics (min/max values) for each column, allowing for query
  **pruning**. This means a query can skip entire blocks of data that don't
  match the filtering criteria, boosting performance even further.

## Disclaimer

The information in this document is provided **for research and educational
purposes only**. I am **not responsible** for how this data, methods, or
guidelines are used. Any misuse, unlawful activity, or harm resulting from
applying this content is the sole responsibility of the individual or
organization using it.

## File Naming Convention

- **Lowercase only**, ASCII (no accents).
- **Separators**:
  - `_` inside blocks (`date_2023_10`)
  - `-` between blocks (`instagram.com-date_2023_10`)
- **Prefix**: always start with the **source name/url** (e.g., `instagram.com`,
  `alien_txt`).
- **Blocks**: each additional part must be prefixed by its block name:

  - `date_YYYY[_MM[_DD]]` → use ISO format (year, or year-month, or full date).
  - `source_*` → origin of the leak (e.g., `scrape`, `dump`, `combo`).
  - `version_v*` → versioning if regenerated or transformed.
  - `notes_*` → optional clarifications.
- **Extension**: always `.parquet`.

**Recommended pattern:**

```txt
{source}-date_{YYYY[_MM[_DD]]}-source_{origin}-version_{vN}-notes_{info}.parquet
```

**Examples:**

```txt
instagram.com-date_2023_10.parquet
alien_txt-date_2022-source_dump.parquet
combo_french-notes_crypto.parquet
```

## Column Naming Convention

- **snake\_case only** (lowercase, `_` separator).

- **No dots (`.`)** in column names (`husband.phone` → `husband_phone`).

- **Allowed characters**: `[a-z0-9_]+` (no spaces, hyphens, or accents).

- **Multiple variants of the same field**:

  - Relations → prefix clearly: `husband_phone`, `mother_last_name`.
  - Multiples of the same type → numbered prefix: `1_phone`, `2_phone`,
    `3_phone`.
  - Always end with the column "type" (e.g., `_phone`, `_last_name`).

- **Rename if mislabeled**: If a `username` column actually contains only emails
  rename it to `email`.

- **Remove irrelevant columns**: Drop meaningless identifiers like `id` or
  fields with no analytical value.

- **Standard columns**: to enable schema alignment across leaks:

  | Column        |
  | ------------- |
  | email         |
  | username      |
  | password      |
  | password_hash |
  | phone         |
  | date          |
  | birth_date    |
  | age           |
  | first_name    |
  | last_name     |
  | full_name     |
  | address       |
  | city          |
  | country       |
  | state         |
  | postal_code   |
  | ip            |
  | url           |
  | city          |

## Standard Column Formatting

- **Email**: lowercase, trimmed, keep only `[^a-z0-9._@-]`.

- **Phone**: keep only `[^0-9]`

- **Names**:

  - Keep `first_name` / `last_name` if present.
  - Generate `full_name = CONCAT(first_name, ' ', last_name)`.
  - If only `name` exists, rename it to `full_name`.

- **Passwords**:

  - Hashes → `password_hash`.
  - Plaintext → `password`.
  - Never mix hashes and plaintext in the same column.

- **NULLs**: always use SQL `NULL` (never `""` or `"NULL"`).

## Deduplication

Deduplication is often **impractical at scale** (billions of rows). Do **not**
attempt to deduplicate at ingestion time. Instead, handle deduplication **after
running a search** to optimize performance and storage.
