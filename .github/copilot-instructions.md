# Code Quality Instructions

- Do not write comments. If you feel the need to write comments, make the code more readable instead.
- Use `slog` to log messages. Use the Debug level for code flow messages, Info for letting the user know of changes like directory/file creation, Warn for recoverable issues, and Error for serious problems.
- Use the methods of the `iac/utils` module wherever possible. If there is no suitable method, consider adding one.
- Consider using the `iac/utils/errors` module for creating custom errors. If there is need to exit the program with a specific exit code, use the `iac/utils/exitcodes` module.
