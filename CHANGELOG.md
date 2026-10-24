# Changelog

## 1.3.1 - 2026-10-24

- Built with Go 1.27.3, which distrusts compromised intermediate certificates
  in `crypto/x509` (GO-2026-6809).
- Dependency updates

## 1.3.0 - 2026-10-09

- Dependency updates

## 1.2.0 - 2025-11-06

- Dependency updates

## 1.1.0 - 2025-05-07

- Added support for setting the log format via the `LOG_FORMAT` env var.
  Accepts either `text` (the default) or `json`.

## 1.0.0 - 2025-05-02

_Initial release._
