
# Domain Processor Utility

This Go program automates the process of handling domain names by reading from a text file, executing domain-specific commands via `shuffledns`, and managing the results in a SQLite database. It is designed to run continuously, processing a list of domains every four hours.

## Features

- **Automated Domain Processing**: Automatically processes domains from a list and performs network operations on them.
- **Logging**: Detailed logs are maintained for each operation, with logs saved in a uniquely named file based on the execution time.
- **Error Handling**: Robust error handling and logging to capture and report issues during execution.
- **Database Integration**: Stores domain results in a SQLite database to ensure uniqueness and manage new entries.
- **Pause/Resume Capability**: Allows pausing and resuming the processing via standard input.
- **Signal Handling**: Gracefully handles termination signals to safely shutdown operations.

## Setup

### Prerequisites

- Go installed on your machine (see Go's [official documentation](https://golang.org/doc/install)).
- SQLite3 installed for database operations.
- `shuffledns` and `massdns` tools must be installed and accessible in your environment.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/z4z4x0/Add2DB_ShuffleDns.git
   ```
2. Navigate to the cloned directory:
   ```bash
   cd Add2DB_ShuffleDns
   ```

### Running the Program

1. Start the program by running:
   ```bash
   go run .
   ```
2. To pause/resume the processing, press the spacebar in the terminal where the program is running.

## Logging

Logs are generated in a new file each time the program is run, named with the current date and time, making it easy to track when each log file is generated.

## Database Schema

The program initializes the following SQLite schema if it doesn't exist:

```sql
CREATE TABLE IF NOT EXISTS domains (
    id INTEGER PRIMARY KEY,
    domain TEXT UNIQUE
);
```

Domains are inserted with an "INSERT OR IGNORE" strategy to avoid duplicates.

## Contributions

Contributions are welcome! Please fork the repository and submit a pull request with your suggested changes.

## License

This project is licensed under GPL-3.0 license. Please see the `LICENSE` file for more details.

## Contact 📧
@z4z4_h1
